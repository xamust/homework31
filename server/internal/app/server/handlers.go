package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"reflect"
	"server/internal/app/model"
	"server/internal/app/store"
	"strconv"
)

type Handlers struct {
	logger *logrus.Logger
	mux    *mux.Router
	userZ  store.UserRepo
}

// Create CreateUser
func (s *Handlers) Create(w http.ResponseWriter, r *http.Request) {

	var u *model.User

	if r.Method == http.MethodPost && r.Header.Get("Content-Type") == "application/json" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		defer r.Body.Close()

		if err := json.Unmarshal(content, &u); err != nil {
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		//весьма спорно... захардкодил
		if reflect.TypeOf(u.Friends).String() == "[]int64" && u.Friends == nil {
			u.Friends = make([]int64, 0)
		}

		if _, err = s.userZ.Create(u); err != nil {
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		//s.store[u.Name] = u //for map[] store
		//GlobalId += 1 //for map[] store
		//u.ID = GlobalId //костылёк с ID, for map[] store

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("id:%d", u.ID)))

		s.logger.Info(fmt.Sprintf("Добавлен пользователь: id=%d, name=%s, age=%d, friends=%v", u.ID, u.Name, u.Age, u.Friends))
		return
	}
	w.WriteHeader(http.StatusBadRequest)

}

// MakeFriends MakeFriends
func (s *Handlers) MakeFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && r.Header.Get("Content-Type") == "application/json" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		defer r.Body.Close()

		var fm *model.FriendsMaker

		if err := json.Unmarshal(content, &fm); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		//userFirst := s.GetUser(fm.SourceId) //for map[] store

		userFirst, err := s.userZ.FindById(fm.SourceId)
		if err != nil {
			s.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		//userSecond := s.GetUser(fm.TargetId)                          //for map[] store

		userSecond, err := s.userZ.FindById(fm.TargetId)
		if err != nil {
			s.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		//check friendship...
		for _, v := range userFirst.Friends {
			if v == userSecond.ID {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Внимание! %s и %s уже друзья!!!", userFirst.Name, userSecond.Name)))
				s.logger.Info(fmt.Sprintf("Внимание! %s и %s уже друзья!!!", userFirst.Name, userSecond.Name))
				return
			}
		}

		//success...
		userFirst.Friends = append(userFirst.Friends, userSecond.ID)
		if _, err := s.userZ.SetFriends(userFirst); err != err {
			s.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		userSecond.Friends = append(userSecond.Friends, userFirst.ID)
		if _, err := s.userZ.SetFriends(userSecond); err != err {
			s.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%s и %s теперь стали друзьями.", userFirst.Name, userSecond.Name)))
		s.logger.Info(fmt.Sprintf("%s и %s теперь стали друзьями.", userFirst.Name, userSecond.Name))

		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

//Delete DeleteUser
func (s *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete && r.Header.Get("Content-Type") == "application/json" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		defer r.Body.Close()

		var ud *model.FriendsMaker

		if err := json.Unmarshal(content, &ud); err != nil {
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		//проверка наличия имени
		user, err := s.userZ.FindById(ud.TargetId)
		if err != nil {
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		//user := s.GetUser(ud.TargetId)
		//_, ok := s.store[user.Name] //проверка наличия имени
		//if ok {
		//	delete(s.store, user.Name) //delete for map[]
		//	s.DeleteByID(user.ID)
		//}

		if err := s.userZ.DeleteByID(ud.TargetId); err != nil {
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%s", user.Name)))
		s.logger.Info(fmt.Sprintf("Пользователь %s удален.", user.Name))
		return
	}
}

// GetFriends GetUserFriends
func (s *Handlers) GetFriends(w http.ResponseWriter, r *http.Request) {

	us := &model.User{}

	//конструкция получения выполнения регулярки
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		s.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if r.Method == http.MethodGet {
		response := ""

		//user := s.GetUser(int64(id)) //for map[] store...

		us, err = s.userZ.FindById(int64(id))
		if err != nil {
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		//проверка наличия имени, for map[] store
		//if user == nil {
		//	w.WriteHeader(http.StatusInternalServerError)
		//	w.Write([]byte("Нет пользователя с таким id"))
		//	return
		//}

		//	for _, userFriend := range us.Friends {
		//		response += fmt.Sprintf("%s\n", s.GetUser(userFriend).Name) //get username by id
		//		}
		//if response == "" {
		//	response = fmt.Sprintf("У %s нет друзей :(\n", us.Name)
		//} else {
		//	response = fmt.Sprintf("Друзья %s:\n%s", us.Name, response)
		//}

		response = fmt.Sprintf("Друзья %s:\n", us.Name)
		if len(us.Friends) > 0 {
			for _, v := range us.Friends {
				usFr, _ := s.userZ.FindById(v)
				response += fmt.Sprintf("%s\n", usFr.Name)
			}
		} else {
			response = fmt.Sprintf("У %s нет друзей :(\n", us.Name)
		}

		s.logger.Info(response)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

// Put UpdateUserAge
func (s *Handlers) Put(w http.ResponseWriter, r *http.Request) {

	us := &model.User{}

	//конструкция получения выполнения регулярки
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		s.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	//check Content-Type
	if r.Method == http.MethodPut && r.Header.Get("Content-Type") == "application/json" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		defer r.Body.Close()

		var na *model.NewAge

		if err = json.Unmarshal(content, &na); err != nil {
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		//находим юзера
		us, err = s.userZ.FindById(int64(userId))
		if err != nil {
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		s.logger.Info(us, us.Age)
		//for map[] storage
		//user := s.GetUser(int64(userId))
		//проверка наличия имени
		//if user == nil {
		//	w.WriteHeader(http.StatusInternalServerError)
		//	w.Write([]byte("Нет пользователя с таким id"))
		//	return
		//}

		//check age...
		if na.NewAge < 0 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Возраст не может быть меньше 0!!!"))
			return
		}

		us.Age = na.NewAge
		s.userZ.UpdateUser(us)
		s.logger.Info(fmt.Sprintf("Возраст %s успешно обновлён до %d", us.Name, us.Age))
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprint("Возраст успешно обновлён")))

		return
	}
}

//
//my handlers by debug (without test)
//

// GetAll GetAllUsersInformation
func (s *Handlers) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		response := ""
		store, err := s.userZ.GetAll()
		if err != nil {
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		for _, user := range *store {
			response += fmt.Sprintf("%d Пользователь %s, возраст %d, друзья %v \n", user.ID, user.Name, user.Age, user.Friends)
		}
		s.logger.Info(response)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

// GetUserInfo GetUserInformation
func (s *Handlers) GetUserInfo(w http.ResponseWriter, r *http.Request) {

	us := &model.User{}

	//конструкция получения выполнения регулярки
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["user_id"])

	if err != nil {
		s.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if r.Method == http.MethodGet {
		us, err = s.userZ.FindById(int64(id))
		if err != nil {
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Нет пользователя с таким id"))
			return
		}

		s.logger.Info(fmt.Sprintf("id %d, name %s, age %d, friends %v\n", us.ID, us.Name, us.Age, us.Friends))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("id %d, name %s, age %d, friends %v\n", us.ID, us.Name, us.Age, us.Friends)))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}
