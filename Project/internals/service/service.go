package service

import (
	"project/helper"
	"project/internals/storage"

	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	storage storage.UserStorage
}

type UserService interface {
	CreateUser(firstName, lastName, phone string, car []storage.CarDTO) (storage.User, error)
	GetUserById(id string) (storage.User, error)
	GetUsers() ([]storage.User, error)
	UpdateUser(id, firstName, lastName, phone string, car []storage.CarDTO) (storage.User, error)
	DeleteUser(id string) error
}

func NewService(storage storage.UserStorage) UserService {
	return &Service{
		storage: storage,
	}
}

func (s *Service) CreateUser(
	firstName,
	lastName,
	phone string,
	car []storage.CarDTO) (storage.User, error) {
	trimF := strings.TrimSpace(firstName)

	if trimF == "" {
		return storage.User{}, helper.Invalid

	}
	trimL := strings.TrimSpace(lastName)

	if trimL == "" {
		return storage.User{}, helper.Invalid

	}
	trimP := strings.TrimSpace(phone)

	if trimP == "" {
		return storage.User{}, helper.Invalid

	}

	id_user := uuid.New().String()

	usersCar := make([]storage.Car, 0, len(car))

	for _, carDTO := range car {
		usersCar = append(usersCar, storage.Car{
			ID:      uuid.New().String(),
			Model:   carDTO.Model,
			CarType: carDTO.Car_type,
			UserId:  id_user,
		})
	}

	timeNow := time.Now()

	newUser := storage.User{

		ID:        id_user,
		FirstName: trimF,
		LastName:  trimL,
		Phone:     trimP,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		Car:       usersCar,
	}

	err := s.storage.CreateUser(newUser.ID, newUser.FirstName, newUser.LastName, newUser.Phone)

	if err != nil {
		return storage.User{}, err
	}

	for _, car := range usersCar {
		if err := s.storage.CreateCar(car.ID, car.Model, car.CarType, car.UserId); err != nil {
			return storage.User{}, err
		}
	}

	return newUser, nil

}
func (s *Service) GetUserById(id string) (storage.User, error) {

	uuid, err := uuid.Parse(id)
	if err != nil {
		return storage.User{}, helper.Invalid
	}

	id = uuid.String()

	user, err := s.storage.GetUser(id)

	if err != nil {
		return storage.User{}, helper.UserNotFound
	}

	return user, nil

}

func (s *Service) GetUsers() ([]storage.User, error) {
	return s.storage.GetUsers()
}

func (s *Service) UpdateUser(
	id,
	firstName,
	lastName,
	phone string,
	cars []storage.CarDTO) (storage.User, error) {

	trimF := strings.TrimSpace(firstName)

	if trimF == "" {
		return storage.User{}, helper.Invalid

	}
	trimL := strings.TrimSpace(lastName)

	if trimL == "" {
		return storage.User{}, helper.Invalid

	}
	trimP := strings.TrimSpace(phone)

	if trimP == "" {
		return storage.User{}, helper.Invalid

	}

	uuid, err := uuid.Parse(id)

	if err != nil {
		return storage.User{}, helper.Invalid

	}

	id = uuid.String()

	err = s.storage.UpdateUser(id, trimF, trimL, trimP)
	if err != nil {
		return storage.User{}, err
	}

	for _, car := range cars {

		err := s.storage.UpdateCar(car.CarDTO_ID, car.Model, car.Car_type, id)

		if err != nil {
			return storage.User{}, err
		}
	}

	user, err := s.storage.GetUser(id)
	if err != nil {
		return storage.User{}, helper.UserNotFound

	}
	return user, nil
}

func (s *Service) DeleteUser(id string) error {

	uuid, err := uuid.Parse(id)
	if err != nil {
		return helper.Invalid
	}

	id = uuid.String()

	err = s.storage.DeleteUser(id)
	if err != nil {
		return err
	}

	return nil

}
