package model

type User struct {
	ID      int64
	Name    string
	Age     int64
	Friends []int64 //The arguments to Scan must be of one of the supported types, or implement the sql.Scanner interface.
}
