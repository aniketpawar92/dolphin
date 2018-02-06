package bolt

import "dolphin/api"

// Migrator defines a service to migrate data after a DockM version update.
type Migrator struct {
	UserService            *UserService
	EndpointService        *EndpointService
	ResourceControlService *ResourceControlService
    SettingsService        *SettingsService
	VersionService         *VersionService
	CurrentDBVersion       int
	store                  *Store
}

// NewMigrator creates a new Migrator.
func NewMigrator(store *Store, version int) *Migrator {
	return &Migrator{
		UserService:            store.UserService,
		EndpointService:        store.EndpointService,
		ResourceControlService: store.ResourceControlService,
        SettingsService:        store.SettingsService,
		VersionService:         store.VersionService,
		CurrentDBVersion:       version,
		store:                  store,
	}
}

// Migrate checks the database version and migrate the existing data to the most recent data model.
func (m *Migrator) Migrate() error {

	// DockM < 1.12
	if m.CurrentDBVersion < 1  {
		err := m.updateAdminUserToDBVersion1()
		if err != nil {
			return err
		}
	}

	// DockM 1.12.x
	if m.CurrentDBVersion < 2 {
		err := m.updateResourceControlsToDBVersion2()
		if err != nil {
			return err
		}
		err = m.updateEndpointsToDBVersion2()
		if err != nil {
			return err
		}
	}

    // dockm 1.13.x
	if m.CurrentDBVersion < 3 {
		err := m.updateSettingsToDBVersion3()
		if err != nil {
			return err
		}
	}

	// dockm 1.14.0
	if m.CurrentDBVersion < 4 {
		err := m.updateEndpointsToDBVersion4()
		if err != nil {
		return err
		}
	}
 
	err := m.VersionService.StoreDBVersion(dockm.DBVersion)
	if err != nil {
		return err
	}
	return nil
}
