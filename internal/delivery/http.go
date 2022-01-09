package delivery

import (
	"encoding/json"
	"forum/internal"
	"forum/internal/models"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"net/http"
	"strconv"
	"strings"
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

	if forumErr != nil && forumErr.Code == models.ForumConflict {
		body, _ := easyjson.Marshal(newForum)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		w.Write(body)
		return
	} else if forumErr != nil {
		body, _ := easyjson.Marshal(forumErr.Err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(int(forumErr.Code))
		w.Write(body)
		return
	}

	body, err := easyjson.Marshal(newForum)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(body)
		return
	}

	body, err := easyjson.Marshal(newForum)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
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
	thread.Forum = slug

	th, forumErr := fu.usecase.CreateThread(*forum, thread)
	if forumErr != nil {
		switch forumErr.Code {
		case http.StatusNotFound:
			body, _ := easyjson.Marshal(forumErr.Err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write(body)
		case http.StatusConflict:
			body, _ := easyjson.Marshal(th)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			w.Write(body)
		}
		return
	}

	body, _ := easyjson.Marshal(th)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
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

	desc := r.URL.Query().Get("desc")
	limit := r.URL.Query().Get("limit")
	since := r.URL.Query().Get("since")

	users, userErr := fh.usecase.GetUsers(forum, limit, since, desc)
	if userErr != nil {
		body, _ := json.Marshal(userErr.Err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(body)
		return
	}

	body, _ := json.Marshal(users)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
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
	desc := r.URL.Query().Get("desc")
	limit := r.URL.Query().Get("limit")
	since := r.URL.Query().Get("since")

	threads, forumErr := fh.usecase.GetThreads(forum, limit, since, desc)
	if forumErr != nil {
		body, _ := easyjson.Marshal(forumErr.Err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(body)
		return
	}

	body, _ := json.Marshal(threads)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (fh ForumHandler) GetThreadInfo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(GetFromVars(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	post := models.Post{Id: int32(id)}

	query := r.URL.Query()
	relateds := query["related"]
	related := []string{}
	if len(relateds) > 0 {
		related = strings.Split(relateds[0], ",")
	}

	postFull, postErr := fh.usecase.GetPostInfo(post, related)
	if postErr != nil {
		body, _ := easyjson.Marshal(&postErr.Err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(body)
		return
	}
	body, _ := easyjson.Marshal(&postFull)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (fh ForumHandler) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(GetFromVars(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	post := models.Post{Id: int32(id)}
	postUpdate := models.PostUpdate{}
	err = easyjson.UnmarshalFromReader(r.Body, &postUpdate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newPost, intErr := fh.usecase.UpdateMessage(post, postUpdate)
	if intErr != nil {
		body, _ := easyjson.Marshal(&intErr.Err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(body)
		return
	}
	body, _ := json.Marshal(&newPost)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
	return
}

func (fh ForumHandler) DropAllInfo(w http.ResponseWriter, r *http.Request) {
	fh.usecase.DropAllInfo()
	w.WriteHeader(http.StatusOK)
}

func (fh ForumHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	st := fh.usecase.GetStatus()
	body, _ := easyjson.Marshal(st)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (fh ForumHandler) CreatePosts(w http.ResponseWriter, r *http.Request) {
	slugOrID := GetFromVars(r, "slug_or_id")
	if slugOrID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	posts := models.Posts{}
	err := json.NewDecoder(r.Body).Decode(&posts)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	intErr := &models.InternalError{}
	thread := models.Thread{}

	postsAnswer := models.Posts{}

	id, err := strconv.ParseInt(slugOrID, 10, 64)
	if err == nil {
		thread.Id = int32(id)
		postsAnswer, intErr = fh.usecase.CreatePostsID(thread, posts)
	} else {
		thread.Slug = slugOrID
		postsAnswer, intErr = fh.usecase.CreatePosts(thread, posts)
	}

	if intErr != nil {
		body, _ := easyjson.Marshal(intErr.Err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(int(intErr.Code))
		w.Write(body)
		return
	}

	body, _ := json.Marshal(postsAnswer)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func (fh ForumHandler) GetThreadInfoBySlug(w http.ResponseWriter, r *http.Request) {
	slug := GetFromVars(r, "slug_or_id")
	if slug == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	thread := models.Thread{
		Slug: slug,
	}

	th, intErr := fh.usecase.GetThreadInfoBySlug(thread)
	if intErr != nil {
		body, _ := easyjson.Marshal(intErr.Err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(body)
		return
	}

	body, _ := easyjson.Marshal(th)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (fh ForumHandler) UpdateThread(w http.ResponseWriter, r *http.Request) {
	slug := GetFromVars(r, "slug_or_id")
	if slug == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	thread := models.Thread{
		Slug: slug,
	}
	updateThread := models.ThreadUpdate{}
	err := easyjson.UnmarshalFromReader(r.Body, &updateThread)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	th, intErr := fh.usecase.UpdateThread(thread, updateThread)
	if intErr != nil {
		body, _ := easyjson.Marshal(intErr.Err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(body)
		return
	}

	body, _ := easyjson.Marshal(th)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (fh ForumHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	slug := GetFromVars(r, "slug_or_id")
	if slug == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	thread := models.Thread{
		Slug: slug,
	}

	desc, _ := strconv.ParseBool(r.URL.Query().Get("desc"))
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 32)
	since, err := strconv.ParseInt(r.URL.Query().Get("since"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sort := r.URL.Query().Get("sort")
	sincePost := models.Post{Id: int32(since)} // int64
	posts, intErr := fh.usecase.GetPosts(thread, int32(limit), sincePost, sort, desc)
	if intErr != nil {
		body, _ := easyjson.Marshal(intErr.Err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(body)
		return
	}

	body, _ := json.Marshal(posts)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (fh ForumHandler) VoteForThread(w http.ResponseWriter, r *http.Request) {
	slug := GetFromVars(r, "slug_or_id")
	if slug == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	thread := models.Thread{
		Slug: slug,
	}
	vote := models.Vote{}
	err := easyjson.UnmarshalFromReader(r.Body, &vote)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	th, intErr := fh.usecase.VoteForThread(thread, vote)
	if intErr != nil {
		body, _ := easyjson.Marshal(intErr.Err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(body)
		return
	}
	body, _ := easyjson.Marshal(th)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (fh ForumHandler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	nickname := GetFromVars(r, "nickname")
	if nickname == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := models.User{}
	err := easyjson.UnmarshalFromReader(r.Body, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user.NickName = nickname
	users, intErr := fh.usecase.CreateProfile(user)
	if intErr != nil {
		body, _ := json.Marshal(users)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		w.Write(body)
		return
	}
	body, _ := easyjson.Marshal(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func (fh ForumHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	nickname := GetFromVars(r, "nickname")
	if nickname == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := models.User{NickName: nickname}

	usr, intErr := fh.usecase.GetProfile(user)
	if intErr != nil {
		body, _ := easyjson.Marshal(intErr.Err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(body)
		return
	}
	body, _ := easyjson.Marshal(usr)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (fh ForumHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	nickname := GetFromVars(r, "nickname")
	if nickname == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := models.User{NickName: nickname}

	userUpdate := models.UserUpdate{}
	err := easyjson.UnmarshalFromReader(r.Body, &userUpdate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user.FullName = userUpdate.FullName
	user.About = userUpdate.About
	user.Email = userUpdate.Email

	usr, intErr := fh.usecase.UpdateProfile(user)
	if intErr != nil {
		body, _ := easyjson.Marshal(intErr.Err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(int(intErr.Code))
		w.Write(body)
		return
	}
	body, _ := easyjson.Marshal(usr)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
