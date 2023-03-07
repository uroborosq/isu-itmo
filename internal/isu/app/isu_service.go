package app

import (
	"github.com/google/uuid"
	"github.com/uroborosq/isu/internal/isu/adapters"
)

type IsuService struct {
	userRepo adapters.UserRepo
}

func NewIsuService(userRepo adapters.UserRepo) IsuService {
	return IsuService{
		userRepo: userRepo,
	}
}

func (s *IsuService) AddUser(user adapters.UserFullData) error {
	err := ValidateFullInfo(user)
	if err != nil {
		return err
	}
	return s.userRepo.AddUser(user)
}

func (s *IsuService) GetPublicInfo(phoneNumber string) (adapters.UserFullData, error) {
	return s.userRepo.FindByPhoneNumber(phoneNumber)
}

func (s *IsuService) UpdatePublicInfo(user adapters.UserFullData) error {
	err := ValidateFullInfo(user)
	if err != nil {
		return err
	}
	return s.userRepo.UpdatePublicInfo(user)
}

func (s *IsuService) GetAllUsers() ([]adapters.UserFullData, error) {
	return s.userRepo.GetAllUsers()
}

func (s *IsuService) GetRole(id uuid.UUID) (adapters.Role, error) {
	return s.userRepo.GetRole(id)
}

func (s *IsuService) UpdateFullInfo(user adapters.UserFullData) error {
	err := ValidateFullInfo(user)
	if err != nil {
		return err
	}
	return s.userRepo.UpdateUser(user)
}
