package core

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func OpenDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("./hyve.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&DBAccount{}, &DBStore{}, &DBSession{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

//----------//
// Accounts //
//----------//

type DBRole int8

const (
	DBRole_SuperAdmin DBRole = 0
	DBRole_Admin      DBRole = 1
)

type DBAccount struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     *string   `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	Role      DBRole    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (DBAccount) TableName() string {
	return "hyve_accounts"
}

//--------//
// Stores //
//--------//

type DBStoreType string

const (
	DBStoreType_SQL   DBStoreType = "sql"
	DBStoreType_Mongo DBStoreType = "mongo"
)

type DBStore struct {
	ID        uint           `json:"id"`
	Type      DBStoreType    `json:"type"`
	Config    map[string]any `json:"config" gorm:"serializer:json"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (DBStore) TableName() string {
	return "hyve_stores"
}

//----------//
// Sessions //
//----------//

type DBSession struct {
	ID        string    `json:"id"`
	AccountID uint      `json:"account_id"`
	Account   DBAccount `json:"account"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (DBSession) TableName() string {
	return "hyve_sessions"
}
