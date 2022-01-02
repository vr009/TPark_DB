package usecase

import (
	"forum/internal"
	"forum/internal/models"
)

type ForumUsecase struct {
	repo internal.Repo
}

func NewForumUsecase(repo *internal.Repo) *ForumUsecase {
	return &ForumUsecase{
		repo: *repo,
	}
}

func (fu ForumUsecase) CreateForum(forum models.Forum) (models.Forum, models.Error) {
	return fu.repo.GetForum(forum)
}

func (fu ForumUsecase) GetForum(forum models.Forum) (models.Forum, models.Error) {
	return fu.repo.GetForum(forum)
}

func (fu ForumUsecase) CreateThread(forum models.Forum, thread models.Thread) (models.Thread, models.Error) {
	return fu.repo.CreateThread(forum, thread)
}

func (fu ForumUsecase) GetUsers(forum models.Forum, limit int32, sinceUser models.User, desc bool) (models.Users, models.Error) {
	return fu.repo.GetUsers(forum, limit, sinceUser, desc)
}

func (fu ForumUsecase) GetThreads(forum models.Forum, limit int32, sinceUser models.User, desc bool) (models.Threads, models.Error) {
	return fu.repo.GetThreads(forum, limit, sinceUser, desc)
}

func (fu ForumUsecase) GetThreadInfo(thread models.Thread, related interface{}) (models.Thread, models.Error) {
	return fu.repo.GetThreadInfo(thread, related)
}

func (fu ForumUsecase) UpdateMessage(post models.Post, update models.PostUpdate) (models.Post, models.Error) {
	return fu.repo.UpdateMessage(post, update)
}

func (fu ForumUsecase) DropAllInfo() models.Error {
	return fu.repo.DropAllData()
}

func (fu ForumUsecase) GetStatus() (models.Status, models.Error) {
	return fu.repo.GetStatus()
}

func (fu ForumUsecase) CreatePosts(thread models.Thread, posts models.Posts) (models.Posts, models.Error) {
	return fu.repo.CreatePosts(thread, posts)
}

func (fu ForumUsecase) GetThreadInfoBySlug(thread models.Thread) (models.Thread, models.Error) {
	return fu.repo.GetThreadInfoBySlug(thread)
}

func (fu ForumUsecase) UpdateThread(thread models.Thread, update models.ThreadUpdate) (models.Thread, models.Error) {
	return fu.repo.UpdateThread(thread, update)
}

func (fu ForumUsecase) GetPosts(forum models.Forum, limit int32, sinceUser models.User, desc bool) (models.Posts, models.Error) {
	return fu.repo.GetPosts(forum, limit, sinceUser, desc)
}

func (fu ForumUsecase) VoteForThread(thread models.Thread, vote models.Vote) (models.Thread, models.Error) {
	return fu.repo.VoteForThread(thread, vote)
}

func (fu ForumUsecase) CreateProfile(user models.User) (models.User, models.Error) {
	return fu.CreateProfile(user)
}

func (fu ForumUsecase) GetProfile(user models.User) (models.User, models.Error) {
	return fu.repo.GetUser(user)
}

func (fu ForumUsecase) UpdateProfile(user models.User) (models.User, models.Error) {
	return fu.repo.UpdateUser(user)
}
