package repo

import (
	"context"
	"fmt"
	"forum/internal/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx"
	v4 "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"strconv"
	"time"
)

type ForumRepo struct {
	db   *pgxpool.Pool
	info models.Status
}

func NewForumRepo(pool *pgxpool.Pool) *ForumRepo {
	return &ForumRepo{
		db: pool,
	}
}

func (r ForumRepo) CreateForum(forum models.Forum) (models.Forum, *models.InternalError) {
	user, code := r.CheckUser(models.User{NickName: forum.User})
	if code != models.OK {
		Error := &models.InternalError{
			Err: models.Error{
				Message: fmt.Sprintf("Can't find thread with id #\n"),
			},
			Code: models.ForumConflict,
		}
		return models.Forum{}, Error
	}

	forum.User = user.NickName
	//уменьшаем количество записей
	query := `INSERT INTO forums (title, author, slug) 
	VALUES ($1, $2, $3) 
	RETURNING slug;`

	//уменьшаем количество выделений памяти
	var results models.Forum
	res := forum
	row := r.db.QueryRow(context.Background(),
		query, forum.Title, forum.User, forum.Slug)
	err := row.Scan(&res.Slug)

	if err != nil {
		if pqError, ok := err.(*pgconn.PgError); ok {
			switch pqError.Code {
			case "23505":
				result, _ := r.GetForum(forum)
				Error := &models.InternalError{
					Err: models.Error{
						Message: fmt.Sprintf("Can't find thread with id #\n"),
					},
					Code: models.ForumConflict,
				}
				return result, Error
			case "23503":
				Error := &models.InternalError{
					Err: models.Error{
						Message: fmt.Sprintf("Can't find thread with id #\n"),
					},
					Code: models.NotFound,
				}
				return results, Error
			default:
				result, _ := r.GetForum(forum)
				Error := &models.InternalError{
					Err: models.Error{
						Message: fmt.Sprintf("Can't find thread with id #\n"),
					},
					Code: models.ForumConflict,
				}
				return result, Error
			}
		}
	}

	r.info.Forum++
	return forum, nil
}

func (r ForumRepo) GetForum(forum models.Forum) (models.Forum, *models.InternalError) {
	query := `SELECT title, author, slug, posts, threads
				FROM forums 
				WHERE slug = $1;`

	row := r.db.QueryRow(context.Background(), query, forum.Slug)

	err := row.Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)

	if err != nil {
		Error := &models.InternalError{
			Err: models.Error{
				Message: fmt.Sprintf("Can't find thread with id #\n"),
			},
			Code: models.NotFound,
		}
		return forum, Error
	}

	return forum, nil
}

func (r ForumRepo) CreateThread(forum models.Forum, thread models.Thread) (models.Thread, *models.InternalError) {
	user, code := r.CheckUser(models.User{NickName: thread.Author})
	if code != models.OK {
		Error := &models.InternalError{
			Err: models.Error{
				Message: fmt.Sprintf("Can't find thread with id #\n"),
			},
			Code: models.NotFound,
		}
		return models.Thread{}, Error
	}
	thread.Author = user.NickName

	f, status := r.ForumCheck(models.Forum{Slug: thread.Forum})
	if status == models.NotFound {
		Error := &models.InternalError{
			Err: models.Error{
				Message: fmt.Sprintf("Can't find thread with id #\n"),
			},
			Code: models.NotFound,
		}
		return models.Thread{}, Error
	}

	thread.Forum = f.Slug

	t := thread

	if thread.Slug != "" {
		thread, status := r.CheckSlug(thread)
		if status == models.OK {
			th, _ := r.GetThreadBySlug(thread.Slug, t)
			Error := &models.InternalError{
				Err: models.Error{
					Message: fmt.Sprintf("Can't find thread with id #\n"),
				},
				Code: models.ForumConflict,
			}
			return th, Error
		}
	}

	query := `INSERT INTO threads (author, message, title, created_at, forum, slug, votes)
				VALUES ($1, $2, $3, $4, $5, $6, $7)	RETURNING id`

	row := r.db.QueryRow(context.Background(), query, thread.Author, thread.Message, thread.Title,
		thread.Created, thread.Forum, thread.Slug, 0)

	err := row.Scan(&t.Id)

	if err != nil {
		if pqError, ok := err.(*pgconn.PgError); ok {
			switch pqError.Code {
			case "23505":
				Error := &models.InternalError{
					Err: models.Error{
						Message: fmt.Sprintf("Can't find thread with id #\n"),
					},
					Code: models.ForumConflict,
				}
				return t, Error
			case "23503":
				Error := &models.InternalError{
					Err: models.Error{
						Message: fmt.Sprintf("Can't find thread with id #\n"),
					},
					Code: models.NotFound,
				}
				return models.Thread{}, Error
			default:
				log.Print(pqError.Code)
				Error := &models.InternalError{
					Err: models.Error{
						Message: fmt.Sprintf("Can't find thread with id #\n"),
					},
					Code: models.NotFound,
				}
				return models.Thread{}, Error
			}
		}
	}
	r.info.Thread++
	query2 := `INSERT INTO forum_users(nickname,
    forum) VALUES ($1, $2) ON CONFLICT DO NOTHING;`
	_, err = r.db.Exec(context.Background(), query2, thread.Author, thread.Forum)
	if err != nil {
		log.Fatal(err)
	}

	query3 := `UPDATE forums SET threads = threads + 1 WHERE slug =$1`
	_, err = r.db.Exec(context.Background(), query3, thread.Forum)
	if err != nil {
		log.Fatal(err)
	}
	return t, nil
}

func (r *ForumRepo) CheckSlug(thread models.Thread) (models.Thread, models.ErrorCode) {
	query := `SELECT slug, author
				FROM threads 
				WHERE slug = $1;`

	row := r.db.QueryRow(context.Background(), query, thread.Slug)

	err := row.Scan(&thread.Slug, &thread.Author)

	if err != nil {
		return thread, models.NotFound
	}

	return thread, models.OK
}

func (r ForumRepo) GetUsers(forum models.Forum, limit, since, desc string) (models.Users, *models.InternalError) {
	us := []models.User{}
	var row v4.Rows
	var err error

	forum, state := r.ForumCheck(forum)
	if state != models.OK {
		Error := &models.InternalError{
			Err: models.Error{
				Message: fmt.Sprintf("Can't find thread with id #\n"),
			},
			Code: models.NotFound,
		}
		return us, Error
	}

	query := ``

	if limit == "" && since == "" {
		if desc == "true" {
			query = `SELECT email,
                      fullname,
                      users.nickname,
                      about
		FROM users JOIN forum_users ON users.nickname = forum_users.nickname WHERE forum = $1
		ORDER BY forum_users.nickname DESC`
		} else {
			query = `SELECT email,
                      fullname,
                      users.nickname,
                      about
		FROM users JOIN forum_users ON users.nickname = forum_users.nickname WHERE forum = $1
		ORDER BY forum_users.nickname ASC`
		}

		row, err = r.db.Query(context.Background(), query, forum.Slug)
	}

	if limit != "" && since == "" {
		if desc == "true" {
			query = `SELECT email,
                      fullname,
                      users.nickname,
                      about
		FROM users JOIN forum_users ON users.nickname = forum_users.nickname WHERE forum = $1
		ORDER BY forum_users.nickname DESC LIMIT $2`
		} else {
			query = `SELECT email,
                      fullname,
                      users.nickname,
                      about
		FROM users JOIN forum_users ON users.nickname = forum_users.nickname WHERE forum = $1
		ORDER BY forum_users.nickname ASC LIMIT $2`
		}

		row, err = r.db.Query(context.Background(), query, forum.Slug, limit)
	}

	if limit == "" && since != "" {
		if desc == "true" {
			query = `SELECT email,
                      fullname,
                      users.nickname,
                      about
		FROM users JOIN forum_users ON users.nickname = forum_users.nickname WHERE forum = $1
		AND forum_users.nickname < $2
		ORDER BY forum_users.nickname  DESC  `
		} else {
			query = `SELECT email,
                      fullname,
                     users.nickname,
                      about
		FROM users JOIN forum_users ON users.nickname = forum_users.nickname WHERE forum = $1
		AND forum_users.nickname > $2
		ORDER BY forum_users.nickname  ASC`
		}

		row, err = r.db.Query(context.Background(), query, forum.Slug, since)
	}

	if limit != "" && since != "" {
		if desc == "true" {
			query = `SELECT email,
                      fullname,
                      users.nickname,
                      about
		FROM users JOIN forum_users ON users.nickname = forum_users.nickname WHERE forum = $1
		AND forum_users.nickname < $2
		ORDER BY forum_users.nickname  DESC  LIMIT $3`
		} else {
			query = `SELECT email,
                      fullname,
                      users.nickname,
                      about
		FROM users JOIN forum_users ON users.nickname = forum_users.nickname WHERE forum = $1
		AND forum_users.nickname > $2
		ORDER BY forum_users.nickname ASC LIMIT $3`
		}

		row, err = r.db.Query(context.Background(), query, forum.Slug, since, limit)
		log.Print(err)
	}

	defer row.Close()

	for row.Next() {
		a := models.User{}
		row.Scan(&a.Email, &a.FullName, &a.NickName, &a.About)
		us = append(us, a)
	}
	return us, nil
}

func (r *ForumRepo) ForumCheck(forum models.Forum) (models.Forum, models.ErrorCode) {
	query := `SELECT slug FROM forums 
				WHERE slug = $1;`

	row := r.db.QueryRow(context.Background(), query, forum.Slug)

	err := row.Scan(&forum.Slug)

	if err != nil {
		return forum, models.NotFound
	}

	return forum, models.OK
}

func (r ForumRepo) GetThreads(forum models.Forum, limit, since, desc string) (models.Threads, *models.InternalError) {
	th := []models.Thread{}
	var row v4.Rows
	var err error

	query := ``

	if limit == "" && since == "" {
		if desc == "" || desc == "false" {
			query = `SELECT id, slug, author, created_at, forum, title, message, votes
						FROM threads
						WHERE forum = $1
						ORDER BY created_at ASC`

		} else {
			query = `SELECT id, slug, author, created_at, forum, title, message, votes
						FROM threads
						WHERE forum = $1
						ORDER BY created_at DESC`
		}
		row, err = r.db.Query(context.Background(), query, forum)
	} else {

		if limit != "" && since == "" {
			if desc == "" || desc == "false" {
				query = `SELECT id, slug, author, created_at, forum, title, message, votes
						FROM threads
						WHERE forum = $1
						ORDER BY created_at ASC  LIMIT $2`

			} else {
				query = `SELECT id, slug, author, created_at, forum, title, message, votes
						FROM threads
						WHERE forum = $1
						ORDER BY created_at DESC  LIMIT $2`
			}

			row, err = r.db.Query(context.Background(), query, forum, limit)
		}

		if since != "" && limit == "" {
			if desc == "" || desc == "false" {
				query = `SELECT id, slug, author, created_at, forum, title, message, votes
						FROM threads
						WHERE forum = $1 AND created_at >= $2
						ORDER BY created_at ASC `
			} else {
				query = `SELECT id, slug, author, created_at, forum, title, message, votes
						FROM threads
						WHERE forum = $1 AND created_at <= $2
						ORDER BY created_at DESC `
			}

			row, err = r.db.Query(context.Background(), query, forum, since)
		}

		if since != "" && limit != "" {

			if desc == "" || desc == "false" {

				query = `SELECT id, slug, author, created_at, forum, title, message, votes
						FROM threads
						WHERE forum = $1 AND created_at >= $2
						ORDER BY created_at ASC LIMIT $3`
			} else {
				query = `SELECT id, slug, author, created_at, forum, title, message, votes
						FROM threads
						WHERE forum = $1 AND created_at <= $2
						ORDER BY created_at DESC LIMIT $3`
			}
			row, err = r.db.Query(context.Background(), query, forum, since, limit)
		}
	}
	defer row.Close()
	for row.Next() {
		t := models.Thread{}
		err = row.Scan(&t.Id, &t.Slug, &t.Author, &t.Created, &forum, &t.Title, &t.Message, &t.Votes)

		th = append(th, t)
	}
	if err == nil {

	}

	if len(th) == 0 {
		_, status := r.GetForum(forum)
		if status != nil {
			Error := &models.InternalError{
				Err: models.Error{
					Message: fmt.Sprintf("Can't find thread with id #\n"),
				},
				Code: models.NotFound,
			}
			return th, Error
		}
		return th, nil
	}

	return th, nil
}

func (r ForumRepo) GetPostInfo(post models.Post, related []string) (models.PostFull, *models.InternalError) {
	pr := models.PostFull{
		Author: nil,
		Forum:  nil,
		Post:   models.Post{},
		Thread: nil,
	}

	p := models.Post{}
	p.Id = post.Id
	query := `SELECT author, post, created_at, forum, isedited, parent, thread
	FROM posts 
	WHERE id = $1`

	times := time.Time{}
	row := r.db.QueryRow(context.Background(), query, post.Id)
	err := row.Scan(&p.Author, &p.Message, &times,
		&p.Forum, &p.IsEdited, &p.Parent, &p.Thread)
	p.Created = times.Format(time.RFC3339)

	if err != nil {
		Error := &models.InternalError{
			Err: models.Error{
				Message: fmt.Sprintf("Can't find thread with id #\n"),
			},
			Code: models.NotFound,
		}
		return pr, Error
	}

	pr.Post = p

	for j := 0; j < len(related); j++ {
		if related[j] == "user" {
			u, _ := r.GetUser(models.User{NickName: p.Author})
			pr.Author = &u
		}
		if related[j] == "forum" {

			f, _ := r.GetForum(models.Forum{Slug: p.Forum})
			pr.Forum = &f

		}
		if related[j] == "thread" {
			t, _ := r.GetThreadByID(int(p.Thread), models.Thread{})
			pr.Thread = &t

		}
	}
	return pr, nil
}

func (r ForumRepo) UpdateMessage(post models.Post, update models.PostUpdate) (models.Post, *models.InternalError) {
	res := models.Post{}
	//проверить наличие поста
	query := `SELECT id, author, post, created_at,
                       forum, isEdited, parent, thread
				FROM posts 
				WHERE id = $1 `

	row := r.db.QueryRow(context.Background(), query, post.Id)

	times := time.Time{}
	err := row.Scan(&res.Id, &res.Author, &res.Message, &times,
		&res.Forum, &res.IsEdited, &res.Parent, &res.Thread)
	res.Created = times.Format(time.RFC3339)
	//поста нет
	if err != nil {
		Error := &models.InternalError{
			Err: models.Error{
				Message: fmt.Sprintf("Can't find thread with id #\n"),
			},
			Code: models.NotFound,
		}
		return models.Post{}, Error
	}

	if update.Message == "" || update.Message == res.Message {
		return res, nil
	}

	queryupdate := `UPDATE posts SET post = $1, isEdited = $2 WHERE id = $3`
	rowup, err := r.db.Exec(context.Background(), queryupdate, update.Message, true, post.Id)

	if err != nil || rowup.RowsAffected() == 0 {
		Error := &models.InternalError{
			Err: models.Error{
				Message: fmt.Sprintf("Can't find thread with id #\n"),
			},
			Code: models.NotFound,
		}
		return models.Post{}, Error
	}

	res.Message = update.Message
	res.IsEdited = true

	return res, nil
}

func (r ForumRepo) DropAllData() {
	query := `TRUNCATE TABLE users, forums, threads, post CASCADE;`
	r.db.Exec(context.Background(), query)

	r.info = models.Status{
		Forum:  0,
		Post:   0,
		Thread: 0,
		User:   0,
	}
}

func (r ForumRepo) GetStatus() models.Status {
	return r.info
}

func (r ForumRepo) CreatePosts(th models.Thread, posts models.Posts) (models.Posts, *models.InternalError) {
	thread := models.Thread{}

	query := `SELECT id, forum
					FROM threads
					WHERE slug = $1`

	row := r.db.QueryRow(context.Background(), query, th.Slug)
	err := row.Scan(&thread.Id, &thread.Forum)

	if err != nil {
		Error := &models.InternalError{
			Err: models.Error{
				Message: fmt.Sprintf("Can't find thread with id #%d\n", thread.Id),
			},
			Code: models.NotFound,
		}
		return []models.Post{}, Error
	}

	times := time.Now()

	if len(posts) == 0 {
		return posts, nil
	}

	tx, err := r.db.Begin(context.Background())
	query = `INSERT INTO posts (author, post, created_at, forum,  isEdited, parent, thread, path) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
			RETURNING  id;`
	ins, _ := tx.Prepare(context.Background(), "insert", query)

	result := []models.Post{}
	for _, p := range posts {
		p.Forum = thread.Forum
		p.Thread = thread.Id

		if p.Parent != 0 {
			old := 0
			query2 := `SELECT thread FROM posts WHERE id = $1`
			row = tx.QueryRow(context.Background(), query2, p.Parent)
			err := row.Scan(&old)
			if err != nil || old != int(p.Thread) {
				Error := &models.InternalError{
					Err: models.Error{
						Message: fmt.Sprintf("Can't find thread with id #%d\n", thread.Id),
					},
					Code: models.ForumConflict,
				}
				return []models.Post{}, Error
			}
		}

		err = tx.QueryRow(context.Background(), ins.Name, p.Author, p.Message, times, thread.Forum, false, p.Parent, thread.Id, []int{}).Scan(&p.Id)

		p.Created = times.Format(time.RFC3339)

		if err != nil {
			tx.Rollback(context.Background())
			if pqError, ok := err.(pgx.PgError); ok {
				switch pqError.Code {
				case "23503":
					Error := &models.InternalError{
						Err: models.Error{
							Message: fmt.Sprintf("Can't find thread with id #%d\n", thread.Id),
						},
						Code: models.NotFound,
					}
					return []models.Post{}, Error
				case "23505":
					Error := &models.InternalError{
						Err: models.Error{
							Message: fmt.Sprintf("Can't find thread with id #%d\n", thread.Id),
						},
						Code: models.ForumConflict,
					}
					return []models.Post{}, Error
				default:
					Error := &models.InternalError{
						Err: models.Error{
							Message: fmt.Sprintf("Can't find thread with id #%d\n", thread.Id),
						},
						Code: models.ForumConflict,
					}
					return []models.Post{}, Error
				}
			}
		}

		query2 := `INSERT INTO forum_users(nickname,
    	forum) VALUES ($1, $2) ON CONFLICT DO NOTHING;`
		_, err = r.db.Exec(context.Background(), query2, p.Author, p.Forum)
		if err != nil {
			Error := &models.InternalError{
				Err: models.Error{
					Message: fmt.Sprintf("Can't find thread with id #%d\n", thread.Id),
				},
				Code: models.NotFound,
			}
			return []models.Post{}, Error
		}
		result = append(result, p)
		r.info.Post++
	}
	tx.Commit(context.Background())
	query3 := `UPDATE forums SET posts = posts + $2 WHERE slug =$1`
	_, err = r.db.Exec(context.Background(), query3, thread.Forum, len(result))
	if err != nil {
		log.Fatal(err)
	}
	return result, nil
}

func (r ForumRepo) GetThreadInfoBySlug(thread models.Thread) (models.Thread, *models.InternalError) {
	th, status := r.GetThreadBySlug(thread.Slug, thread)
	if status == models.NotFound {
		Error := &models.InternalError{
			Err: models.Error{
				Message: fmt.Sprintf("Can't find thread with id #%d\n", thread.Id),
			},
			Code: models.NotFound,
		}
		return models.Thread{}, Error
	}
	return th, nil
}

func (r ForumRepo) UpdateThread(thread models.Thread, update models.ThreadUpdate) (models.Thread, *models.InternalError) {
	check := ""
	if thread.Slug != "" {
		check = thread.Slug
	} else {
		check = strconv.Itoa(int(thread.Id))
	}
	t, status := r.GetThreadBySlugOrId(check, thread)
	if status == models.NotFound {
		Error := &models.InternalError{
			Err: models.Error{
				Message: fmt.Sprintf("Can't find thread with id #%d\n", thread.Id),
			},
			Code: models.NotFound,
		}
		return thread, Error
	}

	if update.Message != "" {
		t.Message = update.Message
	}

	if update.Title != "" {
		t.Title = update.Title
	}

	query := `UPDATE threads
	SET message=$1, title=$2
	WHERE id = $3
	RETURNING id, author, message, title, created_at, forum, slug, votes`

	row := r.db.QueryRow(context.Background(), query, t.Message, t.Title, t.Id)
	res := models.Thread{}

	err := row.Scan(&res.Id, &res.Author, &res.Message, &res.Title,
		&res.Created, &res.Forum, &res.Slug, &res.Votes)
	if err == nil {
	}

	return res, nil
}

func (r *ForumRepo) GetThreadBySlugOrId(check string, thread models.Thread) (models.Thread, models.ErrorCode) {
	var row v4.Row

	if value, err := strconv.Atoi(check); err != nil {
		thread.Slug = check
		query := `SELECT id, author, message, title, created_at, forum, slug, votes
					FROM threads
					WHERE lower(slug) = lower($1)`
		row = r.db.QueryRow(context.Background(), query, thread.Slug)

	} else {
		query := `SELECT id, author, message, title, created_at, forum, slug, votes
					FROM threads
					WHERE id = $1`
		row = r.db.QueryRow(context.Background(), query, value)
	}

	err := row.Scan(&thread.Id, &thread.Author, &thread.Message, &thread.Title,
		&thread.Created, &thread.Forum, &thread.Slug, &thread.Votes)

	if err != nil {
		return thread, models.NotFound
	}

	return thread, models.OK
}

func (r ForumRepo) GetPosts(thread models.Thread, limit int32, since models.Post, sort string, desc bool) (models.Posts, *models.InternalError) {
	var row v4.Rows
	ps := models.Posts{}
	//TODO: получить только id
	thread, status := r.GetThreadBySlug(thread.Slug, models.Thread{})
	if status == models.NotFound {
		Error := &models.InternalError{
			Err: models.Error{
				Message: fmt.Sprintf("Can't find thread with id #%d\n", thread.Id),
			},
			Code: models.NotFound,
		}
		return ps, Error
	}

	switch sort {
	case "flat":
		row = r.getFlat(thread.Id, strconv.Itoa(int(since.Id)), strconv.Itoa(int(limit)), strconv.FormatBool(desc))

	case "tree":
		row = r.getTree(thread.Id, strconv.Itoa(int(since.Id)), strconv.Itoa(int(limit)), strconv.FormatBool(desc))

	case "parent_tree":
		row = r.getParentTree(thread.Id, strconv.Itoa(int(since.Id)), strconv.Itoa(int(limit)), strconv.FormatBool(desc))

	default:
		row = r.getFlat(thread.Id, strconv.Itoa(int(since.Id)), strconv.Itoa(int(limit)), strconv.FormatBool(desc))
	}

	defer row.Close()
	for row.Next() {

		pr := models.Post{}
		times := time.Time{}
		err := row.Scan(&pr.Id, &pr.Author, &pr.Message, &times, &pr.Forum, &pr.IsEdited, &pr.Parent, &pr.Thread)
		pr.Created = times.Format(time.RFC3339)
		if err != nil {
		}
		ps = append(ps, pr)
	}

	return ps, nil
}

func (r *ForumRepo) GetThreadBySlug(check string, thread models.Thread) (models.Thread, models.ErrorCode) {
	query := `SELECT id, author, message, title, created_at, forum, slug, votes
					FROM threads
					WHERE slug = $1`
	row := r.db.QueryRow(context.Background(), query, check)

	err := row.Scan(&thread.Id, &thread.Author, &thread.Message, &thread.Title,
		&thread.Created, &thread.Forum, &thread.Slug, &thread.Votes)

	if err != nil {
		return thread, models.NotFound
	}

	return thread, models.OK
}

func (r *ForumRepo) getFlat(id int32, since, limit, desc string) v4.Rows {
	var rows v4.Rows

	query := `SELECT id, author, post, created_at, forum, isedited, parent, thread
				FROM posts
				WHERE thread = $1`

	if limit == "" && since == "" {
		if desc == "true" {
			query += ` ORDER BY id DESC`
		} else {
			query += ` ORDER BY id ASC`
		}
		rows, _ = r.db.Query(context.Background(), query, id)
	} else {
		if limit != "" && since == "" {
			if desc == "true" {
				query += ` ORDER BY id DESC LIMIT $2`
			} else {
				query += `ORDER BY id ASC LIMIT $2`
			}
			rows, _ = r.db.Query(context.Background(), query, id, limit)
		}

		if limit != "" && since != "" {
			if desc == "true" {
				query += `AND id < $2 ORDER BY id DESC LIMIT $3`
			} else {
				query += `AND id > $2 ORDER BY id ASC LIMIT $3`
			}
			rows, _ = r.db.Query(context.Background(), query, id, since, limit)
		}

		if limit == "" && since != "" {
			if desc == "true" {
				query += `AND id < $2 ORDER BY id DESC`
			} else {
				query += `AND id > $2 ORDER BY id ASC`
			}
			rows, _ = r.db.Query(context.Background(), query, id, since)
		}
	}

	return rows
}

func (r *ForumRepo) getTree(id int32, since, limit, desc string) v4.Rows {

	var rows v4.Rows

	query := ``

	if limit == "" && since == "" {
		if desc == "true" {
			query = `SELECT id, author, post, created_at, forum, isedited, parent, thread
				FROM posts
				WHERE thread = $1 ORDER BY path, id DESC`
		} else {
			query = ` SELECT id, author, post, created_at, forum, isedited, parent, thread
				FROM posts
				WHERE thread = $1 ORDER BY path, id ASC`
		}
		rows, _ = r.db.Query(context.Background(), query, id)
	} else {
		if limit != "" && since == "" {
			if desc == "true" {
				query += `SELECT id, author, post, created_at, forum, isedited, parent, thread
				FROM posts
				WHERE thread = $1 ORDER BY path DESC, id DESC LIMIT $2`
			} else {
				query += `SELECT id, author, post, created_at, forum, isedited, parent, thread
				FROM posts
				WHERE thread = $1 ORDER BY path, id ASC LIMIT $2`
			}
			rows, _ = r.db.Query(context.Background(), query, id, limit)
		}

		if limit != "" && since != "" {
			if desc == "true" {
				query = `SELECT posts.id, posts.author, posts.post, 
				posts.created_at, posts.forum, posts.isedited, posts.parent, posts.thread
				FROM posts JOIN posts parent ON parent.id = $2 WHERE posts.path < parent.path AND  posts.thread = $1
				ORDER BY posts.path DESC, posts.id DESC LIMIT $3`
			} else {
				query = `SELECT posts.id, posts.author, posts.post, 
				posts.created_at, posts.forum, posts.isedited, posts.parent, posts.thread
				FROM posts JOIN posts parent ON parent.id = $2 WHERE posts.path > parent.path AND  posts.thread = $1
				ORDER BY posts.path ASC, posts.id ASC LIMIT $3`
			}
			rows, _ = r.db.Query(context.Background(), query, id, since, limit)
		}

		if limit == "" && since != "" {
			if desc == "true" {
				query = `SELECT posts.id, posts.author, posts.post, 
				posts.created_at, posts.forum, posts.isedited, posts.parent, posts.thread
				FROM posts JOIN posts parent ON parent.id = $2 WHERE posts.path < parent.path AND  posts.thread = $1
				ORDER BY posts.path DESC, posts.id DESC`
			} else {
				query = `SELECT posts.id, posts.author, posts.post, 
				posts.created_at, posts.forum, posts.isedited, posts.parent, posts.thread
				FROM posts JOIN posts parent ON parent.id = $2 WHERE posts.path > parent.path AND  posts.thread = $1
				ORDER BY posts.path ASC, posts.id ASC`
			}
			rows, _ = r.db.Query(context.Background(), query, id, since)
		}
	}

	return rows
}

func (r *ForumRepo) getParentTree(id int32, since, limit, desc string) v4.Rows {
	var rows v4.Rows

	//все корневые посты
	parents := fmt.Sprintf(`SELECT id FROM posts WHERE thread = %d AND parent = 0 `, id)

	if since != "" {
		if desc == "true" {
			parents += ` AND path[1] < ` + fmt.Sprintf(`(SELECT path[1] FROM posts WHERE id = %s) `, since)
		} else {
			parents += ` AND path[1] > ` + fmt.Sprintf(`(SELECT path[1] FROM posts WHERE id = %s) `, since)
		}
	}

	if desc == "true" {
		parents += ` ORDER BY id DESC `
	} else {
		parents += ` ORDER BY id ASC `
	}

	if limit != "" {
		parents += " LIMIT " + limit
	}

	query := fmt.Sprintf(
		`SELECT id, author, post, created_at, forum, isedited, parent, thread FROM posts WHERE path[1] = ANY (%s) `, parents)

	if desc == "true" {
		query += ` ORDER BY path[1] DESC, path,  id `
	} else {
		query += ` ORDER BY path[1] ASC, path,  id `
	}

	rows, _ = r.db.Query(context.Background(), query)
	return rows

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
