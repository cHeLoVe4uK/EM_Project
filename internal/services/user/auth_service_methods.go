package user

import "github.com/cHeLoVe4uK/EM_Project/internal/models"

func (us *UserService) Login(u *models.User) (string, string, error) {
	// Для начала нужно проверить есть ли пользователь в БД
	ok, err := us.userRepo.CheckUserByID(u.ID)
	if err != nil {
		return "", "", err
	}
	if !ok {
		return "", "", UserNotFound
	}

	// Если есть нужно выбить для пользователя токены
	accessToken, refreshToken, err := us.authService.GetTokens(u)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (us *UserService) Logout(u *models.User) error {
	return nil
}
