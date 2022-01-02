package models

type PostFull struct {
	Post   Post   `json:"post"`
	Author User   `json:"author"`
	Thread Thread `json:"thread"`
	Forum  Forum  `json:"forum"`
}
