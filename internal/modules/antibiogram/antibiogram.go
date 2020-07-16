package antibiogram

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/gidyon/antibug/internal/modules/culture"
	"github.com/gidyon/antibug/internal/pkg/auth"
	"github.com/gidyon/antibug/internal/pkg/errs"
	antibiogram "github.com/gidyon/antibug/pkg/api/antibiogram"
	culture_pb "github.com/gidyon/antibug/pkg/api/culture"
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"sort"
	"strings"
	"time"
)

type apiServer struct {
	sqlDB       *gorm.DB
	redisClient *redis.Client
	logger      grpclog.LoggerV2
	authAPI     auth.Interface
}

// Options contains parameters that is passed to NewAntibiogramAPIServer factory
type Options struct {
	SQLDB         *gorm.DB
	RedisDB       *redis.Client
	Logger        grpclog.LoggerV2
	JWTSigningKey string
}

// NewAntibiogramAPIServer is a factory for new AntibiogramAPIServer service
func NewAntibiogramAPIServer(
	ctx context.Context, opt *Options,
) (antibiogram.AntibiogramAPIServer, error) {
	// Validation
	var err error
	switch {
	case ctx == nil:
		err = errs.NilObject("Context")
	case opt.SQLDB == nil:
		err = errs.NilObject("SqlDB")
	case opt.RedisDB == nil:
		err = errs.NilObject("RedisClient")
	case opt.Logger == nil:
		err = errs.NilObject("Logger")
	case opt.JWTSigningKey == "":
		err = errs.MissingField("JWTSigningKey")
	}
	if err != nil {
		return nil, err
	}

	authAPI, err := auth.NewAPI(opt.JWTSigningKey)
	if err != nil {
		return nil, err
	}

	api := &apiServer{
		sqlDB:       opt.SQLDB,
		redisClient: opt.RedisDB,
		logger:      opt.Logger,
		authAPI:     authAPI,
	}

	// Auto migration
	err = api.sqlDB.AutoMigrate(&culture.Culture{}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate cultures table: %v", err)
	}

	return api, nil
}

const (
	sixMonths      = time.Hour * 24 * 30 * 6
	twelveMonths   = time.Hour * 24 * 30 * 12
	twoYears       = time.Hour * 24 * 30 * 24
	fourYears      = time.Hour * 24 * 30 * 48
	eightYears     = time.Hour * 24 * 30 * 96
	sixteenYears   = time.Hour * 24 * 30 * 192
	thirtytwoYears = time.Hour * 24 * 30 * 384
)

func buildQuery(sqlDB *gorm.DB, filter *antibiogram.Filter) *gorm.DB {
	// Duration
	now := time.Now().Unix()
	switch filter.PastDuration {
	case antibiogram.Duration_PAST_SIX_MONTHS:
		dur := now - int64(sixMonths)
		sqlDB = sqlDB.Where("results_timestamp_sec>=?", dur)
	case antibiogram.Duration_PAST_ONE_YEARS:
		dur := now - int64(twelveMonths)
		sqlDB = sqlDB.Where("results_timestamp_sec>=?", dur)
	case antibiogram.Duration_PAST_TWO_YEARS:
		dur := now - int64(twoYears)
		sqlDB = sqlDB.Where("results_timestamp_sec>=?", dur)
	case antibiogram.Duration_PAST_FOUR_YEARS:
		dur := now - int64(fourYears)
		sqlDB = sqlDB.Where("results_timestamp_sec>=?", dur)
	case antibiogram.Duration_PAST_EIGHT_YEARS:
		dur := now - int64(eightYears)
		sqlDB = sqlDB.Where("results_timestamp_sec>=?", dur)
	case antibiogram.Duration_PAST_SIXTEEN_YEARS:
		dur := now - int64(sixteenYears)
		sqlDB = sqlDB.Where("results_timestamp_sec>=?", dur)
	case antibiogram.Duration_PAST_THIRTY_TWO_YEARS:
		dur := now - int64(thirtytwoYears)
		sqlDB = sqlDB.Where("results_timestamp_sec>=?", dur)
	}

	// RegionScope
	switch filter.RegionScope {
	case antibiogram.RegionScope_COUNTRY:
	case antibiogram.RegionScope_COUNTY:
		if len(filter.ScopeValues) > 0 {
			sqlDB = sqlDB.Where("county_code IN(?)", filter.ScopeValues)
		}
	case antibiogram.RegionScope_SUB_COUNTY:
		if len(filter.ScopeValues) > 0 {
			sqlDB = sqlDB.Where("sub_county_code IN(?)", filter.ScopeValues)
		}
	case antibiogram.RegionScope_FACILITY:
		if len(filter.ScopeValues) > 0 {
			sqlDB = sqlDB.Where("hospital_id IN(?)", filter.ScopeValues)
		}
	}

	// Advanced => Gender, Age
	if filter.GetAdvanced() {
		advancedFilter := filter.GetAdvance()
		if advancedFilter != nil {
			// Gender
			switch advancedFilter.GetGender() {
			case antibiogram.Gender_ALL:
				sqlDB = sqlDB.Where("patient_gender IN(?)", []string{"male", "female", "all"})
			case antibiogram.Gender_MALE:
				sqlDB = sqlDB.Where("patient_gender=?", "male")
			case antibiogram.Gender_FEMALE:
				sqlDB = sqlDB.Where("patient_gender=?", "female")
			}

			// Age
			if advancedFilter.GetAgeMinDays() < advancedFilter.GetAgeMaxDays() {
				sqlDB = sqlDB.Where(
					"patient_age BETWEEN ? AND ?",
					advancedFilter.GetAgeMinDays()/(365), advancedFilter.GetAgeMaxDays()/(365),
				)
			}
		}
	}

	return sqlDB
}

func maxLabel(labels map[culture_pb.Label]int32) culture_pb.Label {
	var (
		label culture_pb.Label
		score int32
	)
	for key, value := range labels {
		if value > score {
			label = key
		}
	}
	return label
}

func genFilterHash(filter *antibiogram.Filter) string {
	// Filter criteria
	str := fmt.Sprintf(
		"%d%d%s%t",
		filter.GetPastDuration(),
		filter.GetRegionScope(),
		strings.Join(filter.GetScopeValues(), ","),
		filter.GetAdvanced(),
	)
	// Input values
	if len(filter.GetInputValues()) > 0 {
		inputValues := make([]string, 0, len(filter.InputValues))
		for _, inputValue := range filter.GetInputValues() {
			inputValues = append(inputValues, inputValue.GetId())
		}
		sort.Strings(inputValues)
		str += strings.Join(inputValues, ",")
	}
	// Advance
	if filter.GetAdvance() != nil {
		str += fmt.Sprintf(
			"%d%d%d",
			filter.Advance.GetGender(),
			filter.Advance.GetAgeMaxDays(),
			filter.Advance.GetAgeMinDays(),
		)
	}

	// Apply hash
	sum := sha256.Sum224([]byte(str))

	return string(sum[:])
}

func (api *apiServer) getPathogenAntibiogramFromCache(
	ctx context.Context, filter *antibiogram.Filter, index int,
) (*antibiogram.PathogenAntibiogram, error) {
	// Get hash of filter query
	filterHash := genFilterHash(filter)

	// Check cache if it exists
	data, err := api.redisClient.Get(ctx, filterHash).Result()
	switch {
	case err == nil:
	case errors.Is(err, redis.Nil):
		return api.getPathogenAntibiogram(ctx, filterHash, filter, index)
	default:
		return nil, errs.RedisCmdFailed(err, "GET")
	}

	// Unmarshal data
	pathogenAntibiogram := &antibiogram.PathogenAntibiogram{}
	err = proto.Unmarshal([]byte(data), pathogenAntibiogram)
	if err != nil {
		return nil, errs.WrapErrWithMessage(codes.Internal, err, "failed to proto unmarshal")
	}

	return pathogenAntibiogram, nil
}

func (api *apiServer) getPathogenAntibiogram(
	ctx context.Context, queryHash string, filter *antibiogram.Filter, index int,
) (*antibiogram.PathogenAntibiogram, error) {

	culturesDB := make([]*culture.Culture, 0, 500)

	pathogenPB := filter.GetInputValues()[index]

	// Parse filter
	sqlDB := buildQuery(api.sqlDB, filter).Where("? MEMBER OF(pathogens_found)", pathogenPB.GetId())
	err := sqlDB.Find(&culturesDB).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "SELECT")
	}

	pathogenAntibiogram := &antibiogram.PathogenAntibiogram{
		PathogenName:     pathogenPB.GetName(),
		PathogenId:       pathogenPB.GetId(),
		Susceptibilities: make([]*antibiogram.PathogenSusceptibility, 0),
	}

	pathogenSusceptibilities := make(map[string]*antibiogram.PathogenSusceptibility, 0)

	averageLabel := make(map[culture_pb.Label]int32, 0)

	// Range over results and populate response slice
	for _, cultureDB := range culturesDB {
		culturePB, err := culture.GetCulturePB(cultureDB)
		if err != nil {
			return nil, err
		}

		for _, cultureResult := range culturePB.GetCultureResults() {
			if cultureResult.GetPathogenId() == pathogenPB.GetId() {
				pathogenSusceptibility, ok := pathogenSusceptibilities[pathogenPB.GetId()]
				if !ok {
					pathogenSusceptibilities[pathogenPB.GetId()] = &antibiogram.PathogenSusceptibility{
						AntimicrobialName:   cultureResult.GetAntimicrobialName(),
						AntimicrobialId:     cultureResult.GetAntimicrobialId(),
						Isolates:            1,
						SusceptibilityScore: cultureResult.GetSusceptibilityScore(),
						Label:               cultureResult.Label,
					}
					averageLabel[cultureResult.Label]++
					continue
				}

				pathogenSusceptibility.Isolates++

				// Average susceptibility score
				pathogenSusceptibility.SusceptibilityScore += cultureResult.GetSusceptibilityScore()
				pathogenSusceptibility.SusceptibilityScore = pathogenSusceptibility.SusceptibilityScore / 2
				// Average label score
				averageLabel[cultureResult.Label]++
				pathogenSusceptibility.Label = maxLabel(averageLabel)
			}
		}
	}

	// Add individual susceptibility to list of susceptibilities
	for _, val := range pathogenSusceptibilities {
		pathogenAntibiogram.Susceptibilities = append(pathogenAntibiogram.Susceptibilities, val)
	}

	// Marshal data
	bs, err := proto.Marshal(pathogenAntibiogram)
	if err != nil {
		return nil, errs.WrapErrWithMessage(codes.Internal, err, "failed to proto marshal")
	}

	// Save to cache
	err = api.redisClient.Set(ctx, queryHash, bs, time.Hour*24*7).Err()
	if err != nil {
		return nil, errs.WrapErrWithMessage(codes.Internal, err, "fail to save antibiogram from cache")
	}

	return pathogenAntibiogram, nil
}

func validateFilter(filter *antibiogram.Filter) error {
	var err error
	switch {
	case len(filter.GetInputValues()) == 0:
		err = errs.MissingField("InputValues")
	case filter.RegionScope != antibiogram.RegionScope_COUNTRY && len(filter.GetScopeValues()) == 0:
		err = errs.MissingField("ScopeValues")
	}
	return err
}

func (api *apiServer) GenPathogensAntibiogram(
	ctx context.Context, filter *antibiogram.Filter,
) (*antibiogram.PathogensAntibiogram, error) {
	// Request must not be nil
	if filter == nil {
		return nil, errs.NilObject("Filter")
	}

	// Authentication
	err := api.authAPI.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Validation
	err = validateFilter(filter)
	if err != nil {
		return nil, err
	}

	pathogensAntibiogram := make([]*antibiogram.PathogenAntibiogram, 0, len(filter.GetInputValues()))

	// Get individual susceptibility
	for index := range filter.GetInputValues() {
		pathogenAntibiogram, err := api.getPathogenAntibiogramFromCache(ctx, filter, index)
		if err != nil {
			api.logger.Errorf("error getting pathogen from cache: %v", err)
			continue
		}

		pathogensAntibiogram = append(pathogensAntibiogram, pathogenAntibiogram)
	}

	return &antibiogram.PathogensAntibiogram{
		Antibiograms: pathogensAntibiogram,
	}, nil
}

func (api *apiServer) GenPathogenAntibiogram(
	ctx context.Context, filter *antibiogram.Filter,
) (*antibiogram.PathogenAntibiogram, error) {
	// Request must not be nil
	if filter == nil {
		return nil, errs.NilObject("Filter")
	}

	// Authentication
	err := api.authAPI.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Validation
	err = validateFilter(filter)
	if err != nil {
		return nil, err
	}

	// Get antibiogram from filter
	return api.getPathogenAntibiogramFromCache(ctx, filter, 0)
}

func (api *apiServer) getAntimicrobialAntibiogramFromCache(
	ctx context.Context, filter *antibiogram.Filter, index int,
) (*antibiogram.AntimicrobialAntibiogram, error) {
	// Get hash of filter query
	filterHash := genFilterHash(filter)

	// Check cache if it exists
	data, err := api.redisClient.Get(ctx, filterHash).Result()
	switch {
	case err == nil:
	case errors.Is(err, redis.Nil):
		return api.getAntimicrobialAntibiogram(ctx, filterHash, filter, index)
	default:
		return nil, errs.RedisCmdFailed(err, "GET")
	}

	// Unmarshal data
	antimicrobialAntibiogram := &antibiogram.AntimicrobialAntibiogram{}
	err = proto.Unmarshal([]byte(data), antimicrobialAntibiogram)
	if err != nil {
		return nil, errs.WrapErrWithMessage(codes.Internal, err, "failed to proto unmarshal")
	}

	return antimicrobialAntibiogram, nil
}

func (api *apiServer) getAntimicrobialAntibiogram(
	ctx context.Context, queryHash string, filter *antibiogram.Filter, index int,
) (*antibiogram.AntimicrobialAntibiogram, error) {

	culturesDB := make([]*culture.Culture, 0, 500)

	antimicrobialPB := filter.GetInputValues()[index]

	// Parse filter
	sqlDB := buildQuery(api.sqlDB, filter).Where("? MEMBER OF(antimicrobials_used)", antimicrobialPB.GetId())
	err := sqlDB.Find(&culturesDB).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "SELECT")
	}

	antimicrobialAntibiogram := &antibiogram.AntimicrobialAntibiogram{
		AntimicrobialName: antimicrobialPB.GetName(),
		AntimicrobialId:   antimicrobialPB.GetId(),
		Susceptibilities:  make([]*antibiogram.AntimicrobialSusceptibility, 0),
	}

	antimicrobialSusceptibilities := make(map[string]*antibiogram.AntimicrobialSusceptibility, 0)

	averageLabel := make(map[culture_pb.Label]int32, 0)

	// Range over results and populate response slice
	for _, cultureDB := range culturesDB {
		culturePB, err := culture.GetCulturePB(cultureDB)
		if err != nil {
			return nil, err
		}

		for _, cultureResult := range culturePB.GetCultureResults() {
			if cultureResult.GetAntimicrobialId() == antimicrobialPB.GetId() {
				antimicrobialSusceptibility, ok := antimicrobialSusceptibilities[antimicrobialPB.GetId()]
				if !ok {
					antimicrobialSusceptibilities[antimicrobialPB.GetId()] = &antibiogram.AntimicrobialSusceptibility{
						PathogenName:        cultureResult.GetPathogenName(),
						PathogenId:          cultureResult.GetPathogenId(),
						Isolates:            1,
						SusceptibilityScore: cultureResult.GetSusceptibilityScore(),
						Label:               cultureResult.Label,
					}
					averageLabel[cultureResult.Label]++
					continue
				}

				antimicrobialSusceptibility.Isolates++

				// Average susceptibility score
				antimicrobialSusceptibility.SusceptibilityScore += cultureResult.GetSusceptibilityScore()
				antimicrobialSusceptibility.SusceptibilityScore = antimicrobialSusceptibility.SusceptibilityScore / 2

				// Average label score
				averageLabel[cultureResult.Label]++
				antimicrobialSusceptibility.Label = maxLabel(averageLabel)
			}
		}
	}

	// Add individual susceptibility to list of susceptibilities
	for _, val := range antimicrobialSusceptibilities {
		antimicrobialAntibiogram.Susceptibilities = append(antimicrobialAntibiogram.Susceptibilities, val)
	}

	// Marshal data
	bs, err := proto.Marshal(antimicrobialAntibiogram)
	if err != nil {
		return nil, errs.WrapErrWithMessage(codes.Internal, err, "failed to proto marshal")
	}

	// Save to cache
	err = api.redisClient.Set(ctx, queryHash, bs, time.Hour*24*7).Err()
	if err != nil {
		return nil, errs.WrapErrWithMessage(codes.Internal, err, "fail to save antibiogram from cache")
	}

	return antimicrobialAntibiogram, nil
}

func (api *apiServer) GenAntimicrobialsAntibiogram(
	ctx context.Context, filter *antibiogram.Filter,
) (*antibiogram.AntimicrobialsAntibiogram, error) {
	// Request must not be nil
	if filter == nil {
		return nil, errs.NilObject("Filter")
	}

	// Authentication
	err := api.authAPI.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Validation
	err = validateFilter(filter)
	if err != nil {
		return nil, err
	}

	antimicrobialsAntibiogram := make([]*antibiogram.AntimicrobialAntibiogram, 0, len(filter.GetInputValues()))

	// Get individual susceptibility
	for index := range filter.GetInputValues() {
		antimicrobialAntibiogram, err := api.getAntimicrobialAntibiogramFromCache(ctx, filter, index)
		if err != nil {
			api.logger.Errorf("error getting antimicrobial from cache: %v", err)
			continue
		}

		antimicrobialsAntibiogram = append(antimicrobialsAntibiogram, antimicrobialAntibiogram)
	}

	return &antibiogram.AntimicrobialsAntibiogram{
		Antibiograms: antimicrobialsAntibiogram,
	}, nil
}

func (api *apiServer) GenAntimicrobialAntibiogram(
	ctx context.Context, filter *antibiogram.Filter,
) (*antibiogram.AntimicrobialAntibiogram, error) {
	// Request must not be nil
	if filter == nil {
		return nil, errs.NilObject("Filter")
	}

	// Authentication
	err := api.authAPI.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Validation
	err = validateFilter(filter)
	if err != nil {
		return nil, err
	}

	// Get antibiogram from filter
	return api.getAntimicrobialAntibiogramFromCache(ctx, filter, 0)
}
