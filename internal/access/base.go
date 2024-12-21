package access

import "gorm.io/gorm"

func Registry(db *gorm.DB) error {
	if err := registryAdminMenu(db); err != nil {
		return err
	}
	if err := registryAdminRole(db); err != nil {
		return err
	}
	if err := registryAdminRoleMenu(db); err != nil {
		return err
	}
	if err := registryAdminUser(db); err != nil {
		return err
	}
	if err := registryAdminUserRole(db); err != nil {
		return err
	}
	if err := registryLink(db); err != nil {
		return err
	}
	return nil
}
