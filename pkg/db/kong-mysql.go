package db

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"sync"
)

var Gorm sync.Map

var (
	// ErrRecordNotFound record not found error, happens when haven't find any matched data when looking up with a struct
	ErrRecordNotFound = errors.New("record not found")
	// ErrInvalidSQL invalid SQL error, happens when you passed invalid SQL
	ErrInvalidSQL = errors.New("invalid SQL")
	// ErrInvalidTransaction invalid transaction when you are trying to `Commit` or `Rollback`
	ErrInvalidTransaction = errors.New("no valid transaction")
	// ErrCantStartTransaction can't start transaction when you are trying to start one with `Begin`
	ErrCantStartTransaction = errors.New("can't start transaction")
	// ErrUnaddressable unaddressable value
	ErrUnaddressable = errors.New("using unaddressable value")
)


func GetOrm(DBConnect, DBType string, DBMaxIdle, DBMaxOpen int) (orm *gorm.DB, err error) {

	if DBConnect == "" {
		return nil, errors.New("DBConnect Is Nil")
	}

	v, ok := Gorm.Load(DBConnect)
	if ok {
		orm = v.(*gorm.DB)
	} else {
		orm, err = gorm.Open(DBType, DBConnect)
		if err != nil {
			return nil, err
		}
		if DBMaxIdle != 0 {
			orm.DB().SetMaxIdleConns(DBMaxIdle)
		}
		if DBMaxOpen != 0 {
			orm.DB().SetMaxOpenConns(DBMaxOpen)
		}

		Gorm.LoadOrStore(DBConnect, orm)
	}

	return orm, err
}


func GetSecretByAppId(orm *gorm.DB, SecretSQL, appId string) (secret string, err error) {
	rows, err := orm.Raw(SecretSQL, appId).Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&secret)
		if err != nil {
			return "", err
		}
		return secret, err
	}
	return "", nil
}
