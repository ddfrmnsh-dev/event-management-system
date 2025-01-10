package usecase

import (
	"event-management-system/models"
	"event-management-system/repository"
	"fmt"
)

type RbacUseCase interface {
	// CheckPermission(userID uint, action string, resource string) (bool, error)
	CreateRole(role models.Role) (models.Role, error)
	CreatePermission(permission models.Permission) (models.Permission, error)
	FindAllRole() ([]models.Role, error)
	FindAllPermission() ([]models.Permission, error)
	DeleteRole(id int) (models.Role, error)
	DeletePermission(id int) (models.Permission, error)
	AddPermissionToRole(roleId int, permissionId []int) (models.Role, error)
}

type rbacUseCaseImpl struct {
	rbacRepository repository.RolePermissionRepository
}

func NewRbacUseCase(rbacRepository repository.RolePermissionRepository) RbacUseCase {
	return &rbacUseCaseImpl{rbacRepository: rbacRepository}
}

func (r *rbacUseCaseImpl) FindAllRole() ([]models.Role, error) {
	role, err := r.rbacRepository.FindAllRole()
	if err != nil {
		fmt.Println("err:", err.Error())
		return role, err
	}

	return role, nil
}
func (r *rbacUseCaseImpl) FindAllPermission() ([]models.Permission, error) {
	permission, err := r.rbacRepository.FindAllPermission()
	if err != nil {
		fmt.Println("err:", err.Error())
		return permission, err
	}

	return permission, nil
}

func (r *rbacUseCaseImpl) CreateRole(role models.Role) (models.Role, error) {
	role, err := r.rbacRepository.SaveRole(role)
	if err != nil {
		return role, err
	}

	return role, nil
}

func (r *rbacUseCaseImpl) CreatePermission(permission models.Permission) (models.Permission, error) {
	permission, err := r.rbacRepository.SavePermission(permission)
	if err != nil {
		return permission, err
	}

	return permission, nil
}

func (r *rbacUseCaseImpl) DeleteRole(id int) (models.Role, error) {
	role, err := r.rbacRepository.DeleteRole(id)
	if err != nil {
		return role, err
	}

	return role, nil
}
func (r *rbacUseCaseImpl) DeletePermission(id int) (models.Permission, error) {
	permission, err := r.rbacRepository.DeletePermission(id)
	if err != nil {
		return permission, err
	}

	return permission, nil
}

func (r *rbacUseCaseImpl) AddPermissionToRole(roleId int, permissionId []int) (models.Role, error) {
	role, err := r.rbacRepository.AddPermissionToRole(roleId, permissionId)
	if err != nil {
		return role, err
	}
	return role, nil
}
