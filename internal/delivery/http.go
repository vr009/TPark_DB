package delivery

import (
	"encoding/json"
	"forum/internal"
	"forum/internal/models"
	"github.com/gorilla/mux"
	"github.com/jackc/pgtype"
	"github.com/mailru/easyjson"
	"net/http"
	"strconv"
)

type ForumHandler struct {
	usecase internal.Usecase
}

func GetFromVars(r *http.Request, keyword string) string {
	vars := mux.Vars(r)
	res, found := vars[keyword]
	if !found {
		return ""
	}
	return res
}

func NewForumHandler(usecase internal.Usecase) *ForumHandler {
	return &ForumHandler{
		usecase: usecase,
	}
}

func (fh ForumHandler) CreateForum(w http.ResponseWriter, r *http.Request) {
	forum := models.Forum{}
	err := easyjson.UnmarshalFromReader(r.Body, &forum)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newForum, forumErr := fh.usecase.CreateForum(forum)
	if forumErr != nil {
		w.WriteHeader(int(forumErr.Code))
		return
	}

	body, err := easyjson.Marshal(&newForum)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (fh ForumHandler) GetForum(w http.ResponseWriter, r *http.Request) {
	slug := GetFromVars(r, "slug")
	if slug == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	forum := &models.Forum{
		Slug: slug,
	}

	newForum, forumErr := fh.usecase.GetForum(*forum)
	if forumErr != nil {
		body, _ := easyjson.Marshal(forumErr.Err)
		w.WriteHeader(http.StatusNotFound)
		w.Write(body)
		return
	}

	body, err := easyjson.Marshal(newForum)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (fu ForumHandler) CreateThread(w http.ResponseWriter, r *http.Request) {
	slug := GetFromVars(r, "slug")
	if slug == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	forum := &models.Forum{
		Slug: slug,
	}

	thread := models.Thread{}
	err := easyjson.UnmarshalFromReader(r.Body, &thread)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	th, forumErr := fu.usecase.CreateThread(*forum, thread)
	if forumErr != nil {
		switch forumErr.Code {
		case http.StatusNotFound:
			body, _ := easyjson.Marshal(forumErr.Err)
			w.WriteHeader(http.StatusNotFound)
			w.Write(body)
		case http.StatusConflict:
			body, _ := easyjson.Marshal(th)
			w.WriteHeader(http.StatusConflict)
			w.Write(body)
		}
		return
	}

	body, _ := easyjson.Marshal(th)
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (fh ForumHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	slug := GetFromVars(r, "slug")
	if slug == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	forum := models.Forum{
		Slug: slug,
	}

	desc, _ := strconv.ParseBool(r.URL.Query().Get("desc"))
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 32)
	since := r.URL.Query().Get("since")
	user := models.User{NickName: since}

	users, userErr := fh.usecase.GetUsers(forum, int32(limit), user, desc)
	if userErr != nil {
		body, _ := json.Marshal(userErr.Err)
		w.WriteHeader(http.StatusNotFound)
		w.Write(body)
		return
	}

	body, _ := json.Marshal(users)
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (fh ForumHandler) GetThreads(w http.ResponseWriter, r *http.Request) {
	slug := GetFromVars(r, "slug")
	if slug == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	forum := models.Forum{
		Slug: slug,
	}
	desc, _ := strconv.ParseBool(r.URL.Query().Get("desc"))
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 32)
	since := r.URL.Query().Get("since")

	time, _ := http.ParseTime(since)
	sinceDate := pgtype.Timestamptz{Time: time}

	threads, forumErr := fh.usecase.GetThreads(forum, int32(limit), sinceDate, desc)
	if forumErr != nil {
		body, _ := easyjson.Marshal(forumErr.Err)
		w.WriteHeader(http.StatusNotFound)
		w.Write(body)
		return
	}

	body, _ := json.Marshal(threads)
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (ForumHandler) GetThreadInfo(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (ForumHandler) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (ForumHandler) DropAllInfo(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (ForumHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (ForumHandler) CreatePosts(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (ForumHandler) GetThreadInfoBySlug(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (ForumHandler) UpdateThread(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (ForumHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (ForumHandler) VoteForThread(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (ForumHandler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (ForumHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (ForumHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
