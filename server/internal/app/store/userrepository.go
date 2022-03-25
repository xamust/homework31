package store

import (
	"fmt"
	"github.com/lib/pq"
	"server/internal/app/model"
)

type UserRepositoryInterface interface {
	Create(m *model.User) (*model.User, error)
	UpdateUser(m *model.User) (*model.User, error)
	SetFriends(m *model.User) (*model.User, error)
	FindById(id int64) (*model.User, error)
	FindByName(name string) (*model.User, error)
	GetAll() (*[]model.User, error)
	DeleteByID(id int64) error
	ClearDeleteUserFromFriends(userID int64, id []int64) error
}

type UserRepository struct {
	store *AppStore
}

func (u *UserRepository) Create(m *model.User) (*model.User, error) { //returning id for postgres by default don't return id's
	if err := u.store.db.QueryRow( // after return result (id), Scan(&m.ID) mapping id to User.ID
		"INSERT INTO users (name,age,friends) VALUES ($1,$2,$3) RETURNING id",
		//"INSERT INTO users (name,age) VALUES ($1,$2) RETURNING id",
		m.Name,
		m.Age,
		pq.Array(m.Friends), //pq.Array (friends []int32!!!!!)  unsupported type []int, a slice of int, The arguments to Scan must be of one of the supported types, or implement the sql.Scanner interface.
	).Scan(&m.ID); err != nil {
		return nil, err
	}
	return m, nil
}

func (u *UserRepository) UpdateUser(m *model.User) (*model.User, error) {
	if err := u.store.db.QueryRow(
		// "INSERT INTO users (friends) VALUES ($1) RETURNING id",
		"UPDATE users SET name = $3, age = $4 ,friends = $1 WHERE id = $2 RETURNING id",
		pq.Array(m.Friends), //to array pq
		m.ID,
		m.Name,
		m.Age,
	).Scan(&m.ID); err != nil {
		return nil, err
	}
	return m, nil
}

func (u *UserRepository) SetFriends(m *model.User) (*model.User, error) {
	if err := u.store.db.QueryRow(
		// "INSERT INTO users (friends) VALUES ($1) RETURNING id",
		"UPDATE users SET friends = $1 WHERE id = $2 RETURNING id",
		pq.Array(m.Friends), //to array pq
		m.ID,
	).Scan(&m.ID); err != nil {
		return nil, err
	}
	return m, nil
}

func (u *UserRepository) FindById(id int64) (*model.User, error) {

	//another way to use query....
	query := "SELECT id,name,age,friends FROM users WHERE id = $1"
	us := &model.User{}

	var temp pq.Int64Array

	rows, err := u.store.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// iterate over the result and print out the titles
	for rows.Next() {
		if err := rows.Scan(&us.ID, &us.Name, &us.Age, &temp); err != nil {
			return nil, err
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	us.Friends = temp

	if us.ID == 0 {
		return nil, fmt.Errorf("Пользователя с id %d нет в базе данных", id)
	}

	return us, nil
}

func (u *UserRepository) FindByName(name string) (*model.User, error) {
	us := &model.User{}
	var temp pq.Int64Array
	if err := u.store.db.QueryRow(
		"SELECT id,name,age,friends FROM users WHERE name = $1",
		name,
	).Scan(
		&us.ID, &us.Name, &us.Age, &temp, //The arguments to Scan must be of one of the supported types, or implement the sql.Scanner interface.
	); err != nil {
		return nil, err
	}

	us.Friends = temp

	if us.ID == 0 {
		return nil, fmt.Errorf("Пользователя с именем %s нет в базе данных", name)
	}

	return us, nil
}

func (u *UserRepository) GetAll() (*[]model.User, error) {
	us := &model.User{}
	usRes := &[]model.User{}
	var temp pq.Int64Array
	rows, err := u.store.db.Query(
		"SELECT * FROM users",
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if err := rows.Scan(&us.ID, &us.Name, &us.Age, &temp); err != nil {
			return nil, err
		}
		us.Friends = temp
		*usRes = append(*usRes, *us)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return usRes, nil
}

func (u *UserRepository) DeleteByID(id int64) error {

	//GetFriendsUser...
	var temp pq.Int64Array
	if err := u.store.db.QueryRow("SELECT friends FROM users WHERE id = $1",
		id,
	).Scan(&temp); err != nil {
		return err
	}

	//DeleteUserFromFriends...
	if err := u.ClearDeleteUserFromFriends(id, temp); err != nil {
		return err
	}

	//DeleteUser...
	_, err := u.store.db.Exec("DELETE FROM users WHERE id = $1", id) //The Query() will return a sql.Rows, which reserves a database connection until the sql.Rows is closed.
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) ClearDeleteUserFromFriends(userID int64, id []int64) error {

	var temp pq.Int64Array
	us := &model.User{}

	for _, v := range id {
		if err := u.store.db.QueryRow("SELECT id,name,age,friends FROM users WHERE id = $1",
			v,
		).Scan(&us.ID, &us.Name, &us.Age, &temp); err != nil {
			return err
		}
		for i, k := range temp {
			if k == userID {
				temp = append(temp[:i], temp[i+1:]...)
			}
		}
		us.Friends = temp
		if _, err := u.SetFriends(us); err != nil {
			return err
		}
	}
	return nil
}
