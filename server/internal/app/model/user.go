package model

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
