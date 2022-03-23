package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"server/internal/app/store"
)

type AppServer struct {
	config  *Config
	mux     *mux.Router
	logger  *logrus.Logger
	store   map[string]*User //map ссылочный тип
	storeBD *store.AppStore
}

type User struct {
	ID      int64   //Users ID's
	Name    string  `json:"name"`
	Age     int64   `json:"age"`
	Friends []int64 `json:"friends"`
}

type FriendsMaker struct {
	SourceId int64 `json:"source_id"`
	TargetId int64 `json:"target_id"`
}

type NewAge struct {
	NewAge int64 `json:"new_age"`
}

//init new server
func New(config *Config) *AppServer {
	return &AppServer{
		config: config,
		mux:    mux.NewRouter(), //gorilla/mux
		logger: logrus.New(),    //sirupsen/logrus

	}
}

//configure logrus...
func (s *AppServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *AppServer) Start() error {

	if err := s.configureLogger(); err != nil {
		return err //if logrus configure result err
	}

	s.configureRouter() //configure router

	if err := s.configureStore(); err != nil {
		return err
	}

	s.logger.Info(fmt.Sprintf("Starting server (bind on %v)...", s.config.BindAddr)) // set message Info level about succesfull starting server...
	return http.ListenAndServe(s.config.BindAddr, s.mux)                             //bind addr from config and new gorilla mux
}

//config route...
func (s *AppServer) configureRouter() {

	s.mux.HandleFunc("/create", s.Create)
	s.mux.HandleFunc("/make_friends", s.MakeFriends)
	s.mux.HandleFunc("/user", s.Delete)
	s.mux.HandleFunc("/friends/{id:[0-9]+}", s.GetFriends) //regexp
	s.mux.HandleFunc("/{user_id:[0-9]+}", s.Put)           //regexp

	//my handler for debug
	s.mux.HandleFunc("/get_all", s.GetAll)
	s.mux.HandleFunc("/get/{user_id:[0-9]+}", s.GetUserInfo)
}

func (s *AppServer) configureStore() error {
	newStore := store.New(s.config.Store)
	if err := newStore.Open(); err != nil {
		return err
	}
	s.storeBD = newStore

	return nil
}
