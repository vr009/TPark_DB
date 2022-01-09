package models

type Error struct {
	Message string `json:"message"`
}

type ErrorCode int

const (
	OK            ErrorCode = 200
	ForumConflict           = 409
	NotFound                = 404
)

type InternalError struct {
	Err  Error
	Code ErrorCode
}
