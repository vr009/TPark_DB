package repo

import "forum/internal/models"

type ForumRepo struct{}

func (ForumRepo) CreateForum(forum models.Forum) (models.Forum, models.Error) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) GetForum(forum models.Forum) (models.Forum, models.Error) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) CreateThread(forum models.Forum, thread models.Thread) (models.Thread, models.Error) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) GetUsers(forum models.Forum, limit int32, sinceUser models.User, desc bool) (models.Users, models.Error) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) GetThreads(forum models.Forum, limit int32, sinceUser models.User, desc bool) (models.Threads, models.Error) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) GetThreadInfo(thread models.Thread, related interface{}) (models.Thread, models.Error) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) UpdateMessage(post models.Post, update models.PostUpdate) (models.Post, models.Error) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) DropAllData() models.Error {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) GetStatus() (models.Status, models.Error) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) CreatePosts(thread models.Thread, posts models.Posts) (models.Posts, models.Error) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) GetThreadInfoBySlug(thread models.Thread) (models.Thread, models.Error) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) UpdateThread(thread models.Thread, update models.ThreadUpdate) (models.Thread, models.Error) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) GetPosts(forum models.Forum, limit int32, sinceUser models.User, desc bool) (models.Posts, models.Error) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) VoteForThread(thread models.Thread, vote models.Vote) (models.Thread, models.Error) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) CreateUser(user models.User) (models.User, models.Error) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) GetUser(user models.User) (models.User, models.Error) {
	//TODO implement me
	panic("implement me")
}

func (ForumRepo) UpdateUser(user models.User) (models.User, models.Error) {
	//TODO implement me
	panic("implement me")
}
