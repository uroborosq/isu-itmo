package adapters

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role int

const (
	User Role = iota + 1
	Admin
)

type UserFullData struct {
	ID          uuid.UUID
	Isu         int
	Email       string
	PhoneNumber string
	Role        Role
	FullName    string
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) (UserRepo, error) {
	if !db.Migrator().HasTable(UserFullData{}) {
		err := db.Migrator().CreateTable(UserFullData{})
		if err != nil {
			return UserRepo{}, err
		}
	}
	return UserRepo{db: db}, nil
}

func (r *UserRepo) AddUser(user UserFullData) error {
	result := r.db.Model(UserFullData{}).Create(&user)
	return result.Error
}

func (r *UserRepo) UpdateUser(user UserFullData) error {
	return r.db.Model(UserFullData{}).Where("ID = ?", user.ID).Save(&user).Error
}

func (r *UserRepo) GetAllUsers() ([]UserFullData, error) {
	var users []UserFullData
	result := r.db.Find(&users)
	return users, result.Error
}

func (r *UserRepo) FindByPhoneNumber(phoneNumber string) (UserFullData, error) {
	var user UserFullData
	result := r.db.First(&user, "phone_number = ?", phoneNumber)
	return user, result.Error
}

func (r *UserRepo) GetRole(id uuid.UUID) (Role, error) {
	var user UserFullData
	result := r.db.First(&user, "ID = ?", id)
	return user.Role, result.Error
}

func (r *UserRepo) UpdatePublicInfo(user UserFullData) error {
	return r.db.Model(UserFullData{}).
		Where("ID = ?", user.ID).
		Select("isu", "email", "phone_number", "full_name").
		Updates(user).Error
}
