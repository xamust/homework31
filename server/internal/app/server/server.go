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
	storeBD *store.AppStore
	handl   Handlers
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

	//handlers init...
	s.handl = Handlers{s.logger, s.mux, s.storeBD.User()}

	s.logger.Info(fmt.Sprintf("Starting server (bind on %v)...", s.config.BindAddr)) // set message Info level about succesfull starting server...
	return http.ListenAndServe(s.config.BindAddr, s.mux)                             //bind addr from config and new gorilla mux
}

//config route...
func (s *AppServer) configureRouter() {

	s.mux.HandleFunc("/create", s.handl.Create)
	s.mux.HandleFunc("/make_friends", s.handl.MakeFriends)
	s.mux.HandleFunc("/user", s.handl.Delete)
	s.mux.HandleFunc("/friends/{id:[0-9]+}", s.handl.GetFriends) //regexp
	s.mux.HandleFunc("/{user_id:[0-9]+}", s.handl.Put)           //regexp

	//my handler for debug
	s.mux.HandleFunc("/get_all", s.handl.GetAll)
	s.mux.HandleFunc("/get/{user_id:[0-9]+}", s.handl.GetUserInfo)
}

func (s *AppServer) configureStore() error {
	newStore := store.New(s.config.Store)
	if err := newStore.Open(); err != nil {
		return err
	}
	s.storeBD = newStore

	return nil
}
