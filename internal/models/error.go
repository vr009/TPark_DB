package models

type Error struct {
	Message string `json:"message"`
}

type ErrorCode int

const (
	AlreadyExists ErrorCode = iota
	NotFound
)

type InternalError struct {
	Err  Error
	Code ErrorCode
}
