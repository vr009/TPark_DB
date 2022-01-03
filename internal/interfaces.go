package internal

import (
	"forum/internal/models"
	"github.com/jackc/pgtype"
)

type Repo interface {
	CreateForum(forum models.Forum) (models.Forum, *models.InternalError)
	GetForum(forum models.Forum) (models.Forum, *models.InternalError)
	CreateThread(forum models.Forum, thread models.Thread) (models.Thread, *models.InternalError)
	GetUsers(forum models.Forum, limit int32, sinceUser models.User, desc bool) (models.Users, *models.InternalError)
	GetThreads(forum models.Forum, limit int32, since pgtype.Timestamptz, desc bool) (models.Threads, *models.InternalError)
	GetPostInfo(post models.Post, related string) (models.PostFull, *models.InternalError)
	UpdateMessage(post models.Post, update models.PostUpdate) (models.Post, *models.InternalError)
	DropAllData()
	GetStatus() models.Status
	CreatePosts(thread models.Thread, posts models.Posts) (models.Posts, *models.InternalError)
	GetThreadInfoBySlug(thread models.Thread) (models.Thread, *models.InternalError)
	UpdateThread(thread models.Thread, update models.ThreadUpdate) (models.Thread, *models.InternalError)
	GetPosts(forum models.Thread, limit int32, sincePost models.Post, sort string, desc bool) (models.Posts, *models.InternalError)
	VoteForThread(thread models.Thread, vote models.Vote) (models.Thread, *models.InternalError)
	CreateUser(user models.User) (models.User, *models.InternalError)
	GetUser(user models.User) (models.User, *models.InternalError)
	UpdateUser(user models.User) (models.User, *models.InternalError)
}

type Usecase interface {
	CreateForum(forum models.Forum) (models.Forum, *models.InternalError)
	GetForum(forum models.Forum) (models.Forum, *models.InternalError)
	CreateThread(forum models.Forum, thread models.Thread) (models.Thread, *models.InternalError)
	GetUsers(forum models.Forum, limit int32, sinceUser models.User, desc bool) (models.Users, *models.InternalError)
	GetThreads(forum models.Forum, limit int32, since pgtype.Timestamptz, desc bool) (models.Threads, *models.InternalError)
	GetPostInfo(thread models.Post, related string) (models.PostFull, *models.InternalError)
	UpdateMessage(post models.Post, update models.PostUpdate) (models.Post, *models.InternalError)
	DropAllInfo()
	GetStatus() models.Status
	CreatePosts(thread models.Thread, posts models.Posts) (models.Posts, *models.InternalError)
	GetThreadInfoBySlug(thread models.Thread) (models.Thread, *models.InternalError)
	UpdateThread(thread models.Thread, update models.ThreadUpdate) (models.Thread, *models.InternalError)
	GetPosts(forum models.Thread, limit int32, sincePost models.Post, sort string, desc bool) (models.Posts, *models.InternalError)
	VoteForThread(thread models.Thread, vote models.Vote) (models.Thread, *models.InternalError)
	CreateProfile(user models.User) (models.User, *models.InternalError)
	GetProfile(user models.User) (models.User, *models.InternalError)
	UpdateProfile(user models.User) (models.User, *models.InternalError)
}
