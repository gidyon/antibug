package account

import (
	"context"
	"errors"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/antibug/internal/pkg/auth"
	"github.com/gidyon/antibug/internal/pkg/errs"
	"github.com/gidyon/antibug/pkg/api/account"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"strings"
)

type accountAPIServer struct {
	sqlDB   *gorm.DB
	logger  grpclog.LoggerV2
	authAPI auth.Interface
}

// Options contains parameters for passing to NewAccountAPI
type Options struct {
	SQLDB      *gorm.DB
	Logger     grpclog.LoggerV2
	SigningKey string
}

// NewAccountAPI is factory for creating account APIs
func NewAccountAPI(ctx context.Context, opt *Options) (account.AccountAPIServer, error) {
	var err error
	switch {
	case ctx == nil:
		err = errs.NilObject("Context")
	case opt.SQLDB == nil:
		err = errs.NilObject("SqlDB")
	case opt.Logger == nil:
		err = errs.NilObject("Logger")
	case opt.SigningKey == "":
		err = errs.MissingField("signing key")
	}
	if err != nil {
		return nil, err
	}

	authAPI, err := auth.NewAPI(randomdata.RandStringRunes(32))
	if err != nil {
		return nil, err
	}

	api := &accountAPIServer{
		sqlDB:   opt.SQLDB,
		logger:  opt.Logger,
		authAPI: authAPI,
	}

	// Perform automigration
	err = api.sqlDB.AutoMigrate(&Account{}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to automigrate table: %v", err)
	}

	return api, nil
}

func (api *accountAPIServer) Login(
	ctx context.Context, loginReq *account.LoginRequest,
) (*account.LoginResponse, error) {
	// Request must not be nil
	if loginReq == nil {
		return nil, errs.NilObject("LoginRequest")
	}

	// Validation
	var err error
	switch {
	case loginReq.Username == "":
		err = errs.MissingField("username")
	case loginReq.Password == "":
		err = errs.MissingField("password")
	}
	if err != nil {
		return nil, err
	}

	// Query model
	accountDB := &Account{}
	err = api.sqlDB.First(accountDB, "email=? OR phone=?", loginReq.Username, loginReq.Username).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errs.AccountNotFound(loginReq.Username)
	default:
		return nil, errs.SQLQueryFailed(err, "LOGIN")
	}

	// Check password
	if accountDB.Password == "" {
		return nil, errs.WrapMessage(
			codes.PermissionDenied, "account has no password; please request new password",
		)
	}
	accountPB, err := getAccountPB(accountDB)
	if err != nil {
		return nil, err
	}

	// Check that account is not blocked
	if !accountPB.Active {
		return nil, errs.WrapMessage(
			codes.PermissionDenied, "account is not active; please activate account first",
		)
	}

	// Compare passwords
	err = compareHash(accountDB.Password, loginReq.Password)
	if err != nil {
		return nil, errs.WrapMessage(
			codes.Unauthenticated, "wrong password",
		)
	}

	accountID := fmt.Sprint(accountDB.ID)

	// Generate token
	token, err := api.authAPI.GenToken(ctx, &auth.Payload{
		ID:        accountID,
		FirstName: accountDB.FirstName,
		LastName:  accountDB.LastName,
		Group:     accountDB.Group,
	}, 0)
	if err != nil {
		return nil, errs.FailedToGenToken(err)
	}

	// Populate response
	return &account.LoginResponse{
		Token:        token,
		AccountId:    accountID,
		AccountState: accountDB.Active,
		AccountGroup: accountDB.Group,
	}, nil
}

// generates hashed version of password
func genHash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// compares hashed password with password
func compareHash(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (api *accountAPIServer) CreateAccount(
	ctx context.Context, createReq *account.CreateAccountRequest,
) (*account.CreateAccountResponse, error) {
	// Request must not be nil
	if createReq == nil {
		return nil, errs.NilObject("Account")
	}

	// Validation
	var err error
	accountPB := createReq.GetAccount()
	switch {
	case accountPB == nil:
		err = errs.NilObject("Account")
	case accountPB.FirstName == "":
		err = errs.MissingField("FirstName")
	case accountPB.LastName == "":
		err = errs.MissingField("LastName")
	case accountPB.Phone == "":
		err = errs.MissingField("Phone")
	case accountPB.Email == "":
		err = errs.MissingField("Email")
	case accountPB.Gender == "":
		err = errs.MissingField("Gender")
	}
	if err != nil {
		return nil, err
	}

	// Get model
	accountDB, err := getAccountDB(accountPB)
	if err != nil {
		return nil, err
	}

	// Check passord if not empty
	if createReq.Password != "" {
		// Passowrds must match
		if createReq.Password != createReq.ConfirmPassword {
			return nil, errs.WrapMessage(
				codes.InvalidArgument, "passwords do not match",
			)
		}

		// Hash password
		hashedPass, err := genHash(createReq.Password)
		if err != nil {
			return nil, errs.FailedToGenHashedPass(err)
		}

		accountDB.Password = hashedPass
	}

	accountDB.Active = false

	// Create in database
	err = api.sqlDB.Create(accountDB).Error
	switch {
	case err == nil:
	default:
		errStr := err.Error()
		switch {
		case strings.Contains(errStr, "email"):
			return nil, errs.DuplicateField("email", accountPB.Email)
		case strings.Contains(errStr, "phone"):
			return nil, errs.DuplicateField("phone", accountPB.Phone)
		default:
			return nil, errs.SQLQueryFailed(err, "CREATE")
		}
	}

	return &account.CreateAccountResponse{
		AccountId: fmt.Sprint(accountDB.ID),
	}, nil
}

func (api *accountAPIServer) ActivateAccount(
	ctx context.Context, activateReq *account.ActivateAccountRequest,
) (*empty.Empty, error) {
	// Request must not be nil
	if activateReq == nil {
		return nil, errs.NilObject("ActivateAccountRequest")
	}

	// Authorize request
	_, err := api.authAPI.AuthorizeActor(ctx, activateReq.AccountId)
	if err != nil {
		return nil, err
	}

	// Validation
	switch {
	case activateReq.AccountId == "":
		err = errs.MissingField("account id")
	}
	if err != nil {
		return nil, err
	}

	// Update account state in database
	err = api.sqlDB.Table(accountsTable).Where("id=?", activateReq.AccountId).Update("active", true).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "UPDATE")
	}

	return &empty.Empty{}, nil
}

func (api *accountAPIServer) UpdateAccount(
	ctx context.Context, updateReq *account.UpdateAccountRequest,
) (*empty.Empty, error) {
	// Request must not be nil
	if updateReq == nil {
		return nil, errs.NilObject("UpdateAccountRequest")
	}

	// Authorize request
	_, err := api.authAPI.AuthorizeActor(ctx, updateReq.GetAccountId())
	if err != nil {
		return nil, err
	}

	// Validation
	switch {
	case updateReq.AccountId == "":
		err = errs.MissingField("account id")
	case updateReq.Account == nil:
		err = errs.NilObject("account")
	}
	if err != nil {
		return nil, err
	}

	// Get model
	accountDB, err := getAccountDB(updateReq.Account)
	if err != nil {
		return nil, err
	}

	// Save in model
	err = api.sqlDB.Table(accountsTable).Where("id=?", updateReq.AccountId).
		Updates(accountDB).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "UPDATE")
	}

	return &empty.Empty{}, nil
}

func (api *accountAPIServer) GetAccount(
	ctx context.Context, getReq *account.GetRequest,
) (*account.Account, error) {
	// Request must not be nil
	if getReq == nil {
		return nil, errs.NilObject("GetRequest")
	}

	// Authorize request
	_, err := api.authAPI.AuthorizeActor(ctx, getReq.GetAccountId())
	if err != nil {
		return nil, err
	}

	// Validation
	if getReq.AccountId == "" {
		return nil, errs.MissingField("account id")
	}

	// Get from model
	accountDB := &Account{}

	err = api.sqlDB.First(accountDB, "id=?", getReq.AccountId).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errs.NotFound("account", getReq.AccountId)
	default:
		return nil, errs.SQLQueryFailed(err, "GET")
	}

	accountPB, err := getAccountPB(accountDB)
	if err != nil {
		return nil, err
	}

	return accountPB, nil
}

func (api *accountAPIServer) GetSettings(
	ctx context.Context, getReq *account.GetRequest,
) (*account.Settings, error) {
	//  Request must not be nil
	if getReq == nil {
		return nil, errs.NilObject("GetRequest")
	}

	// Authorize request
	_, err := api.authAPI.AuthorizeActor(ctx, getReq.GetAccountId())
	if err != nil {
		return nil, err
	}

	// Validation
	if getReq.AccountId == "" {
		return nil, errs.MissingField("account id")
	}

	// Get from database
	data := make([]byte, 0)
	err = api.sqlDB.Table(accountsTable).Where("id=?", getReq.AccountId).Select("settings").
		Row().Scan(&data)
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errs.NotFound("account", getReq.AccountId)
	default:
		return nil, errs.SQLQueryFailed(err, "GET")
	}

	settings, err := getSettingsPB(data)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

func (api *accountAPIServer) UpdateSettings(
	ctx context.Context, updateReq *account.UpdateSettingsRequest,
) (*empty.Empty, error) {
	// Request must nt be nil
	if updateReq == nil {
		return nil, errs.NilObject("UpdateSettingsRequest")
	}

	// Authorize request
	_, err := api.authAPI.AuthorizeActor(ctx, updateReq.GetAccountId())
	if err != nil {
		return nil, err
	}

	// Validation
	switch {
	case updateReq.Settings == nil:
		err = errs.NilObject("Settings")
	case updateReq.AccountId == "":
		err = errs.MissingField("AccountId")
	}
	if err != nil {
		return nil, err
	}

	// Marshal settings
	data, err := getSettingsDB(updateReq.Settings)
	if err != nil {
		return nil, err
	}

	// Update model
	err = api.sqlDB.Table(accountsTable).Where("id=?", updateReq.AccountId).
		Update("settings", data).Error
	switch {
	case err == nil:
	default:
		return nil, errs.SQLQueryFailed(err, "UPDATE")
	}

	return &empty.Empty{}, nil
}

func (api *accountAPIServer) GetJobs(
	ctx context.Context, getReq *account.GetRequest,
) (*account.Jobs, error) {
	//  Request must not be nil
	if getReq == nil {
		return nil, errs.NilObject("GetRequest")
	}

	// Authorize request
	_, err := api.authAPI.AuthorizeActor(ctx, getReq.GetAccountId())
	if err != nil {
		return nil, err
	}

	// Validation
	if getReq.AccountId == "" {
		return nil, errs.MissingField("account id")
	}

	// Get from database
	data := make([]byte, 0)
	err = api.sqlDB.Table(accountsTable).Where("id=?", getReq.AccountId).Select("jobs").
		Row().Scan(&data)
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errs.NotFound("account", getReq.AccountId)
	default:
		return nil, errs.SQLQueryFailed(err, "GET")
	}

	jobs, err := getJobsPB(data)
	if err != nil {
		return nil, err
	}

	return &account.Jobs{
		Jobs: jobs,
	}, nil
}

func (api *accountAPIServer) UpdateJobs(
	ctx context.Context, updateReq *account.UpdateJobsRequest,
) (*empty.Empty, error) {
	// Request must nt be nil
	if updateReq == nil {
		return nil, errs.NilObject("UpdateJobsRequest")
	}

	// Authorize request
	_, err := api.authAPI.AuthorizeActor(ctx, updateReq.GetAccountId())
	if err != nil {
		return nil, err
	}

	// Validation
	switch {
	case updateReq.Jobs == nil:
		err = errs.NilObject("Jobs")
	case updateReq.AccountId == "":
		err = errs.MissingField("AccountId")
	}
	if err != nil {
		return nil, err
	}

	// Validate passed jobs
	for index, job := range updateReq.Jobs {
		switch {
		case job.GetFacilityId() == "":
			err = errs.MissingField(fmt.Sprintf("facility id at index %d", index))
		case job.GetFacilityName() == "":
			err = errs.MissingField(fmt.Sprintf("facility name at index %d", index))
		case job.GetRole() == "":
			err = errs.MissingField(fmt.Sprintf("job role at index %d", index))
		}
		if err != nil {
			return nil, err
		}
		job.JobId = uuid.New().String()
	}

	// Marshal settings
	data, err := getJobsDB(updateReq.Jobs)
	if err != nil {
		return nil, err
	}

	// Update model
	err = api.sqlDB.Table(accountsTable).Where("id=?", updateReq.AccountId).
		Update("jobs", data).Error
	switch {
	case err == nil:
	default:
		return nil, errs.SQLQueryFailed(err, "UPDATE")
	}

	return &empty.Empty{}, nil
}

func (api *accountAPIServer) GetStarredFacilities(
	ctx context.Context, getReq *account.GetRequest,
) (*account.StarredFacilities, error) {
	//  Request must not be nil
	if getReq == nil {
		return nil, errs.NilObject("GetRequest")
	}

	// Authorize request
	_, err := api.authAPI.AuthorizeActor(ctx, getReq.GetAccountId())
	if err != nil {
		return nil, err
	}

	// Validation
	if getReq.AccountId == "" {
		return nil, errs.MissingField("account id")
	}

	// Get from database
	data := make([]byte, 0)
	err = api.sqlDB.Table(accountsTable).Where("id=?", getReq.AccountId).Select("starred_facilities").
		Row().Scan(&data)
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errs.NotFound("account", getReq.AccountId)
	default:
		return nil, errs.SQLQueryFailed(err, "GET")
	}

	starredFacilities, err := getStarredFacilityPB(data)
	if err != nil {
		return nil, err
	}

	return &account.StarredFacilities{
		Facilities: starredFacilities,
	}, nil
}

func (api *accountAPIServer) UpdateStarredFacilities(
	ctx context.Context, updateReq *account.UpdateStarredFacilitiesRequest,
) (*empty.Empty, error) {
	// Request must nt be nil
	if updateReq == nil {
		return nil, errs.NilObject("UpdateStarredFacilitiesRequest")
	}

	// Authorize request
	_, err := api.authAPI.AuthorizeActor(ctx, updateReq.GetAccountId())
	if err != nil {
		return nil, err
	}

	// Validation
	switch {
	case len(updateReq.Facilities) == 0:
		err = errs.NilObject("Facilities")
	case updateReq.AccountId == "":
		err = errs.MissingField("AccountId")
	}
	if err != nil {
		return nil, err
	}

	// Validate passed facilities
	for index, facilty := range updateReq.Facilities {
		switch {
		case facilty.GetFacilityId() == "":
			err = errs.MissingField(fmt.Sprintf("facility id at index %d", index))
		case facilty.GetFacilityId() == "":
			err = errs.MissingField(fmt.Sprintf("facility name at index %d", index))
		}
		if err != nil {
			return nil, err
		}
	}

	// Marshal settings
	data, err := getStarredFacilityDB(updateReq.Facilities)
	if err != nil {
		return nil, err
	}

	// Update model
	err = api.sqlDB.Table(accountsTable).Where("id=?", updateReq.AccountId).
		Update("starred_facilities", data).Error
	switch {
	case err == nil:
	default:
		return nil, errs.SQLQueryFailed(err, "UPDATE")
	}

	return &empty.Empty{}, nil
}
