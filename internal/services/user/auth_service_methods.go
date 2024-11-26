package user

import "github.com/cHeLoVe4uK/EM_Project/internal/models"

// Вход пользователя
func (us *UserService) Login(u *models.User) (string, string, error) {
	// Проверка наличия пользователя в БД
	ok, err := us.userRepo.CheckUserByID(u.ID)
	if err != nil {
		return "", "", err
	}
	if !ok {
		return "", "", ErrUserNotFound
	}

	// Если есть выбиваем токены
	accessToken, refreshToken, err := us.authService.GetTokens(u)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (us *UserService) Logout(u *models.User) error {
	return nil
}
