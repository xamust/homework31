package proxy

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type AppProxy struct {
	config *Config
	mux    *mux.Router
	logger *logrus.Logger
}

func New(config *Config) *AppProxy {
	return &AppProxy{
		config: config,
		mux:    mux.NewRouter(),
		logger: logrus.New(),
	}
}

var COUNT int = 0

func (p *AppProxy) Start() error {
	if err := p.configureLogger(); err != nil {
		return err
	}
	p.configureRouter()
	p.logger.Info(fmt.Sprintf("Starting proxy (from %v to first instance on %v and second instance %v)...", p.config.BindAddr, p.config.FirstInst, p.config.SecondInst)) // set message Info level about succesfull starting proxy...
	return http.ListenAndServe(p.config.BindAddr, p.mux)                                                                                                                 //bind addr from config and new gorilla mux
}

func (p *AppProxy) configureLogger() error {
	level, err := logrus.ParseLevel(p.config.LogLevel)
	if err != nil {
		return err
	}
	p.logger.SetLevel(level)
	return nil
}

func (p *AppProxy) configureRouter() {

	p.mux.HandleFunc("/create", p.Create)
	p.mux.HandleFunc("/make_friends", p.MakeFriends)
	p.mux.HandleFunc("/user", p.Delete)
	p.mux.HandleFunc("/friends/{id:[0-9]+}", p.GetFriends) //regexp
	p.mux.HandleFunc("/{user_id:[0-9]+}", p.Put)           //regexp
	//my handler for debug
	p.mux.HandleFunc("/get_all", p.GetAll)
	p.mux.HandleFunc("/get/{user_id:[0-9]+}", p.GetUserInfo)
}
