package store_test

import (
	"github.com/stretchr/testify/assert"
	"server/internal/app/model"
	"server/internal/app/store"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	//инициализируем тестовую базу...
	s, teardown := store.TestStore(t, dataBaseURL)
	defer teardown("users")

	//ожидаем ошибку..
	_, err := s.User().Create(&model.User{
		Name: "testName",
		Age:  999,
	})
	assert.Error(t, err)

	//не ожидаем ошибку..
	u, err := s.User().Create(&model.User{
		Name:    "testName",
		Age:     999,
		Friends: []int64{},
	})
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_UpdateUser(t *testing.T) {
	//инициализируем тестовую базу...
	s, teardown := store.TestStore(t, dataBaseURL)
	defer teardown("users")

	//первый кейс, ожидаем ошибку...
	m := &model.User{
		Name:    "testName1",
		Age:     777,
		Friends: []int64{},
	}
	_, err := s.User().UpdateUser(m)
	assert.Error(t, err)

	//второй кейс, не ожидаем ошибку...
	mU := &model.User{
		Name:    "testName2",
		Age:     777,
		Friends: []int64{},
	}
	s.User().Create(mU)
	us, err := s.User().UpdateUser(mU)
	assert.NoError(t, err)
	assert.NotNil(t, us)

}

func TestUserRepository_SetFriends(t *testing.T) {
	//инициализируем тестовую базу...
	s, teardown := store.TestStore(t, dataBaseURL)
	defer teardown("users")

	//первый кейс, ожидаем ошибку...
	m := &model.User{
		Name:    "testName1",
		Age:     777,
		Friends: []int64{},
	}
	_, err := s.User().SetFriends(m)
	assert.Error(t, err)

	//второй кейс, не ожидаем ошибку...
	mU := &model.User{
		Name:    "testName2",
		Age:     777,
		Friends: []int64{},
	}
	s.User().Create(mU)
	us, err := s.User().SetFriends(mU)
	assert.NoError(t, err)
	assert.NotNil(t, us)

}

func TestUserRepository_FindById(t *testing.T) {
	//инициализируем тестовую базу...
	s, teardown := store.TestStore(t, dataBaseURL)
	defer teardown("users")

	//инициализируем экземпляр user...
	mU := &model.User{
		Name:    "testName2",
		Age:     777,
		Friends: []int64{},
	}
	_, err := s.User().Create(mU)

	us, err := s.User().FindById(mU.ID)
	//не ожидаем ошибки
	assert.NoError(t, err)

	//mU равен us
	assert.Equal(t, mU, us)
}

func TestUserRepository_FindByName(t *testing.T) {
	//инициализируем тестовую базу...
	s, teardown := store.TestStore(t, dataBaseURL)
	defer teardown("users")

	//инициализируем экземпляр user...
	mU := &model.User{
		Name:    "testName2",
		Age:     777,
		Friends: []int64{},
	}
	_, err := s.User().Create(mU)

	us, err := s.User().FindByName(mU.Name)
	//не ожидаем ошибки
	assert.NoError(t, err)

	//mU равен us
	assert.Equal(t, mU, us)
}

func TestUserRepository_GetAll(t *testing.T) {
	//инициализируем тестовую базу...
	s, teardown := store.TestStore(t, dataBaseURL)
	defer teardown("users")

	//инициализируем экземпляры user...
	testUsers := []model.User{
		{
			Name:    "testName1",
			Age:     777,
			Friends: []int64{},
		},
		{
			Name:    "testName2",
			Age:     888,
			Friends: []int64{},
		},
		{
			Name:    "testName3",
			Age:     999,
			Friends: []int64{},
		},
	}

	for _, us := range testUsers {
		_, err := s.User().Create(&us)

		//не ожидаем ошибки
		assert.NoError(t, err)
	}

	mass, err := s.User().GetAll()
	//не ожидаем ошибки
	assert.NoError(t, err)
	assert.NotNil(t, mass)

	//проверка по именам в исходном и полученном массивах users, не ожидаем ошибки
	for i, res := range *mass {
		assert.Equal(t, res.Name, testUsers[i].Name)
	}
	//длина исходном и полученном массивах users равна, не ожидаем ошибки
	assert.Equal(t, len(testUsers), len(*mass))
}

func TestUserRepository_DeleteByID(t *testing.T) {
	//инициализируем тестовую базу...
	s, teardown := store.TestStore(t, dataBaseURL)
	defer teardown("users")

	//инициализируем экземпляр user...
	m := &model.User{
		Name:    "testName1",
		Age:     777,
		Friends: []int64{},
	}
	us, err := s.User().Create(m)
	assert.NoError(t, err)

	//первый кейс, не ожидаем ошибку...
	err = s.User().DeleteByID(us.ID)
	assert.NoError(t, err)

	//второй кейс, ожидаем ошибку...
	err = s.User().DeleteByID(0)
	assert.Error(t, err)

}

func TestUserRepository_ClearDeleteUserFromFriends(t *testing.T) {
	//инициализируем тестовую базу...
	s, teardown := store.TestStore(t, dataBaseURL)
	defer teardown("users")

	//инициализируем экземпляры user...
	testUsers := []model.User{
		{
			Name:    "testName1",
			Age:     777,
			Friends: []int64{},
		},
		{
			Name:    "testName2",
			Age:     888,
			Friends: []int64{},
		},
		{
			Name:    "testName3",
			Age:     999,
			Friends: []int64{},
		},
	}

	//create id mass...
	idMass := make([]int64, 0)

	for _, us := range testUsers {
		v, err := s.User().Create(&us)

		//create id mass...
		idMass = append(idMass, v.ID)

		//не ожидаем ошибки
		assert.NoError(t, err)
	}

	//get ID user testName1
	usTest1, err := s.User().FindByName("testName1")
	//get ID user testName2
	usTest2, err := s.User().FindByName("testName2")
	//get ID user testName3
	usTest3, err := s.User().FindByName("testName3")

	//make friends
	usTest1.Friends = []int64{usTest2.ID}
	_, err = s.User().UpdateUser(usTest1)
	assert.NoError(t, err)

	usTest2.Friends = []int64{usTest1.ID, usTest3.ID}
	_, err = s.User().UpdateUser(usTest2)
	assert.NoError(t, err)

	usTest3.Friends = []int64{usTest2.ID}
	_, err = s.User().UpdateUser(usTest3)
	assert.NoError(t, err)

	err = s.User().ClearDeleteUserFromFriends(usTest2.ID, idMass)
	//не ожидаем ошибки
	assert.NoError(t, err)
}
