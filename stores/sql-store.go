package stores

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type SQLStore struct {
	Driver     string
	DataSource string

	db *sql.DB
}

func NewSQLStore() *SQLStore {
	driver := os.Getenv("STORE_DRIVER")
	if driver == "" {
		driver = "sqlite3"
	}

	dataSource := os.Getenv("STORE_DATA_SOURCE")
	if dataSource == "" {
		log.Fatal("STORE_DATA_SOURCE environment variable must be set when using STORE=sql")
	}

	return &SQLStore{
		Driver:     driver,
		DataSource: dataSource,
	}
}

func (store *SQLStore) Load() error {
	if !slices.Contains([]string{"sqlite3", "postgres"}, store.Driver) {
		log.Fatal("SQL driver not supported")
	}

	conn, err := sql.Open(store.Driver, store.DataSource)
	if err != nil {
		return err
	}
	store.db = conn

	err = conn.Ping()
	if err != nil {
		return err
	}

	err = createTable(store.db)
	if err != nil {
		return err
	}

	return nil
}

func (store SQLStore) AddCertificate(sni string, cert Certificate) error {
	if existsInDB(sni, store.db) {
		return nil
	}

	query := "INSERT INTO hyve_certificates (sni, certificate, private_key, issuer, expires_at, created_at) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := store.db.Exec(query, sni, cert.CertificateData, cert.PrivateKeyData, cert.Issuer, cert.ExpiresAt, cert.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (store SQLStore) GetCertificate(sni string) (*Certificate, error) {
	query := "SELECT certificate, private_key, issuer, expires_at, created_at FROM hyve_certificates WHERE sni = ?"
	row := store.db.QueryRow(query, sni)

	var cert Certificate
	err := row.Scan(&cert.CertificateData, &cert.PrivateKeyData, &cert.Issuer, &cert.ExpiresAt, &cert.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &cert, nil
}

func (store SQLStore) UpdateCertificate(sni string, cert Certificate) error {
	fields := []string{}
	values := []any{}

	if cert.CertificateData != nil {
		fields = append(fields, "certificate = ?")
		values = append(values, cert.CertificateData)
	}
	if cert.PrivateKeyData != nil {
		fields = append(fields, "private_key = ?")
		values = append(values, cert.PrivateKeyData)
	}
	if cert.Issuer != "" {
		fields = append(fields, "issuer = ?")
		values = append(values, cert.Issuer)
	}
	if cert.ExpiresAt != 0 {
		fields = append(fields, "expires_at = ?")
		values = append(values, cert.ExpiresAt)
	}
	if cert.CreatedAt != 0 {
		fields = append(fields, "created_at = ?")
		values = append(values, cert.CreatedAt)
	}

	values = append(values, sni)

	query := fmt.Sprintf("UPDATE hyve_certificates SET %s WHERE sni = ?", strings.Join(fields, ", "))
	_, err := store.db.Exec(query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (store SQLStore) RemoveCertificate(sni string) error {
	query := "DELETE FROM hyve_certificates WHERE sni = ?"
	_, err := store.db.Exec(query, sni)
	if err != nil {
		return err
	}

	return nil
}

//---------//
// Helpers //
//---------//

func existsInDB(sni string, db *sql.DB) bool {
	query := "SELECT sni FROM hyve_certificates WHERE sni = ?"

	var result sql.NullString
	err := db.QueryRow(query, sni).Scan(&result)

	if err == sql.ErrNoRows {
		return false
	}

	return result.Valid
}

func createTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS hyve_certificates (
		sni VARCHAR(255) PRIMARY KEY,
		certificate TEXT,
		private_key TEXT,
		issuer VARCHAR(255),
		expires_at BIGINT,
		created_at BIGINT
	)`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
