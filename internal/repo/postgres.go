package repo

import (
	"forum/internal/models"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgtype/pgxtype"
)

type ForumRepo struct {
	pool pgxtype.Querier
}

func NewForumRepo(pool pgxtype.Querier) *ForumRepo {
	return &ForumRepo{
		pool: pool,
	}
}

func (ForumRepo) CreateForum(forum models.Forum) (models.Forum, *models.InternalError) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) GetForum(forum models.Forum) (models.Forum, *models.InternalError) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) CreateThread(forum models.Forum, thread models.Thread) (models.Thread, *models.InternalError) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) GetUsers(forum models.Forum, limit int32, sinceUser models.User, desc bool) (models.Users, *models.InternalError) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) GetThreads(forum models.Forum, limit int32, since pgtype.Timestamptz, desc bool) (models.Threads, *models.InternalError) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) GetPostInfo(post models.Post, related string) (models.PostFull, *models.InternalError) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) UpdateMessage(post models.Post, update models.PostUpdate) (models.Post, *models.InternalError) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) DropAllData() {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) GetStatus() models.Status {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) CreatePosts(thread models.Thread, posts models.Posts) (models.Posts, *models.InternalError) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) GetThreadInfoBySlug(thread models.Thread) (models.Thread, *models.InternalError) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) UpdateThread(thread models.Thread, update models.ThreadUpdate) (models.Thread, *models.InternalError) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) GetPosts(thread models.Thread, limit int32, sincePost models.Post, sort string, desc bool) (models.Posts, *models.InternalError) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) VoteForThread(thread models.Thread, vote models.Vote) (models.Thread, *models.InternalError) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) CreateUser(user models.User) (models.User, *models.InternalError) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) GetUser(user models.User) (models.User, *models.InternalError) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) UpdateUser(user models.User) (models.User, *models.InternalError) {
	//TODO implement me
	panic("implement me")
}
