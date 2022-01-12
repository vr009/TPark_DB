package models

type User struct {
	ID       int64  `json:"ID,omitempty"`
	NickName string `json:"nickname"`
	FullName string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}
