package repository_test

import (
	"testing"

	"github.com/Bek0sh/online-market-auth/internal/models"
	"github.com/Bek0sh/online-market-auth/pkg/utils"
	"github.com/stretchr/testify/require"
)

func createUser(t *testing.T) (int, *models.RegisterUser) {

	hashedPassword, err := utils.HashPassword(utils.RandomString(12))
	require.NoError(t, err)

	user := &models.RegisterUser{
		Username:    utils.RandomUsername(),
		Surname:     utils.RandomUsername(),
		PhoneNumber: utils.RandomPhoneNumber(),
		Password:    hashedPassword,
	}

	id, err := repo.CreateUser(
		user,
	)

	require.NoError(t, err)
	require.NotEmpty(t, id)

	return id, user
}

func TestCreateUser(t *testing.T) {
	createUser(t)
}

func TestGetUserById(t *testing.T) {
	id, user := createUser(t)

	userWithId, err := repo.GetUserById(id)
	require.NoError(t, err)
	require.NotEmpty(t, userWithId)

	require.Equal(t, user.Username, userWithId.Username)
	require.Equal(t, user.Surname, userWithId.Surname)
	require.Equal(t, user.PhoneNumber, userWithId.PhoneNumber)
}

func TestGetUserByPhoneNumber(t *testing.T) {
	_, user := createUser(t)

	userWithPhoneNumber, err := repo.GetUserByPhoneNumber(user.PhoneNumber)
	require.NoError(t, err)
	require.NotEmpty(t, userWithPhoneNumber)

	require.Equal(t, user.Username, userWithPhoneNumber.Username)
	require.Equal(t, user.Surname, userWithPhoneNumber.Surname)
	require.Equal(t, user.PhoneNumber, userWithPhoneNumber.PhoneNumber)
}

func TestUpdateUser(t *testing.T) {
	id, user := createUser(t)

	u1 := &models.User{
		Id:       id,
		Username: utils.RandomUsername(),
		Surname:  utils.RandomUsername(),
	}

	resp, err := repo.UpdateUser(u1)

	require.NoError(t, err)
	require.NotEmpty(t, resp)

	require.Equal(t, resp.Username, u1.Username)
	require.Equal(t, resp.Surname, u1.Surname)
	require.Equal(t, resp.PhoneNumber, user.PhoneNumber)
	require.Equal(t, resp.Id, id)

}

func TestDeleteUser(t *testing.T) {
	id, _ := createUser(t)

	err := repo.DeleteUser(id)
	require.NoError(t, err)

	u1, err := repo.GetUserById(id)
	require.Empty(t, u1)
	require.Error(t, err)
}
