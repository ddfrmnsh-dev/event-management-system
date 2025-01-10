package repository

import (
	"event-management-system/models"
	"fmt"

	"gorm.io/gorm"
)

type RolePermissionRepository interface {
	FindAllRole() ([]models.Role, error)
	FindAllPermission() ([]models.Permission, error)
	// FindById(id int) (RolePermission, error)
	SaveRole(role models.Role) (models.Role, error)
	SavePermission(permission models.Permission) (models.Permission, error)
	// Update(rolePermission RolePermission) (RolePermission, error)
	DeleteRole(id int) (models.Role, error)
	DeletePermission(id int) (models.Permission, error)
	AddPermissionToRole(roleId int, permissionId []int) (models.Role, error)
}

type rolePermissionRepositoryImpl struct {
	db *gorm.DB
}

func NewRolePermissionRepository(db *gorm.DB) *rolePermissionRepositoryImpl {
	return &rolePermissionRepositoryImpl{db: db}
}

func (r *rolePermissionRepositoryImpl) FindAllRole() ([]models.Role, error) {
	var roles []models.Role
	res := r.db.Preload("Permissions").Find(&roles)

	if res.Error != nil {
		fmt.Println("err:", res.Error)
		return nil, res.Error
	}
	return roles, nil
}

func (r *rolePermissionRepositoryImpl) FindAllPermission() ([]models.Permission, error) {
	var permissions []models.Permission
	res := r.db.Find(&permissions)

	if res.Error != nil {
		fmt.Println("err:", res.Error)
		return nil, res.Error
	}
	return permissions, nil
}

func (r *rolePermissionRepositoryImpl) SaveRole(role models.Role) (models.Role, error) {
	res := r.db.Create(&role)

	if res.Error != nil {
		return role, res.Error
	}
	return role, nil
}

func (r *rolePermissionRepositoryImpl) SavePermission(permission models.Permission) (models.Permission, error) {
	res := r.db.Create(&permission)
	if res.Error != nil {
		return permission, res.Error
	}
	return permission, nil
}

func (r *rolePermissionRepositoryImpl) DeleteRole(id int) (models.Role, error) {
	var role models.Role
	if err := r.db.Preload("Permissions").First(&role, id).Error; err != nil {
		return models.Role{}, err
	}
	return models.Role{}, r.db.Delete(&role).Error
}
func (r *rolePermissionRepositoryImpl) DeletePermission(id int) (models.Permission, error) {
	var permission models.Permission

	err := r.db.Delete(&permission, id).Error
	if err != nil {
		return permission, err
	}

	return permission, nil
}

func (r *rolePermissionRepositoryImpl) AddPermissionToRole(roleId int, permissionIds []int) (models.Role, error) {
	var role models.Role
	if err := r.db.First(&role, roleId).Error; err != nil {
		return role, err
	}

	var permissions []models.Permission
	if err := r.db.Where("id IN ?", permissionIds).Find(&permissions).Error; err != nil {
		return role, err
	}

	if err := r.db.Model(&role).Association("Permissions").Append(&permissions); err != nil {
		return role, err
	}

	if err := r.db.Preload("Permissions").First(&role, roleId).Error; err != nil {
		return role, err
	}

	return role, nil
}

// func (r *rolePermissionRepositoryImpl) AddPermissionToRole(roleId int, permissionIds []int) (models.Role, error) {
// 	// Cari role berdasarkan ID
// 	var role models.Role
// 	if err := r.db.First(&role, roleId).Error; err != nil {
// 		return role, err
// 	}

// 	// Ambil data permissions yang sudah ada pada role
// 	var existingPermissions []models.Permission
// 	if err := r.db.Model(&role).Association("Permissions").Find(&existingPermissions); err != nil {
// 		return role, err
// 	}

// 	// Buat map untuk cek existing permissions
// 	existingPermissionMap := make(map[int]bool)
// 	for _, perm := range existingPermissions {
// 		existingPermissionMap[perm.ID] = true
// 	}

// 	// Filter permission IDs untuk hanya yang belum ada
// 	var newPermissionIds []int
// 	for _, id := range permissionIds {
// 		if !existingPermissionMap[id] {
// 			newPermissionIds = append(newPermissionIds, id)
// 		}
// 	}

// 	// Jika tidak ada permission baru, kembalikan role tanpa perubahan
// 	if len(newPermissionIds) == 0 {
// 		return role, nil
// 	}

// 	// Ambil data permission berdasarkan ID baru
// 	var newPermissions []models.Permission
// 	if err := r.db.Where("id IN ?", newPermissionIds).Find(&newPermissions).Error; err != nil {
// 		return role, err
// 	}

// 	// Tambahkan permission baru ke role
// 	if err := r.db.Model(&role).Association("Permissions").Append(&newPermissions); err != nil {
// 		return role, err
// 	}

// 	// Kembalikan role dengan data terbaru
// 	if err := r.db.Preload("Permissions").First(&role, roleId).Error; err != nil {
// 		return role, err
// 	}

// 	return role, nil
// }
