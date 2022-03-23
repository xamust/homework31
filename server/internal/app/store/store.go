package store

import (
	"database/sql"
	_ "github.com/lib/pq" //anonym import
)

type AppStore struct {
	config         *Config
	db             *sql.DB
	userRepository *UserRepository
}

func New(config *Config) *AppStore {
	return &AppStore{
		config: config,
	}
}

func (s *AppStore) Open() error {
	db, err := sql.Open("postgres", s.config.DataBaseUrl)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil { //ping connect DB
		return err
	}

	s.db = db

	return nil
}

func (s *AppStore) Close() {
	s.db.Close()
}

//store.User().Create
func (s *AppStore) User() *UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
