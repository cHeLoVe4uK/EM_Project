package user_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/cHeLoVe4uK/EM_Project/internal/repository/user_repository"
	"github.com/cHeLoVe4uK/EM_Project/internal/services/auth"
	"github.com/cHeLoVe4uK/EM_Project/internal/services/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	// Подготовка к тесту
	os.Setenv("TOKEN_SALT", "test_secret_key")
	os.Setenv("ACCESS_TOKEN_EXP", "24")
	defer os.Unsetenv("TOKEN_SALT")
	defer os.Unsetenv("ACCESS_TOKEN_EXP")
	ttl, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXP"))
	require.NoError(t, err, "no error should occuer while converting ACCESS_TOKEN_EXP to int")

	// Тест имитирующий проблемы соединения с БД (независящие от нашего приложения)
	usRepo := &Mock{}
	userService := user.NewUserService(usRepo, auth.NewService(os.Getenv("TOKEN_SALT"), ttl))
	_, err = userService.Register(context.Background(), models.User{})
	require.EqualError(t, err, ErrWithDB.Error(), fmt.Sprintf("err: %s, expected err: %s - must be equal", err, ErrWithDB))

	// Тесты с доступом к БД
	userService = user.NewUserService(New(), auth.NewService(os.Getenv("TOKEN_SALT"), ttl))
	// Верные тест-кейсы
	rightCases := []struct {
		Name string
		models.User
	}{
		{
			Name: "first_right_test",
			User: models.User{
				Email:    "example@mail.ru",
				Username: "Tetete",
				Password: "tut_parol_1",
			},
		},
		{
			Name: "second_right_test",
			User: models.User{
				Email:    "example2@mail.ru",
				Username: "Tututu",
				Password: "tut_parol_2",
			},
		},
		{
			Name: "third_right_test",
			User: models.User{
				Email:    "example3@mail.ru",
				Username: "Tatata",
				Password: "tut_parol_3",
			},
		},
	}

	for _, v := range rightCases {
		ctx := context.Background()
		t.Run(v.Name, func(t *testing.T) {
			id := uuid.NewString()
			v.ID = id
			res, err := userService.Register(ctx, v.User)
			require.Equal(t, res, v.ID, fmt.Sprintf("id: %s, v.user.ID: %s - mast be equal", res, v.ID))
			require.NoError(t, err, "no error should occuer while register user")
		})
	}

	// Тест-кейсы на ошибки
	badCases := []struct {
		Name string
		models.User
		Err error
	}{
		{
			// Попробуем добавить пользователя который уже существует
			Name: "first_bad_test",
			User: models.User{
				Email:    "example@mail.ru",
				Username: "Tetete",
				Password: "tut_parol_1",
			},
			Err: user.ErrUserExists,
		},
		{
			// Попробуем добавить пользователя с недопустимой длиной пароля
			Name: "second_bad_test",
			User: models.User{
				Email:    "example4@mail.ru",
				Username: "Nevazno",
				Password: "Специально тут надо сделать пароль более семидесяти двух байт, чтобы пароль не смог сгенерироваться",
			},
			Err: user.ErrHashPassword,
		},
	}

	for _, v := range badCases {
		ctx := context.Background()
		t.Run(v.Name, func(t *testing.T) {
			id := uuid.NewString()
			v.ID = id
			_, err := userService.Register(ctx, v.User)
			require.EqualError(t, err, v.Err.Error(), fmt.Sprintf("err: %s, expexted err: %s - must be equal", err, v.Err))
		})
	}
}

func TestUpdateUser(t *testing.T) {
	// Подготовка к тесту
	os.Setenv("TOKEN_SALT", "test_secret_key")
	os.Setenv("ACCESS_TOKEN_EXP", "24")
	defer os.Unsetenv("TOKEN_SALT")
	defer os.Unsetenv("ACCESS_TOKEN_EXP")
	ttl, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXP"))
	require.NoError(t, err, "no error should occuer while converting ACCESS_TOKEN_EXP to int")

	userService := user.NewUserService(New(), auth.NewService(os.Getenv("TOKEN_SALT"), ttl))

	// Добавим сначала пользователя перед тем как его изменять
	user := models.User{
		Email:    "example@mail.ru",
		Username: "Tetete",
		Password: "tut_parol_1",
	}
	id := uuid.NewString()
	user.ID = id
	userService.Register(context.Background(), user)

	// Тест с изменением существующего пользователя
	user = models.User{
		ID:       id,
		Email:    "example@mail.ru",
		Username: "Ratata",
		Password: "teper_drugoi_parol",
	}

	err = userService.UpdateUser(context.Background(), user)
	require.NoError(t, err, "no error should occuer while update existing user")

	// Тест с изменением несуществующего пользователя (bad case)
	user2 := models.User{
		Email:    "example34@mail.ru",
		Username: "Bad",
		Password: "teper_drugoi_parol_voob4e",
	}

	id2 := uuid.NewString()
	user2.ID = id2

	err = userService.UpdateUser(context.Background(), user2)
	require.EqualError(t, err, user_repository.ErrUserNotFound.Error(), "err: %s, expected err: %s - must be equal", err, user_repository.ErrUserNotFound)
}

func TestDeleteUser(t *testing.T) {
	// Подготовка к тесту
	os.Setenv("TOKEN_SALT", "test_secret_key")
	os.Setenv("ACCESS_TOKEN_EXP", "24")
	defer os.Unsetenv("TOKEN_SALT")
	defer os.Unsetenv("ACCESS_TOKEN_EXP")
	ttl, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXP"))
	require.NoError(t, err, "no error should occuer while converting ACCESS_TOKEN_EXP to int")

	userService := user.NewUserService(New(), auth.NewService(os.Getenv("TOKEN_SALT"), ttl))

	// Добавим сначала пользователя перед тем как его удалять
	user := models.User{
		Email:    "example@mail.ru",
		Username: "Tetete",
		Password: "tut_parol_1",
	}
	id := uuid.NewString()
	user.ID = id
	userService.Register(context.Background(), user)

	// Тест с удалением существующего пользователя
	userService.DeleteUser(context.Background(), user)
	require.NoError(t, err, "No error should occuer while delete existing user")

	// Тест с удалением несуществующего пользователя (bad case)
	user2 := models.User{
		Email:    "example34@mail.ru",
		Username: "Bad",
		Password: "teper_drugoi_parol_voob4e",
	}

	id2 := uuid.NewString()
	user2.ID = id2

	err = userService.DeleteUser(context.Background(), user2)
	require.EqualError(t, err, user_repository.ErrUserNotFound.Error(), "err: %s, expected err: %s - must be equal", err, user_repository.ErrUserNotFound)
}

func TestLogin(t *testing.T) {
	// Подготовка к тесту
	os.Setenv("TOKEN_SALT", "test_secret_key")
	os.Setenv("ACCESS_TOKEN_EXP", "24")
	defer os.Unsetenv("TOKEN_SALT")
	defer os.Unsetenv("ACCESS_TOKEN_EXP")
	ttl, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXP"))
	require.NoError(t, err, "no error should occuer while converting ACCESS_TOKEN_EXP to int")

	userService := user.NewUserService(New(), auth.NewService(os.Getenv("TOKEN_SALT"), ttl))

	// Добавим сначала пользователя перед тем как позволить ему логиниться
	u := models.User{
		Email:    "example@mail.ru",
		Username: "Tetete",
		Password: "tut_parol_1",
	}
	id := uuid.NewString()
	u.ID = id
	userService.Register(context.Background(), u)

	// Теперь сам тест
	tokens, err := userService.Login(context.Background(), u)
	require.NotEmpty(t, tokens, "tokens can't be empty")
	require.NoError(t, err, "no error should occuer while login")

	// Тест с неправильным паролем
	u.Password = "Izmenili_parol"

	_, err = userService.Login(context.Background(), u)
	require.EqualError(t, err, user.ErrInvalidPassword.Error(), "err: %s, expected err: %s - must be equal", err, user.ErrInvalidPassword)

	// Тест с несуществующим пользователем
	user2 := models.User{
		Email:    "exampleeeeee@mail.ru",
		Username: "Xvatit",
		Password: "da_kakaz_raznisa_4nj_tut",
	}
	id2 := uuid.NewString()
	user2.ID = id2
	_, err = userService.Login(context.Background(), user2)
	require.EqualError(t, err, user_repository.ErrUserNotFound.Error(), "err: %s, expected err: %s - must be equal", err, user_repository.ErrUserNotFound)
}
