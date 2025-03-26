package core

import (
	"errors"

	"github.com/jonasroussel/hyve/stores"
)

func (app *HyveApp) GetStore() (*DBStore, error) {
	var store DBStore

	err := app.DB.First(&store).Error
	if err != nil {
		return nil, err
	}

	return &store, nil
}

func (app *HyveApp) SetStore(dbStore DBStore) error {
	if app.Store != nil {
		if err := app.Store.Close(); err != nil {
			return err
		}
	}

	var store stores.Store
	var err error

	if dbStore.Type == DBStoreType_SQL {
		store, err = stores.NewSQLStore(dbStore.Config["driver"].(string), dbStore.Config["datasource"].(string))
		if err != nil {
			return err
		}
	}

	if dbStore.Type == DBStoreType_Mongo {
		store, err = stores.NewMongoStore(dbStore.Config["uri"].(string), dbStore.Config["database"].(string))
		if err != nil {
			return err
		}
	}

	if store == nil {
		return errors.New("unknown store type")
	}

	err = store.Load()
	if err != nil {
		return err
	}

	app.Store = store

	return nil
}
