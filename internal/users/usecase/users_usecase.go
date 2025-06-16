package usecase

import (
	"errors"

	// "github.com/nightnice1st/testGridWhiz/internal/pkg/validator"
	"github.com/nightnice1st/testGridWhiz/internal/users/domain"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u *userUsecase) GetProfile(userID string) (*domain.User, error) {
	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Clear password before returning
	user.Password = ""
	return user, nil
}

func (u *userUsecase) UpdateProfile(userID, name string) (*domain.User, error) {
	// Get existing user
	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update fields if provided
	if name != "" {
		user.Name = name
	}

	// if email != "" && email != user.Email {
	// 	// Validate new email
	// 	if err := validator.ValidateEmail(email); err != nil {
	// 		return nil, err
	// 	}

	// 	// Check if email already exists
	// 	existingUser, _ := u.userRepo.FindByEmail(email)
	// 	if existingUser != nil && existingUser.ID != userID {
	// 		return nil, errors.New("email already in use")
	// 	}

	// 	user.Email = email
	// }

	// Update in database
	if err := u.userRepo.Update(user); err != nil {
		return nil, err
	}

	// Clear password before returning
	user.Password = ""
	return user, nil
}

func (u *userUsecase) DeleteProfile(userID string) error {
	// Check if user exists
	_, err := u.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	return u.userRepo.SoftDelete(userID)
}

func (u *userUsecase) ListUsers(page, limit int, nameFilter, emailFilter string) ([]*domain.User, int, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := u.userRepo.List(page, limit, nameFilter, emailFilter)
	if err != nil {
		return nil, 0, err
	}

	// Clear passwords before returning
	for _, user := range users {
		user.Password = ""
	}

	return users, total, nil
}
