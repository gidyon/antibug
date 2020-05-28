package pathogen

import (
	"context"
	"errors"
	"fmt"
	"github.com/gidyon/antibug/internal/modules"
	"github.com/gidyon/antibug/internal/pkg/errs"
	"github.com/gidyon/antibug/pkg/api/pathogen"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/grpclog"
	"strings"
)

type pathogenAPIServer struct {
	sqlDB  *gorm.DB
	logger grpclog.LoggerV2
}

// Options contains parameters for NewPathogenAPI
type Options struct {
	SQLDB  *gorm.DB
	Logger grpclog.LoggerV2
}

// NewPathogenAPI creates a new pathogen API server
func NewPathogenAPI(ctx context.Context, opt *Options) (pathogen.PathogenAPIServer, error) {
	// Validation
	var err error
	switch {
	case opt.SQLDB == nil:
		err = errs.NilObject("SqlDB")
	case opt.Logger == nil:
		err = errs.NilObject("Logger")
	case ctx == nil:
		err = errs.NilObject("Context")
	}
	if err != nil {
		return nil, err
	}

	papi := &pathogenAPIServer{
		sqlDB:  opt.SQLDB,
		logger: opt.Logger,
	}

	// Perform auto migration
	err = papi.sqlDB.AutoMigrate(&Pathogen{}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to automigrate pathogens table: %w", err)
	}

	// Create a full text search index
	err = modules.CreateFullTextIndex(papi.sqlDB, pathogensTable, "pathogen_name")
	if err != nil {
		return nil, fmt.Errorf("failed to create full text index: %v", err)
	}

	return papi, nil
}

func (papi *pathogenAPIServer) CreatePathogen(
	ctx context.Context, createReq *pathogen.CreatePathogenRequest,
) (*pathogen.CreatePathogenResponse, error) {
	// Request must not be nil
	if createReq == nil {
		return nil, errs.NilObject("CreatePathogenRequest")
	}

	pathogenPB := createReq.GetPathogen()

	// Validation
	err := func() error {
		var err error
		switch {
		case strings.TrimSpace(pathogenPB.PathogenName) == "":
			err = errs.MissingField("Pathogen Name")
		case strings.TrimSpace(pathogenPB.Category) == "":
			err = errs.MissingField("Pathogen Category")
		case strings.TrimSpace(pathogenPB.GeneralInformation) == "":
			err = errs.MissingField("Pathogen GeneralInformation")
		case pathogenPB.GetEpidemology() == nil || len(pathogenPB.Epidemology.Values) == 0:
			err = errs.MissingField("Pathogen Epidemology")
		case pathogenPB.GetSymptoms() == nil || len(pathogenPB.Symptoms.Values) == 0:
			err = errs.MissingField("Pathogen Symptoms")
		case pathogenPB.GetAdditionalInfo() == nil || len(pathogenPB.AdditionalInfo.Values) == 0:
			err = errs.MissingField("Pathogen AdditionalInfo")
		case pathogenPB.GetGeneralSusceptibilities() == nil || len(pathogenPB.GeneralSusceptibilities.Susceptibilities) == 0:
			err = errs.MissingField("Pathogen GeneralSusceptibilities")
		}
		return err
	}()
	if err != nil {
		return nil, err
	}

	// Get database model
	pathogenDB, err := getPathogenDB(createReq.GetPathogen())
	if err != nil {
		return nil, err
	}

	// Create in database
	err = papi.sqlDB.Create(pathogenDB).Error
	switch {
	case err == nil:
	case strings.Contains(strings.ToLower(err.Error()), "duplicate"):
		return nil, errs.DuplicateField("pathogen name", pathogenDB.PathogenName)
	default:
		return nil, errs.SQLQueryFailed(err, "CREATE")
	}

	// That's it!
	return &pathogen.CreatePathogenResponse{
		PathogenId: fmt.Sprint(int64(pathogenDB.ID)),
	}, nil
}

func (papi *pathogenAPIServer) UpdatePathogen(
	ctx context.Context, updateReq *pathogen.UpdatePathogenRequest,
) (*empty.Empty, error) {
	// Request must not be nil
	if updateReq == nil {
		return nil, errs.NilObject("UpdatePathogenRequest")
	}

	// Validation
	if updateReq.GetPathogenId() == "" {
		return nil, errs.MissingField("")
	}

	// Get database model
	pathogenDB, err := getPathogenDB(updateReq.GetPathogen())
	if err != nil {
		return nil, err
	}

	// Update model
	err = papi.sqlDB.Table(pathogensTable).Where("id=?", updateReq.PathogenId).
		Updates(pathogenDB).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "UPDATE")
	}

	return &empty.Empty{}, nil
}

func (papi *pathogenAPIServer) DeletePathogen(
	ctx context.Context, delReq *pathogen.DeletePathogenRequest,
) (*empty.Empty, error) {
	// Request must not be nil
	if delReq == nil {
		return nil, errs.NilObject("DeletePathogenRequest")
	}

	// Validation
	if delReq.GetPathogenId() == "" {
		return nil, errs.MissingField("pathogen id")
	}

	// Delete in database
	err := papi.sqlDB.Table(pathogensTable).Delete(&Pathogen{}, "id=?", delReq.PathogenId).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "DELETE")
	}

	return &empty.Empty{}, nil
}

func (papi *pathogenAPIServer) ListPathogens(
	ctx context.Context, listReq *pathogen.ListPathogensRequest,
) (*pathogen.Pathogens, error) {
	// Request must not be nil
	if listReq == nil {
		return nil, errs.NilObject("ListPathogensRequest")
	}

	// Normalize page
	pageNumber, pageSize := modules.NormalizePage(listReq.PageToken, listReq.PageSize)
	offset := pageNumber*pageSize - pageSize

	pathogensDB := make([]*Pathogen, 0, pageSize)
	err := papi.sqlDB.Order("created_at DESC").Offset(offset).Limit(pageSize).
		Find(&pathogensDB).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "LIST")
	}

	pathogensPB := make([]*pathogen.Pathogen, 0, len(pathogensDB))
	for _, pathogenDB := range pathogensDB {
		pathogenPB, err := getPathogenPB(pathogenDB)
		if err != nil {
			return nil, err
		}
		pathogensPB = append(pathogensPB, getPathogenView(pathogenPB, listReq.View))
	}

	return &pathogen.Pathogens{
		Pathogens: pathogensPB,
	}, nil
}

func (papi *pathogenAPIServer) SearchPathogens(
	ctx context.Context, searchReq *pathogen.SearchPathogensRequest,
) (*pathogen.Pathogens, error) {
	// Request must not be nil
	if searchReq == nil {
		return nil, errs.NilObject("SearchPathogensRequest")
	}

	// For empty queries
	if searchReq.Query == "" {
		return &pathogen.Pathogens{
			Pathogens: []*pathogen.Pathogen{},
		}, nil
	}

	pageNumber, pageSize := modules.NormalizePage(searchReq.GetPageToken(), searchReq.GetPageSize())
	offset := (pageNumber * pageSize) - pageSize

	parsedQuery := modules.ParseQuery(searchReq.Query, " pathogens", "pathogen")

	pathogensDB := make([]*Pathogen, 0, pageSize)

	err := papi.sqlDB.Unscoped().Offset(offset).Limit(pageSize).
		Find(&pathogensDB, "MATCH(pathogen_name) AGAINST(? IN BOOLEAN MODE)", parsedQuery).Error
	switch {
	case err == nil:
	default:
		return nil, errs.SQLQueryFailed(err, "SELECT")
	}

	// Populate response
	pathogensPB := make([]*pathogen.Pathogen, 0, len(pathogensDB))

	for _, pathogenDB := range pathogensDB {
		pathogenPB, err := getPathogenPB(pathogenDB)
		if err != nil {
			return nil, err
		}
		pathogensPB = append(pathogensPB, getPathogenView(pathogenPB, searchReq.GetView()))
	}

	return &pathogen.Pathogens{
		NextPageToken: int32(pageNumber + 1),
		Pathogens:     pathogensPB,
	}, nil
}

func (papi *pathogenAPIServer) GetPathogen(
	ctx context.Context, getReq *pathogen.GetPathogenRequest,
) (*pathogen.Pathogen, error) {
	// Request must not be nil
	if getReq == nil {
		return nil, errs.NilObject("GetPathogenRequest")
	}

	// Validation
	if getReq.PathogenId == "" {
		return nil, errs.MissingField("pathogen id")
	}

	pathogenDB := &Pathogen{}

	err := papi.sqlDB.First(pathogenDB, "id=?", getReq.PathogenId).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errs.NotFound("pathogen", getReq.PathogenId)
	default:
		return nil, errs.SQLQueryFailed(err, "GET")
	}

	// Get pathogen pb
	pathogenPB, err := getPathogenPB(pathogenDB)

	return getPathogenView(pathogenPB, getReq.View), nil
}

func getPathogenView(
	pathogenPB *pathogen.Pathogen,
	view pathogen.PathogenView,
) *pathogen.Pathogen {
	pathogenView := &pathogen.Pathogen{}

	switch view {
	case pathogen.PathogenView_LIST:
		// Server response include pathogen_id, pathogen_name, and general_usage
		pathogenView.PathogenId = pathogenPB.PathogenId
		pathogenView.PathogenName = pathogenPB.PathogenName
		pathogenView.Category = pathogenPB.Category
		pathogenView.GeneralInformation = pathogenPB.GeneralInformation
	default:
		pathogenView = pathogenPB
	}

	return pathogenView
}
