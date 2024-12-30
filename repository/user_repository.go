package repository

import (
	"event-management-system/models"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll() ([]models.User, error)
	FindAllUserEvent() ([]models.User, error)
	FindById(id int) (models.User, error)
	Save(user models.User) (models.User, error)
	Update(user models.User) (models.User, error)
	Delete(id int) (models.User, error)
	FindByUsername(username string) (models.User, error)
	FindByEmail(email string) (models.User, error)
	FindByArray(params string, value bool) ([]models.User, error)
	FindBySingle(column, value string) (models.User, error)
}

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepositoryImpl {
	return &userRepositoryImpl{db: db}
}

func (u *userRepositoryImpl) FindAll() ([]models.User, error) {
	var users []models.User

	res := u.db.Find(&users)
	if res.Error != nil {
		fmt.Println("err:", res.Error)
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		fmt.Println("no users found")
		return users, nil
	}

	return users, nil
}

func (u *userRepositoryImpl) FindById(id int) (models.User, error) {
	var user models.User

	res := u.db.Preload("Orders.OrderDetails.Ticket.Event").First(&user, id)
	if res.Error != nil {
		fmt.Println("err:", res.Error)
		return user, res.Error
	}

	if res.RowsAffected == 0 {
		fmt.Println("no users found")
		return user, nil
	}

	return user, nil
}

func (u *userRepositoryImpl) FindByArray(params string, value bool) ([]models.User, error) {
	var user []models.User

	res := u.db.Where(fmt.Sprintf("%s = ?", params), value).Find(&user)
	if res.Error != nil {
		return user, res.Error
	}

	if res.RowsAffected == 0 {
		return user, gorm.ErrRecordNotFound
	}

	return user, nil
}
func (u *userRepositoryImpl) FindBySingle(column, value string) (models.User, error) {
	var user models.User

	res := u.db.Where(fmt.Sprintf("%s = ?", column), value).First(&user)
	if res.Error != nil {
		return user, res.Error
	}

	if res.RowsAffected == 0 {
		return user, gorm.ErrRecordNotFound
	}

	return user, nil
}

func (u *userRepositoryImpl) FindByUsername(username string) (models.User, error) {
	var user models.User

	res := u.db.Where("username = ?", username).First(&user)
	if res.Error != nil {
		return user, res.Error
	}

	if res.RowsAffected == 0 {
		return user, gorm.ErrRecordNotFound
	}

	return user, nil
}

func (u *userRepositoryImpl) FindByEmail(email string) (models.User, error) {
	var user models.User

	res := u.db.Where("email = ?", email).First(&user)
	if res.Error != nil {
		return user, res.Error
	}

	if res.RowsAffected == 0 {
		return user, gorm.ErrRecordNotFound
	}

	return user, nil
}

func (u *userRepositoryImpl) Save(user models.User) (models.User, error) {
	res := u.db.Create(&user)

	if res.Error != nil {
		fmt.Println("err:", res.Error)
		return user, res.Error
	}

	return user, nil
}

func (u *userRepositoryImpl) Update(user models.User) (models.User, error) {
	res := u.db.Updates(&user)

	if res.Error != nil {
		fmt.Println("err:", res.Error)
		return user, res.Error
	}

	return user, nil
}

func (u *userRepositoryImpl) Delete(id int) (models.User, error) {
	checkId, err := u.FindById(id)

	if err != nil {
		return checkId, err
	}

	var user models.User

	res := u.db.Delete(&user, id)

	if res.Error != nil {
		fmt.Println("err:", res.Error)
		return user, res.Error
	}

	return checkId, nil
}

func (u *userRepositoryImpl) FindAllUserEvent() ([]models.User, error) {
	var users []models.User

	res := u.db.Preload("Events.Tickets").Find(&users)
	// res := u.db.Preload("Events", func(db *gorm.DB) *gorm.DB {
	// 	return db.Omit("User")
	// }).Find(&users)
	if res.Error != nil {
		fmt.Println("err:", res.Error)
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		fmt.Println("no event users found")
		return users, nil
	}

	return users, nil
}
