package stores

import (
	"context"
	"log"
	"os"
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	ConnectionString string
	DatabaseName     string

	db *mongo.Database
}

func NewMongoStore() *MongoStore {
	connStr := os.Getenv("STORE_CONNECTION_URI")
	if connStr == "" {
		log.Fatal("STORE_CONNECTION_URI environment variable must be set when using STORE=mongo")
	}

	dbName := os.Getenv("STORE_DATABASE_NAME")
	if dbName == "" {
		log.Fatal("STORE_DATABASE_NAME environment variable must be set when using STORE=mongo")
	}

	return &MongoStore{
		ConnectionString: connStr,
		DatabaseName:     dbName,
	}
}

func (store *MongoStore) Load() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(store.ConnectionString))
	if err != nil {
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	store.db = client.Database(store.DatabaseName)

	err = createMongoCollection(store.db)
	if err != nil {
		return err
	}

	return nil
}

func (store MongoStore) AddCertificate(domain string, cert Certificate) error {
	if existsInMongoCollection(domain, store.db) {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := store.db.Collection("hyve_certificates").InsertOne(ctx, cert)
	if err != nil {
		return err
	}

	return nil
}

func (store MongoStore) GetCertificate(domain string) (*Certificate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var cert Certificate
	err := store.db.Collection("hyve_certificates").FindOne(ctx, bson.M{"domain": domain}).Decode(&cert)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &cert, nil
}

func (store MongoStore) GetAllCertificates(exp int64) []Certificate {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := store.db.Collection("hyve_certificates").Find(ctx, bson.M{"expires_at": bson.M{"$lte": exp}})
	if err != nil {
		return nil
	}

	var certs []Certificate

	for cursor.Next(ctx) {
		var cert Certificate
		err = cursor.Decode(&cert)
		if err != nil {
			continue
		}

		certs = append(certs, cert)
	}

	return certs
}

func (store MongoStore) UpdateCertificate(domain string, cert Certificate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	updates := bson.M{}

	if cert.CertificateData != nil {
		updates["certificate"] = cert.CertificateData
	}
	if cert.PrivateKeyData != nil {
		updates["private_key"] = cert.PrivateKeyData
	}
	if cert.Issuer != "" {
		updates["issuer"] = cert.Issuer
	}
	if cert.ExpiresAt != 0 {
		updates["expires_at"] = cert.ExpiresAt
	}
	if cert.CreatedAt != 0 {
		updates["created_at"] = cert.CreatedAt
	}

	_, err := store.db.Collection("hyve_certificates").UpdateOne(ctx, bson.M{"domain": domain}, bson.M{"$set": updates})
	if err != nil {
		return err
	}

	return nil
}

func (store MongoStore) RemoveCertificate(domain string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := store.db.Collection("hyve_certificates").DeleteOne(ctx, bson.M{"domain": domain})
	if err != nil {
		return err
	}

	return nil
}

//---------//
// Helpers //
//---------//

func existsInMongoCollection(domain string, db *mongo.Database) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.Collection("hyve_certificates")

	count, err := collection.CountDocuments(ctx, bson.M{"domain": domain})
	if err != nil {
		return false
	}

	return count > 0
}

func createMongoCollection(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	colls, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil || slices.Contains(colls, "hyve_certificates") {
		return err
	}

	err = db.CreateCollection(ctx, "hyve_certificates")
	if err != nil {
		return err
	}

	opts := options.IndexOptions{}
	opts.SetUnique(true)

	_, err = db.Collection("hyve_certificates").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"domain": 1},
		Options: &opts,
	})
	if err != nil {
		return err
	}

	return nil
}
