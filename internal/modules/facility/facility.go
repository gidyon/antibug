package facility

import (
	"context"
	"errors"
	"fmt"
	"github.com/gidyon/antibug/internal/modules"
	"github.com/gidyon/antibug/internal/pkg/auth"
	"github.com/gidyon/antibug/internal/pkg/errs"
	"github.com/gidyon/antibug/pkg/api/facility"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/grpclog"
	"strings"
)

type facilityAPIServer struct {
	sqlDB       *gorm.DB
	logger      grpclog.LoggerV2
	authAPI     auth.Interface
	counties    []*facility.County
	subCounties []*facility.SubCounty
	data        map[string]*facility.SubCounty
}

// Options contains parameters to new facility API
type Options struct {
	SQLDB            *gorm.DB
	Logger           grpclog.LoggerV2
	JWTSigningKey    string
	CountiesDataFile string
}

// NewFacilityAPI creates a new facility API server
func NewFacilityAPI(ctx context.Context, opt *Options) (facility.FacilityAPIServer, error) {
	var err error
	switch {
	case ctx == nil:
		err = errs.NilObject("Context")
	case opt.SQLDB == nil:
		err = errs.NilObject("SqlDB")
	case opt.Logger == nil:
		err = errs.NilObject("Logger")
	case opt.JWTSigningKey == "":
		err = errs.MissingField("JWTSigning Key")
	}
	if err != nil {
		return nil, err
	}

	authAPI, err := auth.NewAPI(opt.JWTSigningKey)
	if err != nil {
		return nil, err
	}

	fapi := &facilityAPIServer{
		sqlDB:       opt.SQLDB,
		logger:      opt.Logger,
		authAPI:     authAPI,
		counties:    make([]*facility.County, 0),
		subCounties: make([]*facility.SubCounty, 0),
		data:        make(map[string]*facility.SubCounty, 0),
	}

	// Perform auto migration
	err = fapi.sqlDB.AutoMigrate(&Facility{}, &County{}, &SubCounty{}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to automigrate facilities table: %w", err)
	}

	// Load counties and subcounties
	err = fapi.loadCountiesData()
	if err != nil {
		return nil, err
	}

	// Create a full text search index
	err = modules.CreateFullTextIndex(fapi.sqlDB, facilitiesTable, "facility_name")
	if err != nil {
		return nil, fmt.Errorf("failed to create full text index: %v", err)
	}

	return fapi, nil
}

func (fapi *facilityAPIServer) AddFacility(
	ctx context.Context, addReq *facility.AddFacilityRequest,
) (*facility.AddFacilityResponse, error) {
	// Request must not be nil
	if addReq == nil {
		return nil, errs.NilObject("AddFacilityRequest")
	}

	// Authorize request
	_, err := fapi.authAPI.AuthorizeGroup(ctx, auth.Admin)
	if err != nil {
		return nil, err
	}

	facilityPB := addReq.GetFacility()

	// Validation
	switch {
	case strings.TrimSpace(facilityPB.FacilityName) == "":
		err = errs.MissingField("facility name")
	case strings.TrimSpace(facilityPB.County) == "":
		err = errs.MissingField("count name")
	case facilityPB.CountyCode == 0:
		err = errs.MissingField("county code")
	case strings.TrimSpace(facilityPB.SubCounty) == "":
		err = errs.MissingField("sub county")
	case facilityPB.SubCountyCode == 0:
		err = errs.MissingField("sub county code")
	}
	if err != nil {
		return nil, err
	}

	// Get database model
	facilityDB, err := getFacilityDB(facilityPB)
	if err != nil {
		return nil, err
	}

	// Create in database
	err = fapi.sqlDB.Create(facilityDB).Error
	switch {
	case err == nil:
	case strings.Contains(strings.ToLower(err.Error()), "duplicate"):
		return nil, errs.DuplicateField("facility name", facilityDB.FacilityName)
	default:
		return nil, errs.SQLQueryFailed(err, "CREATE")
	}

	// That's it!
	return &facility.AddFacilityResponse{
		FacilityId: fmt.Sprint(int64(facilityDB.ID)),
	}, nil
}

func (fapi *facilityAPIServer) RemoveFacility(
	ctx context.Context, delReq *facility.RemoveFacilityRequest,
) (*empty.Empty, error) {
	// Request must not be nil
	if delReq == nil {
		return nil, errs.NilObject("RemoveFacilityRequest")
	}

	// Authorize request
	_, err := fapi.authAPI.AuthorizeGroup(ctx, auth.Admin)
	if err != nil {
		return nil, err
	}

	if delReq.GetFacilityId() == "" {
		return nil, errs.MissingField("facility id")
	}

	// Delete in database
	err = fapi.sqlDB.Table(facilitiesTable).Delete(&Facility{}, "id=?", delReq.FacilityId).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "DELETE")
	}

	return &empty.Empty{}, nil
}

func (fapi *facilityAPIServer) GetFacility(
	ctx context.Context, getReq *facility.GetFacilityRequest,
) (*facility.Facility, error) {
	// Request must not be nil
	if getReq == nil {
		return nil, errs.NilObject("GetFacilityRequest")
	}

	// Authorize request
	err := fapi.authAPI.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Request must not be nil
	if getReq == nil {
		return nil, errs.NilObject("GetFacilityRequest")
	}

	// Validation
	if getReq.FacilityId == "" {
		return nil, errs.MissingField("facility id")
	}

	facilityDB := &Facility{}

	err = fapi.sqlDB.First(facilityDB, "id=?", getReq.FacilityId).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errs.NotFound("facility", getReq.FacilityId)
	default:
		return nil, errs.SQLQueryFailed(err, "GET")
	}

	// Get facility pb
	facilityPB, err := getFacilityPB(facilityDB)
	if err != nil {
		return nil, err
	}

	return facilityPB, nil
}

const defaultPageSize = 50

func normalizePageSixe(pageToken, pageSize int32) (int, int) {
	if pageToken <= 0 {
		pageToken = 0
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if pageSize > defaultPageSize {
		pageSize = defaultPageSize
	}
	return int(pageToken), int(pageSize)
}

func (fapi *facilityAPIServer) ListFacilities(
	ctx context.Context, listReq *facility.ListFacilitiesRequest,
) (*facility.Facilities, error) {
	// Request must not be nil
	if listReq == nil {
		return nil, errs.NilObject("ListFacilitiesRequest")
	}

	// Authorize request
	err := fapi.authAPI.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Normalize page
	pageToken, pageSize := normalizePageSixe(listReq.PageToken, listReq.PageSize)

	facilitiesDB := make([]*Facility, 0, pageSize)
	err = fapi.sqlDB.Order("created_at DESC").Limit(pageSize).Order("id, created_at ASC").
		Where("id>?", pageToken).Find(&facilitiesDB).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "LIST")
	}

	facilitiesPB := make([]*facility.Facility, 0, len(facilitiesDB))

	for _, facilityDB := range facilitiesDB {
		facilityPB, err := getFacilityPB(facilityDB)
		if err != nil {
			return nil, err
		}
		facilitiesPB = append(facilitiesPB, facilityPB)
		pageToken = int(facilityDB.ID)
	}

	return &facility.Facilities{
		Facilities:    facilitiesPB,
		NextPageToken: int32(pageToken),
	}, nil
}

func (fapi *facilityAPIServer) SearchFacilities(
	ctx context.Context, searchReq *facility.SearchFacilitiesRequest,
) (*facility.Facilities, error) {
	// Request must not be nil
	if searchReq == nil {
		return nil, errs.NilObject("SearchFacilitiesRequest")
	}

	// Authorize request
	err := fapi.authAPI.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// For empty queries
	if searchReq.Query == "" {
		return &facility.Facilities{
			Facilities: []*facility.Facility{},
		}, nil
	}

	pageToken, pageSize := normalizePageSixe(searchReq.PageToken, searchReq.PageSize)

	parsedQuery := modules.ParseQuery(searchReq.Query, " facilities", "facilities")

	facilitiesDB := make([]*Facility, 0, pageSize)

	err = fapi.sqlDB.Unscoped().Limit(pageSize).Order("id, created_at ASC").
		Where("id>?", pageToken).
		Find(&facilitiesDB, "MATCH(facility_name) AGAINST(? IN BOOLEAN MODE)", parsedQuery).Error
	switch {
	case err == nil:
	default:
		return nil, errs.SQLQueryFailed(err, "SELECT")
	}

	// Populate response
	facilitiesPB := make([]*facility.Facility, 0, len(facilitiesDB))

	for _, facilityDB := range facilitiesDB {
		facilityPB, err := getFacilityPB(facilityDB)
		if err != nil {
			return nil, err
		}
		facilitiesPB = append(facilitiesPB, facilityPB)
		pageToken = int(facilityDB.ID)
	}

	return &facility.Facilities{
		NextPageToken: int32(pageToken),
		Facilities:    facilitiesPB,
	}, nil
}

func (fapi *facilityAPIServer) ListCounties(
	ctx context.Context, _ *empty.Empty,
) (*facility.Counties, error) {
	// Authorize request
	err := fapi.authAPI.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	return &facility.Counties{
		Counties: fapi.counties,
	}, nil
}

func (fapi *facilityAPIServer) ListSubCounties(
	ctx context.Context, _ *empty.Empty,
) (*facility.SubCounties, error) {
	// Authorize request
	err := fapi.authAPI.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	return &facility.SubCounties{
		SubCounties: fapi.subCounties,
	}, nil
}
