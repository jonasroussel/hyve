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

	err = conn.Ping()
	if err != nil {
		return err
	}

	store.db = conn

	err = createSQLTable(store.db)
	if err != nil {
		return err
	}

	return nil
}

func (store SQLStore) AddCertificate(domain string, cert Certificate) error {
	if existsInSQLTable(domain, store.db) {
		return nil
	}

	query := "INSERT INTO hyve_certificates (domain, certificate, private_key, issuer, expires_at, created_at) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := store.db.Exec(query, domain, cert.CertificateData, cert.PrivateKeyData, cert.Issuer, cert.ExpiresAt, cert.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (store SQLStore) GetCertificate(domain string) (*Certificate, error) {
	query := "SELECT domain, certificate, private_key, issuer, expires_at, created_at FROM hyve_certificates WHERE domain = ?"
	row := store.db.QueryRow(query, domain)

	var cert Certificate
	err := row.Scan(&cert.Domain, &cert.CertificateData, &cert.PrivateKeyData, &cert.Issuer, &cert.ExpiresAt, &cert.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &cert, nil
}

func (store SQLStore) GetAllCertificates(exp int64) []Certificate {
	query := "SELECT domain, certificate, private_key, issuer, expires_at, created_at FROM hyve_certificates WHERE expires_at <= ?"
	rows, err := store.db.Query(query, exp)
	if err != nil {
		return nil
	}

	var certs []Certificate
	for rows.Next() {
		var cert Certificate
		err = rows.Scan(&cert.Domain, &cert.CertificateData, &cert.PrivateKeyData, &cert.Issuer, &cert.ExpiresAt, &cert.CreatedAt)
		if err != nil {
			return nil
		}

		certs = append(certs, cert)
	}

	return certs
}

func (store SQLStore) UpdateCertificate(domain string, cert Certificate) error {
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

	values = append(values, domain)

	query := fmt.Sprintf("UPDATE hyve_certificates SET %s WHERE domain = ?", strings.Join(fields, ", "))
	_, err := store.db.Exec(query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (store SQLStore) RemoveCertificate(domain string) error {
	query := "DELETE FROM hyve_certificates WHERE domain = ?"
	_, err := store.db.Exec(query, domain)
	if err != nil {
		return err
	}

	return nil
}

//---------//
// Helpers //
//---------//

func existsInSQLTable(domain string, db *sql.DB) bool {
	query := "SELECT domain FROM hyve_certificates WHERE domain = ?"

	var result sql.NullString
	err := db.QueryRow(query, domain).Scan(&result)

	if err == sql.ErrNoRows {
		return false
	}

	return result.Valid
}

func createSQLTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS hyve_certificates (
		domain VARCHAR(255) PRIMARY KEY,
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
