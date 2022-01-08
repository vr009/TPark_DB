package repo

import (
	"context"
	"fmt"
	"forum/internal/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgtype/pgxtype"
	v4 "github.com/jackc/pgx/v4"
	"log"
)

type ForumRepo struct {
	db   pgxtype.Querier
	info models.Status
}

func NewForumRepo(pool pgxtype.Querier) *ForumRepo {
	return &ForumRepo{
		db: pool,
	}
}

func (r ForumRepo) CreateForum(forum models.Forum) (models.Forum, *models.InternalError) {
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

func (r ForumRepo) VoteForThread(thread models.Thread, vote models.Vote) (models.Thread, *models.InternalError) {
	user, code := r.CheckUser(models.User{NickName: vote.NickName})
	if code != models.OK {
		Error := &models.InternalError{
			Err: models.Error{Message: fmt.Sprintf("Can't find user with id #%d\n", user.ID)},
		}
		return models.Thread{}, Error
	}
	vote.NickName = user.NickName

	thread, status := r.GetSlugID(thread.Slug, models.Thread{})
	if status == models.NotFound {
		Error := &models.InternalError{
			Err: models.Error{Message: fmt.Sprintf("Can't find user with id #%d\n", user.ID)},
		}
		return thread, Error
	}

	query := `INSERT INTO VOTES (author, vote, thread) VALUES ($1, $2, $3) RETURNING *`
	row := r.db.QueryRow(context.Background(), query, vote.NickName, vote.Voice, thread.Id)

	value := 0

	err := row.Scan(&vote.NickName, &vote.Voice, &thread.Id)
	if err != nil {
		if pqError, ok := err.(*pgconn.PgError); ok {
			switch pqError.Code {
			case "23503":
				Error := &models.InternalError{
					Err: models.Error{Message: fmt.Sprintf("Can't find user with id #%d\n", user.ID)},
				}
				return thread, Error
			case "23505":
				upd := "WITH u AS ( SELECT vote FROM votes WHERE author = $2 AND thread = $3)" +
					"UPDATE votes SET vote =  $1 WHERE author = $2 AND thread = $3 " +
					"RETURNING vote, (SELECT vote FROM u)"
				row := r.db.QueryRow(context.Background(), upd, vote.Voice, vote.NickName, thread.Id)
				err := row.Scan(&vote.Voice, &value)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	query = `UPDATE threads SET votes=votes+$1 WHERE id = $2;`
	_, err = r.db.Exec(context.Background(), query, vote.Voice-int32(value), thread.Id)

	thread, status = r.GetThreadByID(int(thread.Id), models.Thread{})
	return thread, nil
}

func (r *ForumRepo) GetThreadByID(id int, thread models.Thread) (models.Thread, models.ErrorCode) {
	var row v4.Row

	query := `SELECT id, author, message, title, created_at, forum, slug, votes
					FROM threads
					WHERE id = $1`
	row = r.db.QueryRow(context.Background(), query, id)

	err := row.Scan(&thread.Id, &thread.Author, &thread.Message, &thread.Title,
		&thread.Created, &thread.Forum, &thread.Slug, &thread.Votes)

	if err != nil {
		return thread, models.NotFound
	}

	return thread, models.OK
}

func (r *ForumRepo) GetSlugID(check string, thread models.Thread) (models.Thread, models.ErrorCode) {
	var row v4.Row

	query := `SELECT id FROM threads WHERE lower(slug) = lower($1)`
	row = r.db.QueryRow(context.Background(), query, check)

	err := row.Scan(&thread.Id)

	if err != nil {
		return thread, models.NotFound
	}

	return thread, models.OK
}

func (r ForumRepo) CreateUser(user models.User) ([]models.User, *models.InternalError) {
	results := []models.User{}
	result := user

	query := `INSERT INTO users (email, fullname, nickname, about) 
			VALUES ($1, $2, $3, $4) RETURNING nickname`

	_, err := r.db.Exec(context.Background(), query, user.Email, user.FullName, user.NickName, user.About)
	if err != nil {
		if pqError, ok := err.(*pgconn.PgError); ok {
			switch pqError.Code {
			case "23505":
				us, _ := r.GetUserOnConflict(user)
				return us, &models.InternalError{Code: models.ForumConflict}
			default:
				us, _ := r.GetUserOnConflict(user)
				return us, &models.InternalError{Code: models.ForumConflict}
			}
		}
	}

	r.info.User++
	results = append(results, result)
	return results, nil
}

func (r ForumRepo) GetUser(user models.User) (models.User, *models.InternalError) {
	result := models.User{}
	query := `SELECT email, fullname, nickname, about
	FROM users
	WHERE  nickname =  $1`

	rows := r.db.QueryRow(context.Background(), query, user.NickName)

	err := rows.Scan(&result.Email, &result.FullName, &result.NickName, &result.About)
	if err != nil {
		errorCode := &models.InternalError{
			Err: models.Error{
				Message: fmt.Sprintf("Can't find user with id #%d\n", user.ID),
			},
			Code: models.NotFound,
		}
		return result, errorCode
	}
	return result, nil
}

func (r ForumRepo) UpdateUser(user models.User) (models.User, *models.InternalError) {
	us, status := r.GetUser(user)
	if status.Code == models.NotFound {
		return us, status
	}

	if user.FullName != "" {
		us.FullName = user.FullName
	}
	if user.Email != "" {
		us.Email = user.Email
	}
	if user.About != "" {
		us.About = user.About
	}

	query := `UPDATE users 
	SET fullname=$1, email=$2, about=$3 
	WHERE nickname = $4 
	RETURNING nickname, fullname, about, email;`

	rows := r.db.QueryRow(context.Background(), query, us.FullName, us.Email, us.About, us.NickName)
	err := rows.Scan(&us.NickName, &us.FullName, &us.About, &us.Email)

	if err != nil {
		Error := &models.InternalError{
			Err: models.Error{
				Message: fmt.Sprintf("Can't find user with id #%s\n", user.ID),
			},
		}
		if pqError, ok := err.(*pgconn.PgError); ok {
			switch pqError.Code {
			case "23505":
				Error.Code = models.ForumConflict
				return us, Error
			case "23503":
				Error.Code = models.NotFound
				return us, Error
			}
		}
	}

	return us, nil
}

func (r *ForumRepo) CheckUser(user models.User) (models.User, models.ErrorCode) {
	result := models.User{}
	query := `SELECT nickname
	FROM users
	WHERE nickname =  $1`

	rows := r.db.QueryRow(context.Background(), query, user.NickName)

	err := rows.Scan(&result.NickName)
	if err != nil {
		return result, models.NotFound
	}
	return result, models.OK
}

func (r *ForumRepo) GetUserOnConflict(user models.User) ([]models.User, models.ErrorCode) {
	results := []models.User{}
	query := `SELECT email, fullname, nickname, about
	FROM users
	WHERE email = $1 or nickname =  $2`

	rows, _ := r.db.Query(context.Background(), query, user.Email, user.NickName)
	defer rows.Close()

	for rows.Next() {
		result := models.User{}
		rows.Scan(&result.Email, &result.FullName, &result.NickName, &result.About)
		results = append(results, result)
	}
	return results, models.OK
}
