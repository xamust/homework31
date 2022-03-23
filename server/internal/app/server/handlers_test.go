package server

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"server/internal/app/model"
	"server/internal/app/store"
	"testing"
)

func TestAppServer_Create(t *testing.T) {

	s, serverTest, teardown := initTestDB(t)
	defer teardown("users")

	_, _, usTest3 := initData(t, s)

	//можно без массива структур, для примера
	testCases := []struct {
		name   string
		method string
		body   string
		want   []byte
	}{
		{
			name:   "Testing create user handler...",
			method: http.MethodPost,
			//
			body: fmt.Sprintf(`{"name" : "%s", "age" : %d}`, "testName4", 333),
			want: []byte(fmt.Sprintf("id:%d", usTest3.ID+1)),
		},
	}

	handlerTest := http.HandlerFunc(serverTest.Create)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest(tc.method, "/create", bytes.NewBuffer([]byte(tc.body)))
			req.Header.Add("Content-Type", "application/json")
			handlerTest.ServeHTTP(rec, req)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, rec.Body.Bytes())
		})
	}
}

func TestAppServer_MakeFriends(t *testing.T) {

	s, serverTest, teardown := initTestDB(t)
	defer teardown("users")
	usTest1, usTest2, usTest3 := initData(t, s)

	testCases := []struct {
		name   string
		method string
		body   string
		want   []byte
	}{
		{
			name:   fmt.Sprintf("Testing make friends %s and %s, handler...", usTest1.Name, usTest2.Name),
			method: http.MethodPost,
			body:   fmt.Sprintf(`{"source_id" : %d, "target_id" : %d}`, usTest1.ID, usTest2.ID),
			want:   []byte("теперь стали друзьями."),
		},
		{
			name:   fmt.Sprintf("Testing make friends %s and %s, handler...", usTest2.Name, usTest3.Name),
			method: http.MethodPost,
			body:   fmt.Sprintf(`{"source_id" : %d, "target_id" : %d}`, usTest2.ID, usTest3.ID),
			want:   []byte("теперь стали друзьями."),
		},
	}
	handlerTest := http.HandlerFunc(serverTest.MakeFriends)
	//не ожидаем ошибки (новые друзья)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest(tc.method, "/make_friends", bytes.NewBuffer([]byte(tc.body)))
			req.Header.Add("Content-Type", "application/json")
			handlerTest.ServeHTTP(rec, req)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, rec.Body.Bytes()[len(rec.Body.Bytes())-len(tc.want):])
		})
	}
	//не ожидаем ошибку (уже друзья)
	testCase2 := struct {
		name   string
		method string
		body   string
		want   []byte
	}{
		name:   fmt.Sprintf("Testing make friends %s and %s, handler...", usTest1.Name, usTest2.Name),
		method: http.MethodPost,
		body:   fmt.Sprintf(`{"source_id" : %d, "target_id" : %d}`, usTest1.ID, usTest2.ID),
		want:   []byte("уже друзья!!!"),
	}

	t.Run(testCase2.name, func(t *testing.T) {
		rec := httptest.NewRecorder()
		req, err := http.NewRequest(testCase2.method, "/make_friends", bytes.NewBuffer([]byte(testCase2.body)))
		req.Header.Add("Content-Type", "application/json")
		handlerTest.ServeHTTP(rec, req)
		assert.NotNil(t, rec.Body.Bytes())
		assert.NoError(t, err)
		assert.Equal(t, testCase2.want, rec.Body.Bytes()[len(rec.Body.Bytes())-len(testCase2.want):])
	})
}

func TestAppServer_Delete(t *testing.T) {
	s, serverTest, teardown := initTestDB(t)
	defer teardown("users")
	usTest1, _, _ := initData(t, s)

	handlerTest := http.HandlerFunc(serverTest.Delete)

	//не ожидаем ошибку
	testCase := struct {
		name   string
		method string
		body   string
		want   []byte
	}{
		name:   fmt.Sprintf("Testing delete %s, handler...", usTest1.Name),
		method: http.MethodDelete,
		body:   fmt.Sprintf(`{"target_id" : %d}`, usTest1.ID),
		want:   []byte(fmt.Sprintf("%s", usTest1.Name)),
	}

	t.Run(testCase.name, func(t *testing.T) {
		rec := httptest.NewRecorder()
		req, err := http.NewRequest(testCase.method, "/user", bytes.NewBuffer([]byte(testCase.body)))
		req.Header.Add("Content-Type", "application/json")
		handlerTest.ServeHTTP(rec, req)
		assert.NoError(t, err)
		assert.Equal(t, testCase.want, rec.Body.Bytes())
	})

	//не ожидаем ошибку (ошибка обработана в store)
	testCase = struct {
		name   string
		method string
		body   string
		want   []byte
	}{
		name:   fmt.Sprintf("Testing delete %s, handler...", ""),
		method: http.MethodDelete,
		body:   fmt.Sprintf(`{"target_id" : %d}`, 0),
		want:   []byte(fmt.Sprintf("нет в базе данных")),
	}

	t.Run(testCase.name, func(t *testing.T) {
		rec := httptest.NewRecorder()
		req, err := http.NewRequest(testCase.method, "/user", bytes.NewBuffer([]byte(testCase.body)))
		req.Header.Add("Content-Type", "application/json")
		handlerTest.ServeHTTP(rec, req)
		//TODO Изучить другие asserts!!!
		//assert.HTTPError(t, proxy.Delete, http.MethodDelete, "/user", url.Values{"target_id": []string{"100"}}, "")
		assert.NoError(t, err)
		assert.Contains(t, string(rec.Body.Bytes()), string(testCase.want))
		//можно и assert, но нужно выделять конкретные байты, а это костыль
		//assert.Equal(t, testCase.want, rec.Body.Bytes())
	})

}

func TestAppServer_GetFriends(t *testing.T) {
	s, serverTest, teardown := initTestDB(t)
	defer teardown("users")
	usTest1, usTest2, usTest3 := initData(t, s)

	//make usTest1 and usTest2 friends, usTest3 no friends
	usTest1.Friends = append(usTest1.Friends, usTest2.ID)
	_, err := s.User().UpdateUser(usTest1)
	assert.NoError(t, err)
	usTest2.Friends = append(usTest2.Friends, usTest1.ID)
	_, err = s.User().UpdateUser(usTest2)
	assert.NoError(t, err)

	testCases := []struct {
		name   string
		id     int64
		method string
		want   []byte
	}{
		{
			name:   fmt.Sprintf("Testing get friends %s handler...", usTest1.Name),
			id:     usTest1.ID,
			method: http.MethodGet,
			want:   []byte("Друзья"),
		},
		{
			name:   fmt.Sprintf("Testing get friends %s handler...", usTest2.Name),
			id:     usTest2.ID,
			method: http.MethodGet,
			want:   []byte("Друзья"),
		},
		{
			name:   fmt.Sprintf("Testing get friends %s handler...", usTest3.Name),
			id:     usTest3.ID,
			method: http.MethodGet,
			want:   []byte("нет друзей"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest(tc.method, fmt.Sprintf("/friends/%d", tc.id), nil)

			//another way... with gorilla mux
			serverTest.mux.HandleFunc("/friends/{id:[0-9]+}", serverTest.GetFriends)
			serverTest.mux.ServeHTTP(rec, req)

			assert.NoError(t, err)
			assert.Contains(t, string(rec.Body.Bytes()), string(tc.want))
		})
	}
}

func TestAppServer_Put(t *testing.T) {
	s, serverTest, teardown := initTestDB(t)
	defer teardown("users")
	usTest1, usTest2, usTest3 := initData(t, s)

	testCases := []struct {
		name   string
		id     int64
		method string
		want   []byte
		body   string
	}{
		{
			name:   fmt.Sprintf("Testing put %s handler...", usTest1.Name),
			id:     usTest1.ID,
			method: http.MethodPut,
			want:   []byte("Возраст успешно обновлён"),
			body:   fmt.Sprintf(`{"new_age" : %d}`, 100),
		},
		{
			name:   fmt.Sprintf("Testing put %s handler...", usTest2.Name),
			id:     usTest2.ID,
			method: http.MethodPut,
			want:   []byte("Возраст успешно обновлён"),
			body:   fmt.Sprintf(`{"new_age" : %d}`, 102),
		},
		{
			name:   fmt.Sprintf("Testing put %s handler...", usTest3.Name),
			id:     usTest3.ID,
			method: http.MethodPut,
			want:   []byte("Возраст не может быть меньше 0!!!"),
			body:   fmt.Sprintf(`{"new_age" : %d}`, -88),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest(tc.method, fmt.Sprintf("/%d", tc.id), bytes.NewBuffer([]byte(tc.body)))
			req.Header.Add("Content-Type", "application/json")

			//another way... with gorilla mux
			serverTest.mux.HandleFunc("/{user_id:[0-9]+}", serverTest.Put)
			serverTest.mux.ServeHTTP(rec, req)

			assert.NoError(t, err)
			assert.Equal(t, string(tc.want), string(rec.Body.Bytes()))
		})
	}

}

func initTestDB(t *testing.T) (*store.AppStore, *AppServer, func(...string)) {

	dataBaseURL := os.Getenv("DATABASE_URL")
	if dataBaseURL == "" {
		dataBaseURL = "dbname=testDB user=postgres password=example host=localhost port=5433 sslmode=disable"
	}
	//инициализируем тестовую базу...
	s, teardown := store.TestStore(t, dataBaseURL)

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

	//инициализируем экзепляр server, с требуемыми параметрами
	serverTest := &AppServer{
		logger:  logrus.New(),
		mux:     mux.NewRouter(),
		storeBD: s,
	}
	//level panic for hide log output
	serverTest.logger.SetLevel(logrus.PanicLevel)
	return s, serverTest, teardown
}

func initData(t *testing.T, s *store.AppStore) (*model.User, *model.User, *model.User) {
	//get ID user testName1
	usTest1, err := s.User().FindByName("testName1")
	assert.NoError(t, err)
	//get ID user testName2
	usTest2, err := s.User().FindByName("testName2")
	assert.NoError(t, err)
	//get ID user testName3
	usTest3, err := s.User().FindByName("testName3")
	assert.NoError(t, err)
	return usTest1, usTest2, usTest3
}
