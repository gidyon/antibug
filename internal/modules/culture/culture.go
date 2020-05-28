package culture

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gidyon/antibug/internal/modules"
	"github.com/gidyon/antibug/internal/pkg/errs"
	"github.com/gidyon/antibug/pkg/api/culture"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/grpclog"
	"strings"
)

type cultureAPIServer struct {
	sqlDB  *gorm.DB
	logger grpclog.LoggerV2
}

// Options contains parameters to NewCultureAPI
type Options struct {
	SQLDB  *gorm.DB
	Logger grpclog.LoggerV2
}

// NewCultureAPI is factory for creating culture APIs
func NewCultureAPI(ctx context.Context, opt *Options) (culture.CultureAPIServer, error) {
	// Validation
	var err error
	switch {
	case ctx == nil:
		err = errs.NilObject("Context")
	case opt.SQLDB == nil:
		err = errs.NilObject("SqlDB")
	case opt.Logger == nil:
		err = errs.NilObject("Logger")
	}
	if err != nil {
		return nil, err
	}

	capi := &cultureAPIServer{
		sqlDB:  opt.SQLDB,
		logger: opt.Logger,
	}

	// Perform automigration
	err = capi.sqlDB.AutoMigrate(&Culture{}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to perform automigration: %v", err)
	}

	return capi, nil
}

func (capi *cultureAPIServer) CreateCulture(
	ctx context.Context, createReq *culture.CreateCultureRequest,
) (*culture.CreateCultureResponse, error) {
	// Request must not be nil
	if createReq == nil {
		return nil, errs.NilObject("CreateCultureRequest")
	}

	// Validation
	var err error
	culturePB := createReq.GetCulture()
	switch {
	case strings.TrimSpace(culturePB.LabTechId) == "":
		err = errs.MissingField("lab tech id")
	case strings.TrimSpace(culturePB.HospitalId) == "":
		err = errs.MissingField("hospital id")
	case strings.TrimSpace(culturePB.CountyCode) == "":
		err = errs.MissingField("county code")
	case strings.TrimSpace(culturePB.SubCountyCode) == "":
		err = errs.MissingField("sub county code")
	case strings.TrimSpace(culturePB.PatientId) == "":
		err = errs.MissingField("patient id")
	case strings.TrimSpace(culturePB.PatientGender) == "":
		err = errs.MissingField("patient gender")
	case strings.TrimSpace(culturePB.PatientAge) == "":
		err = errs.MissingField("patient age")
	case strings.TrimSpace(culturePB.CultureSource) == "":
		err = errs.MissingField("culture source")
	case len(culturePB.PathogensFound) == 0:
		err = errs.MissingField("pathogens found")
	case len(culturePB.AntimicrobialsUsed) == 0:
		err = errs.MissingField("antimicrobials used")
	case len(culturePB.CultureResults) == 0:
		err = errs.MissingField("culture results")
	}
	if err != nil {
		return nil, err
	}

	culturePB.Editors = []string{culturePB.LabTechId}

	// Get culture model
	cultureDB, err := getCultureDB(culturePB)
	if err != nil {
		return nil, err
	}

	// Save to database
	err = capi.sqlDB.Create(cultureDB).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "SAVE")
	}

	return &culture.CreateCultureResponse{
		CultureId: fmt.Sprint(cultureDB.ID),
	}, nil
}

func (capi *cultureAPIServer) UpdateCulture(
	ctx context.Context, updateReq *culture.UpdateCultureRequest,
) (*empty.Empty, error) {
	// Request must not be nil
	if updateReq == nil {
		return nil, errs.NilObject("UpdateCultureRequest")
	}

	// Validation
	culturePB := updateReq.GetCulture()
	var err error
	switch {
	case culturePB == nil:
		err = errs.NilObject("culture")
	case updateReq.CultureId == "":
		err = errs.MissingField("culture id")
	case updateReq.EditorId == "":
		err = errs.MissingField("editor id")
	}
	if err != nil {
		return nil, err
	}

	cultureDB := &Culture{}

	// Check if record exist in db, we also need culture editors
	if notFound := capi.sqlDB.Select("editors").
		First(cultureDB, "id=?", updateReq.CultureId).
		RecordNotFound(); notFound {
		return nil, errs.NotFound("culture", updateReq.CultureId)
	}

	// Unmarshal the editors
	if len(cultureDB.Editors) > 0 {
		err = json.Unmarshal(cultureDB.Editors, &culturePB.Editors)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "editors")
		}
	} else {
		culturePB.Editors = make([]string, 0)
	}

	// Add the actor to list of editors
	culturePB.Editors = append(culturePB.Editors, updateReq.EditorId)

	cultureDB, err = getCultureDB(culturePB)
	if err != nil {
		return nil, err
	}

	// Update model in database
	err = capi.sqlDB.Table(culturesTable).Where("id=?", updateReq.CultureId).
		Updates(cultureDB).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "UPDATE")
	}

	return &empty.Empty{}, nil
}

func (capi *cultureAPIServer) DeleteCulture(
	ctx context.Context, delReq *culture.DeleteCultureRequest,
) (*empty.Empty, error) {
	// Request must not be nil
	if delReq == nil {
		return nil, errs.NilObject("DeleteCultureRequest")
	}

	// Validation
	if delReq.CultureId == "" {
		return nil, errs.MissingField("culture id")
	}

	// Delete in database
	err := capi.sqlDB.Delete(&Culture{}, "id=?", delReq.CultureId).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "DELETE")
	}

	return &empty.Empty{}, nil
}

func (capi *cultureAPIServer) ListCultures(
	ctx context.Context, listReq *culture.ListCulturesRequest,
) (*culture.Cultures, error) {
	// Request must not be nil
	if listReq == nil {
		return nil, errs.NilObject("ListCulturesRequest")
	}

	db := capi.sqlDB
	if listReq.Filter != nil {
		// Target filter
		switch listReq.Filter.GetListTarget() {
		case culture.ListTarget_ALL:
		case culture.ListTarget_COUNTY:
			if len(listReq.Filter.GetTargetIds()) > 0 {
				db = db.Where("county_code IN (?)", listReq.Filter.GetTargetIds())
			}
		case culture.ListTarget_SUB_COUNTY:
			if len(listReq.Filter.GetTargetIds()) > 0 {
				db = db.Where("sub_county_code IN (?)", listReq.Filter.GetTargetIds())
			}
		case culture.ListTarget_HOSPITAL:
			if len(listReq.Filter.GetTargetIds()) > 0 {
				db = db.Where("hospital_id IN (?)", listReq.Filter.GetTargetIds())
			}
		case culture.ListTarget_PATIENT:
			if len(listReq.Filter.GetTargetIds()) > 0 {
				db = db.Where("patient_id IN (?)", listReq.Filter.GetTargetIds())
			}
		case culture.ListTarget_LAB_TECHNICIAN:
			if len(listReq.Filter.GetTargetIds()) > 0 {
				db = db.Where("lab_tech_id IN (?)", listReq.Filter.GetTargetIds())
			}
		}

		// Date filter
		if listReq.Filter.GetDateFilter() != nil && listReq.Filter.GetDateFilter().GetFilter() {
			startTimestamp := listReq.Filter.DateFilter.GetStartTimestampSec()
			endTimestamp := listReq.Filter.DateFilter.GetEndTimestampSec()
			switch {
			case startTimestamp < endTimestamp:
				db = db.Where("results_timestamp_sec BETWEEN ? AND ?", startTimestamp, endTimestamp)
			case startTimestamp > endTimestamp:
				db = db.Where("results_timestamp_sec > ?", startTimestamp)
			}
		}
	}

	// Normalize page
	pageNumber, pageSize := modules.NormalizePage(listReq.PageToken, listReq.PageSize)
	offset := pageNumber*pageSize - pageSize

	culturesDB := make([]*Culture, 0, pageSize)
	err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).
		Find(&culturesDB).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "LIST")
	}

	culturesPB := make([]*culture.Culture, 0, len(culturesDB))

	for _, cultureDB := range culturesDB {
		culturePB, err := getCulturePB(cultureDB)
		if err != nil {
			return nil, err
		}
		culturesPB = append(culturesPB, culturePB)
	}

	return &culture.Cultures{
		Cultures: culturesPB,
	}, nil
}

func (capi *cultureAPIServer) GetCulture(
	ctx context.Context, getReq *culture.GetCultureRequest,
) (*culture.Culture, error) {
	// Request must not be nil
	if getReq == nil {
		return nil, errs.NilObject("GetCultureRequest")
	}

	// Validation
	if getReq.CultureId == "" {
		return nil, errs.MissingField("culture id")
	}

	// Get culture from db
	cultureDB := &Culture{}
	err := capi.sqlDB.First(cultureDB, "id=?", getReq.CultureId).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errs.NotFound("culture", getReq.CultureId)
	default:
		return nil, errs.SQLQueryFailed(err, "GET")
	}

	// Get culture from model
	culturePB, err := getCulturePB(cultureDB)
	if err != nil {
		return nil, err
	}

	return culturePB, nil
}
