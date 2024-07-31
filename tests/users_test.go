package all_test

import (
	"context"
	"log"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/nrednav/cuid2"
	db "github.com/GoogleCloudPlatform/golang-samples/run/helloworld/db/sqlc"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/utils"
)

func createRandomUser(t *testing.T) db.User {
	hashedPassword, err := utils.GenerateHashPassword("testing")
	if err != nil {
		log.Fatal("Failed to generate hash password", err)
	}

	generate, err := cuid2.Init(
        cuid2.WithLength(32),
    )
    if err != nil {
		log.Fatal("Failed to generate CUID", err.Error())
    }

    // This function generates an id with a length of 32
    id := generate()
	

	arg := db.CreateUserParams{
		ID: id,
		Lastname:       utils.RandomName(),
		Firstname:      utils.RandomName(),
		Email:          utils.RandomEmail(),
		Phone:          utils.RandomPhone(),
		Address:        "shdhd,hdhd, jdjdjd",
		HashedPassword: hashedPassword,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, user)
	assert.Equal(t, user.Email, arg.Email)
	assert.Equal(t, user.Firstname, arg.Firstname)
	assert.Equal(t, user.Lastname, arg.Lastname)
	assert.Equal(t, user.ID, arg.ID)
	assert.Len(t, user.ID, 32)
	assert.NotZero(t, user.ID)
	assert.WithinDuration(t, user.CreatedAt, time.Now(), 2*time.Second)
	assert.WithinDuration(t, user.UpdatedAt, time.Now(), 2*time.Second)
	assert.Equal(t, user.HashedPassword, arg.HashedPassword)

	return user
}
func TestCreateUser(t *testing.T) {
	userPlate := createRandomUser(t)

	user, err := testQueries.CreateUser(context.Background(), db.CreateUserParams{
		Lastname:       userPlate.Lastname,
		Firstname:      userPlate.Firstname,
		Email:          userPlate.Email,
		Phone:          userPlate.Phone,
		Address:        "shdhd,hdhd, jdjdjd",
		HashedPassword: userPlate.HashedPassword,
	})
	assert.Error(t, err)
	assert.Empty(t, user)

}

func TestUpdateUserPassword(t *testing.T) {

	user := createRandomUser(t)

	newHashedPassword, err := utils.GenerateHashPassword("new Password")
	if err != nil {
		log.Fatal("Failed to generate hash password", err)
	}

	arg := db.UpdateUserPasswordParams{
		ID:             user.ID,
		HashedPassword: newHashedPassword,
		UpdatedAt:      time.Now(),
	}

	UpdatedUser, err := testQueries.UpdateUserPassword(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, UpdatedUser)
	assert.Equal(t, UpdatedUser.HashedPassword, arg.HashedPassword)
	assert.WithinDuration(t, UpdatedUser.UpdatedAt, time.Now(), 2*time.Second)

}

func TestUpdateUser(t *testing.T) {

	user := createRandomUser(t)

	arg := db.UpdateUserParams{
		ID:        user.ID,
		Email:     utils.RandomEmail(),
		Phone:     utils.RandomPhone(),
		Address:   "74 Avenue Suite, idiroko, yanibo, ajah",
		UpdatedAt: time.Now(),
	}

	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, user)
	assert.Equal(t, updatedUser.Email, arg.Email)
	assert.Equal(t, updatedUser.Phone, arg.Phone)
	assert.Equal(t, updatedUser.Address, arg.Address)
	assert.WithinDuration(t, updatedUser.UpdatedAt, time.Now(), 2*time.Second)

}

func TestGetUserById(t *testing.T) {
	user := createRandomUser(t)

	getUser, err := testQueries.GetUserById(context.Background(), user.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, getUser)
	assert.Equal(t, getUser.ID, user.ID)
	assert.Equal(t, getUser.Firstname, user.Firstname)
}

func TestGetUserByEmail(t *testing.T) {
	user := createRandomUser(t)

	getUser, err := testQueries.GetUserByEmail(context.Background(), user.Email)
	assert.NoError(t, err)
	assert.NotEmpty(t, getUser)
	assert.Equal(t, getUser.Email, user.Email)
	assert.Equal(t, getUser.Firstname, user.Firstname)
}

func TestListAllUsers(t *testing.T) {

	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}
	arg := db.ListAllUsersParams{
		Limit:  10,
		Offset: 0,
	}

	allUsers, err := testQueries.ListAllUsers(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, allUsers)
	assert.Equal(t, int32(len(allUsers)), arg.Limit)
}

func TestDeleteUser(t *testing.T) {
	user := createRandomUser(t)

	err := testQueries.DeleteUser(context.Background(), user.ID)
	assert.NoError(t, err)

	getUser, err := testQueries.GetUserById(context.Background(), user.ID)
	assert.Error(t, err)
	assert.Empty(t, getUser)

}

func TestDeleteAllUser(t *testing.T) {
	user := createRandomUser(t)

	err := testQueries.DeleteAllUsers(context.Background())
	assert.NoError(t, err)

	getUser, err := testQueries.GetUserById(context.Background(), user.ID)
	assert.Error(t, err)
	assert.Empty(t, getUser)

}
