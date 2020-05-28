package account

import (
	"encoding/json"
	"github.com/gidyon/antibug/internal/pkg/errs"
	"github.com/gidyon/antibug/pkg/api/account"
	"github.com/jinzhu/gorm"
)

const accountsTable = "accounts"

// Account is model for app user
type Account struct {
	FirstName         string `gorm:"varchar(50);not null"`
	LastName          string `gorm:"varchar(50);not null"`
	Email             string `gorm:"varchar(50);not null"`
	Phone             string `gorm:"varchar(15);not null"`
	Gender            string `gorm:"varchar(8);not null"`
	Group             string `gorm:"varchar(50);not null"`
	ProfileURL        string `gorm:"varchar(256);not null"`
	Password          string `gorm:"varchar(256);not null"`
	DeviceToken       string `gorm:"varchar(256);not null"`
	Active            bool   `gorm:"tinyint(1);default:0"`
	Jobs              []byte `gorm:"type:json"`
	StarredFacilities []byte `gorm:"type:json"`
	Settings          []byte `gorm:"type:json"`
	gorm.Model
}

// TableName ...
func (*Account) TableName() string {
	return accountsTable
}

func getAccountDB(accountPB *account.Account) (*Account, error) {
	if accountPB == nil {
		return nil, errs.NilObject("AccountPB")
	}

	accountDB := &Account{
		FirstName:   accountPB.FirstName,
		LastName:    accountPB.LastName,
		Email:       accountPB.Email,
		Phone:       accountPB.Phone,
		Gender:      accountPB.Gender,
		Group:       accountPB.Group,
		ProfileURL:  accountPB.ProfileUrl,
		DeviceToken: accountPB.DeviceToken,
		Active:      accountPB.Active,
	}

	return accountDB, nil
}

func getAccountPB(accountDB *Account) (*account.Account, error) {
	if accountDB == nil {
		return nil, errs.NilObject("AccountDB")
	}

	accountPB := &account.Account{
		FirstName:   accountDB.FirstName,
		LastName:    accountDB.LastName,
		Email:       accountDB.Email,
		Phone:       accountDB.Phone,
		Gender:      accountDB.Gender,
		Group:       accountDB.Group,
		ProfileUrl:  accountDB.ProfileURL,
		DeviceToken: accountDB.DeviceToken,
		Active:      accountDB.Active,
	}

	return accountPB, nil
}

func getJobsPB(data []byte) ([]*account.Job, error) {
	jobs := make([]*account.Job, 0)
	if len(data) > 0 {
		err := json.Unmarshal(data, &jobs)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "Jobs")
		}
	}
	return jobs, nil
}

func getJobsDB(jobs []*account.Job) ([]byte, error) {
	if len(jobs) > 0 {
		data, err := json.Marshal(jobs)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "Jobs")
		}
		return data, nil
	}
	return []byte{}, nil
}

func getStarredFacilityPB(data []byte) ([]*account.Facility, error) {
	facilities := make([]*account.Facility, 0)
	if len(data) > 0 {
		err := json.Unmarshal(data, &facilities)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "StarredFacilities")
		}
	}
	return facilities, nil
}

func getStarredFacilityDB(facilities []*account.Facility) ([]byte, error) {
	if len(facilities) > 0 {
		data, err := json.Marshal(facilities)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "StarredFacilities")
		}
		return data, nil
	}
	return []byte{}, nil
}

func getSettingsPB(data []byte) (*account.Settings, error) {
	settings := &account.Settings{}
	if len(data) > 0 {
		err := json.Unmarshal(data, settings)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "Settings")
		}
	}
	return settings, nil
}

func getSettingsDB(settings *account.Settings) ([]byte, error) {
	if len(settings.GetSettings()) > 0 {
		data, err := json.Marshal(settings)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "Settings")
		}
		return data, nil
	}
	return []byte{}, nil
}
