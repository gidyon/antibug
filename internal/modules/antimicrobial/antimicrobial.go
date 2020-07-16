package antimicrobial

import (
	"context"
	"errors"
	"fmt"
	"github.com/gidyon/antibug/internal/modules"
	"github.com/gidyon/antibug/internal/pkg/auth"
	"github.com/gidyon/antibug/internal/pkg/errs"
	"github.com/gidyon/antibug/pkg/api/antimicrobial"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/grpclog"
	"strings"
)

var (
	createAllowedGroups = []string{auth.Physician, auth.Researcher, auth.LabTechnician, auth.Admin, auth.Pharmacist}
	deleteAllowedGroups = []string{auth.Physician, auth.Researcher, auth.Admin}
)

type antimicrobialAPIServer struct {
	sqlDB   *gorm.DB
	logger  grpclog.LoggerV2
	authAPI auth.Interface
}

// Options contains parameters for NewAntimicrobialAPI
type Options struct {
	SQLDB      *gorm.DB
	Logger     grpclog.LoggerV2
	SigningKey string
}

// NewAntimicrobialAPI creates a new antimicrobial API server
func NewAntimicrobialAPI(ctx context.Context, opt *Options) (antimicrobial.AntimicrobialAPIServer, error) {
	// Validation
	var err error
	switch {
	case opt.SQLDB == nil:
		err = errs.NilObject("SqlDB")
	case opt.Logger == nil:
		err = errs.NilObject("Logger")
	case opt.SigningKey == "":
		err = errs.MissingField("Jwt SigningKey")
	case ctx == nil:
		err = errs.NilObject("Context")
	}
	if err != nil {
		return nil, err
	}

	authAPI, err := auth.NewAPI(opt.SigningKey)
	if err != nil {
		return nil, err
	}

	papi := &antimicrobialAPIServer{
		sqlDB:   opt.SQLDB,
		logger:  opt.Logger,
		authAPI: authAPI,
	}

	// Perform auto migration
	err = papi.sqlDB.AutoMigrate(&Antimicrobial{}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to automigrate antimicrobials table: %w", err)
	}

	// Create a full text search index
	err = modules.CreateFullTextIndex(papi.sqlDB, antimicrobialsTable, "antimicrobial_name")
	if err != nil {
		return nil, fmt.Errorf("failed to create full text index: %v", err)
	}

	return papi, nil
}

func (papi *antimicrobialAPIServer) CreateAntimicrobial(
	ctx context.Context, createReq *antimicrobial.CreateAntimicrobialRequest,
) (*antimicrobial.CreateAntimicrobialResponse, error) {
	// Request must not be nil
	if createReq == nil {
		return nil, errs.NilObject("CreateAntimicrobialRequest")
	}

	// Authorize request
	_, err := papi.authAPI.AuthorizeGroup(ctx, createAllowedGroups...)
	if err != nil {
		return nil, err
	}

	antimicrobialPB := createReq.GetAntimicrobial()

	// Validation
	switch {
	case strings.TrimSpace(antimicrobialPB.AntimicrobialName) == "":
		err = errs.MissingField("AntimicrobialName")
	case strings.TrimSpace(antimicrobialPB.CDiff) == "":
		err = errs.MissingField("CDiff")
	case strings.TrimSpace(antimicrobialPB.OralBioavailability) == "":
		err = errs.MissingField("OralBioavailability")
	case strings.TrimSpace(antimicrobialPB.ApproximateCost) == "":
		err = errs.MissingField("ApproximateCost")
	case antimicrobialPB.GeneralUsage == nil ||
		len(antimicrobialPB.GeneralUsage.Values) == 0:
		err = errs.MissingField("GeneralUsage.Values")
	case antimicrobialPB.DrugMonitoring == nil ||
		len(antimicrobialPB.DrugMonitoring.Values) == 0:
		err = errs.MissingField("DrugMonitoring.Values")
	case antimicrobialPB.AdverseEffects == nil ||
		len(antimicrobialPB.AdverseEffects.Values) == 0:
		err = errs.MissingField("AdverseEffects.Values")
	case antimicrobialPB.MajorInteractions == nil ||
		len(antimicrobialPB.MajorInteractions.Values) == 0:
		err = errs.MissingField("MajorInteractions.Values")
	case antimicrobialPB.Pharmacology == nil ||
		len(antimicrobialPB.Pharmacology.PharmacologyInfos) == 0:
		err = errs.MissingField("Pharmacology")
	case antimicrobialPB.AdditionalInformation == nil ||
		len(antimicrobialPB.AdditionalInformation.Values) == 0:
		err = errs.MissingField("AdditionalInformation")
	case antimicrobialPB.ActivitySpectrum == nil ||
		len(antimicrobialPB.ActivitySpectrum.Spectrum) == 0:
		err = errs.MissingField("ActivitySpectrum")
	}
	if err != nil {
		return nil, err
	}

	// Get database model
	antimicrobialDB, err := getAntimicrobialDB(antimicrobialPB)
	if err != nil {
		return nil, err
	}

	// Create in database
	err = papi.sqlDB.Create(antimicrobialDB).Error
	switch {
	case err == nil:
	case strings.Contains(strings.ToLower(err.Error()), "duplicate"):
		return nil, errs.DuplicateField("antimicrobial name", antimicrobialDB.AntimicrobialName)
	default:
		return nil, errs.SQLQueryFailed(err, "CREATE")
	}

	// That's it!
	return &antimicrobial.CreateAntimicrobialResponse{
		AntimicrobialId: fmt.Sprint(int64(antimicrobialDB.ID)),
	}, nil
}

func (papi *antimicrobialAPIServer) UpdateAntimicrobial(
	ctx context.Context, updateReq *antimicrobial.UpdateAntimicrobialRequest,
) (*empty.Empty, error) {
	// Request must not be nil
	if updateReq == nil {
		return nil, errs.NilObject("UpdateAntimicrobialRequest")
	}

	// Authorize request
	_, err := papi.authAPI.AuthorizeGroup(ctx, createAllowedGroups...)
	if err != nil {
		return nil, err
	}

	// Validation
	if updateReq.GetAntimicrobialId() == "" {
		return nil, errs.MissingField("")
	}

	// Get database model
	antimicrobialDB, err := getAntimicrobialDB(updateReq.GetAntimicrobial())
	if err != nil {
		return nil, err
	}

	// Update model
	err = papi.sqlDB.Table(antimicrobialsTable).Where("id=?", updateReq.AntimicrobialId).
		Updates(antimicrobialDB).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "UPDATE")
	}

	return &empty.Empty{}, nil
}

func (papi *antimicrobialAPIServer) DeleteAntimicrobial(
	ctx context.Context, delReq *antimicrobial.DeleteAntimicrobialRequest,
) (*empty.Empty, error) {
	// Request must not be nil
	if delReq == nil {
		return nil, errs.NilObject("DeleteAntimicrobialRequest")
	}

	// Authorize request
	_, err := papi.authAPI.AuthorizeGroup(ctx, deleteAllowedGroups...)
	if err != nil {
		return nil, err
	}

	// Validation
	if delReq.GetAntimicrobialId() == "" {
		return nil, errs.MissingField("antimicrobial id")
	}

	// Delete in database
	err = papi.sqlDB.Table(antimicrobialsTable).Delete(&Antimicrobial{}, "id=?", delReq.AntimicrobialId).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "DELETE")
	}

	return &empty.Empty{}, nil
}

func (papi *antimicrobialAPIServer) ListAntimicrobials(
	ctx context.Context, listReq *antimicrobial.ListAntimicrobialsRequest,
) (*antimicrobial.Antimicrobials, error) {
	// Request must not be nil
	if listReq == nil {
		return nil, errs.NilObject("ListAntimicrobialsRequest")
	}

	// Authenticate request
	err := papi.authAPI.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Normalize page
	pageNumber, pageSize := modules.NormalizePage(listReq.PageToken, listReq.PageSize)
	offset := pageNumber*pageSize - pageSize

	antimicrobialsDB := make([]*Antimicrobial, 0, pageSize)
	err = papi.sqlDB.Order("created_at DESC").Offset(offset).Limit(pageSize).
		Find(&antimicrobialsDB).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "LIST")
	}

	antimicrobialsPB := make([]*antimicrobial.Antimicrobial, 0, len(antimicrobialsDB))
	for _, antimicrobialDB := range antimicrobialsDB {
		antimicrobialPB, err := getAntimicrobialPB(antimicrobialDB)
		if err != nil {
			return nil, err
		}
		antimicrobialsPB = append(antimicrobialsPB, getAntimicrobialView(antimicrobialPB, listReq.View))
	}

	return &antimicrobial.Antimicrobials{
		Antimicrobials: antimicrobialsPB,
	}, nil
}

func (papi *antimicrobialAPIServer) SearchAntimicrobials(
	ctx context.Context, searchReq *antimicrobial.SearchAntimicrobialsRequest,
) (*antimicrobial.Antimicrobials, error) {
	// Request must not be nil
	if searchReq == nil {
		return nil, errs.NilObject("SearchAntimicrobialsRequest")
	}

	// Authenticate request
	err := papi.authAPI.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// For empty queries
	if searchReq.Query == "" {
		return &antimicrobial.Antimicrobials{
			Antimicrobials: []*antimicrobial.Antimicrobial{},
		}, nil
	}

	pageNumber, pageSize := modules.NormalizePage(searchReq.GetPageToken(), searchReq.GetPageSize())
	offset := (pageNumber * pageSize) - pageSize

	parsedQuery := modules.ParseQuery(searchReq.Query, " antimicrobials", "antimicrobial")

	antimicrobialsDB := make([]*Antimicrobial, 0, pageSize)

	err = papi.sqlDB.Unscoped().Offset(offset).Limit(pageSize).
		Find(&antimicrobialsDB, "MATCH(antimicrobial_name) AGAINST(? IN BOOLEAN MODE)", parsedQuery).Error
	switch {
	case err == nil:
	default:
		return nil, errs.SQLQueryFailed(err, "SELECT")
	}

	// Populate response
	antimicrobialsPB := make([]*antimicrobial.Antimicrobial, 0, len(antimicrobialsDB))

	for _, antimicrobialDB := range antimicrobialsDB {
		antimicrobialPB, err := getAntimicrobialPB(antimicrobialDB)
		if err != nil {
			return nil, err
		}
		antimicrobialsPB = append(antimicrobialsPB, getAntimicrobialView(antimicrobialPB, searchReq.GetView()))
	}

	return &antimicrobial.Antimicrobials{
		NextPageToken:  int32(pageNumber + 1),
		Antimicrobials: antimicrobialsPB,
	}, nil
}

func (papi *antimicrobialAPIServer) GetAntimicrobial(
	ctx context.Context, getReq *antimicrobial.GetAntimicrobialRequest,
) (*antimicrobial.Antimicrobial, error) {
	// Request must not be nil
	if getReq == nil {
		return nil, errs.NilObject("GetAntimicrobialRequest")
	}

	// Authenticate request
	err := papi.authAPI.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Validation
	if getReq.AntimicrobialId == "" {
		return nil, errs.MissingField("antimicrobial id")
	}

	antimicrobialDB := &Antimicrobial{}

	err = papi.sqlDB.First(antimicrobialDB, "id=?", getReq.AntimicrobialId).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errs.NotFound("antimicrobial", getReq.AntimicrobialId)
	default:
		return nil, errs.SQLQueryFailed(err, "GET")
	}

	// Get antimicrobial pb
	antimicrobialPB, err := getAntimicrobialPB(antimicrobialDB)
	if err != nil {
		return nil, err
	}

	return getAntimicrobialView(antimicrobialPB, getReq.View), nil
}

func getAntimicrobialView(
	antimicrobialPB *antimicrobial.Antimicrobial,
	view antimicrobial.AntimicrobialView,
) *antimicrobial.Antimicrobial {
	antimicrobialView := &antimicrobial.Antimicrobial{}

	switch view {
	case antimicrobial.AntimicrobialView_LIST:
		// Server response include antimicrobial_id, antimicrobial_name, and general_usage
		antimicrobialView.AntimicrobialId = antimicrobialPB.AntimicrobialId
		antimicrobialView.AntimicrobialName = antimicrobialPB.AntimicrobialName
		antimicrobialView.GeneralUsage = antimicrobialPB.GeneralUsage
	default:
		antimicrobialView = antimicrobialPB
	}

	return antimicrobialView
}
