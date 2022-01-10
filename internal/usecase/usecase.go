package usecase

import (
	"forum/internal"
	"forum/internal/models"
)

type ForumUsecase struct {
	repo internal.Repo
}

func NewForumUsecase(repo internal.Repo) *ForumUsecase {
	return &ForumUsecase{
		repo: repo,
	}
}

func (fu ForumUsecase) CreateForum(forum models.Forum) (models.Forum, *models.InternalError) {
	return fu.repo.CreateForum(forum)
}

func (fu ForumUsecase) GetForum(forum models.Forum) (models.Forum, *models.InternalError) {
	return fu.repo.GetForum(forum)
}

func (fu ForumUsecase) CreateThread(forum models.Forum, thread models.Thread) (models.Thread, *models.InternalError) {
	return fu.repo.CreateThread(forum, thread)
}

func (fu ForumUsecase) GetUsers(forum models.Forum, limit, since, desc string) (models.Users, *models.InternalError) {
	return fu.repo.GetUsers(forum, limit, since, desc)
}

func (fu ForumUsecase) GetThreads(forum models.Forum, limit, since, desc string) (models.Threads, *models.InternalError) {
	return fu.repo.GetThreads(forum, limit, since, desc)
}

func (fu ForumUsecase) GetPostInfo(post models.Post, related []string) (models.PostFull, *models.InternalError) {
	return fu.repo.GetPostInfo(post, related)
}

func (fu ForumUsecase) UpdateMessage(post models.Post, update models.PostUpdate) (models.Post, *models.InternalError) {
	return fu.repo.UpdateMessage(post, update)
}

func (fu ForumUsecase) DropAllInfo() {
	fu.repo.DropAllData()
}

func (fu ForumUsecase) GetStatus() models.Status {
	return fu.repo.GetStatus()
}

func (fu ForumUsecase) CreatePosts(thread models.Thread, posts models.Posts) (models.Posts, *models.InternalError) {
	return fu.repo.CreatePosts(thread, posts)
}

func (fu ForumUsecase) CreatePostsID(thread models.Thread, posts models.Posts) (models.Posts, *models.InternalError) {
	return fu.repo.CreatePostsID(thread, posts)
}

func (fu ForumUsecase) GetThreadInfoBySlug(thread models.Thread) (models.Thread, *models.InternalError) {
	return fu.repo.GetThreadInfoBySlug(thread)
}

func (fu ForumUsecase) GetThreadInfoByID(thread models.Thread) (models.Thread, *models.InternalError) {
	return fu.repo.GetThreadInfoByID(thread)
}

func (fu ForumUsecase) UpdateThread(thread models.Thread, update models.ThreadUpdate) (models.Thread, *models.InternalError) {
	return fu.repo.UpdateThread(thread, update)
}

func (fu ForumUsecase) GetPosts(thread models.Thread, limit, sincePost, sort, desc string) (models.Posts, *models.InternalError) {
	return fu.repo.GetPosts(thread, limit, sincePost, sort, desc)
}

func (fu ForumUsecase) GetPostsID(thread models.Thread, limit, sincePost, sort, desc string) (models.Posts, *models.InternalError) {
	return fu.repo.GetPostsID(thread, limit, sincePost, sort, desc)
}

func (fu ForumUsecase) VoteForThread(thread models.Thread, vote models.Vote) (models.Thread, *models.InternalError) {
	return fu.repo.VoteForThread(thread, vote)
}

func (fu ForumUsecase) VoteForThreadID(thread models.Thread, vote models.Vote) (models.Thread, *models.InternalError) {
	return fu.repo.VoteForThreadID(thread, vote)
}

func (fu ForumUsecase) CreateProfile(user models.User) ([]models.User, *models.InternalError) {
	return fu.repo.CreateUser(user)
}

func (fu ForumUsecase) GetProfile(user models.User) (models.User, *models.InternalError) {
	return fu.repo.GetUser(user)
}

func (fu ForumUsecase) UpdateProfile(user models.User) (models.User, *models.InternalError) {
	return fu.repo.UpdateUser(user)
}
