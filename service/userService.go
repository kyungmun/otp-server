package service

// OTP 비즈니스 로직만 처리하는 서비스

import (
	"github.com/kyungmun/otp-server/models"
	"github.com/kyungmun/otp-server/repository"
)

func NewUserServices(r *repository.UserRepository) (*UserServices, error) {
	return &UserServices{repo: r}, nil
}

type UserServices struct {
	repo *repository.UserRepository
}

func (s *UserServices) GetAll(page, pageSize int) (*[]models.User, error) {
	users, err := s.repo.GetAll(page, pageSize)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserServices) GetByID(username string) (*models.User, error) {

	user, err := s.repo.GetByID(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserServices) GetByEmail(email string) (*models.User, error) {

	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserServices) UpdateRecord(user *models.User) (*models.User, error) {

	userNew, err := s.repo.UpdateRecord(user)

	if err != nil {
		return nil, err
	}

	return userNew, nil
}

func (s *UserServices) PatchRecord(username string, jsonData map[string]interface{}) (*models.User, error) {

	user, err := s.repo.PatchRecord(username, jsonData)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserServices) DeleteRecord(username string) error {

	err := s.repo.DeleteByID(username)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserServices) CreateRecord(user *models.User) (*models.User, error) {

	userNew, err := s.repo.Create(user)

	if err != nil {
		return nil, err
	}

	return userNew, nil
}
