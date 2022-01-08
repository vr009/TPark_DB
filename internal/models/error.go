package models

type Error struct {
	Message string `json:"message"`
}

type ErrorCode int

const (
	OK ErrorCode = iota
	ForumConflict
	NotFound
)

type InternalError struct {
	Err  Error
	Code ErrorCode
}
