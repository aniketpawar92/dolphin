package bolt

import (
	"dolphin/api"
	"dolphin/api/bolt/internal"

	"github.com/boltdb/bolt"
)

// SettingsService represents a service to manage application settings.
type SettingsService struct {
	store *Store
}

const (
	dbSettingsKey = "SETTINGS"
)

// Settings retrieve the settings object.
func (service *SettingsService) Settings() (*dockm.Settings, error) {
	var data []byte
	err := service.store.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(settingsBucketName))
		value := bucket.Get([]byte(dbSettingsKey))
		if value == nil {
			return dockm.ErrSettingsNotFound
		}

		data = make([]byte, len(value))
		copy(data, value)
		return nil
	})
	if err != nil {
		return nil, err
	}

	var settings dockm.Settings
	err = internal.UnmarshalSettings(data, &settings)
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

// StoreSettings persists a Settings object.
func (service *SettingsService) StoreSettings(settings *dockm.Settings) error {
	return service.store.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(settingsBucketName))

		data, err := internal.MarshalSettings(settings)
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(dbSettingsKey), data)
		if err != nil {
			return err
		}
		return nil
	})
}
