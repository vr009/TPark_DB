package models

type Post struct {
	Id       int32  `json:"id"`
	Parent   int32  `json:"parent"`
	Author   string `json:"author"`
	Message  string `json:"message"`
	IsEdited bool   `json:"isEdited"`
	Forum    string `json:"forum"`
	Thread   int32  `json:"thread"`
	Created  string `json:"created"`
}
