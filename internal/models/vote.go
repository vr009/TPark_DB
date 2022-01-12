package models

const (
	Like    = 1
	DisLike = -1
)

type Vote struct {
	NickName string `json:"nickname"`
	Voice    int32  `json:"voice"`
	Existed  bool
}
