package dbaccess

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

type SourceFetchUpdates string

const (
	FetchUpdatesOn  SourceFetchUpdates = "on"
	FetchUpdatesOff SourceFetchUpdates = "off"
)

type Source struct {
	Id              int
	Link            string
	FetchUpdates    SourceFetchUpdates
	DefaultInterval int
}

type Post struct {
	Id      int
	Source  int
	Title   string
	Content string
	PubTime int64
	Link    string
	Guid    string
}

func New(dsn string) (*Store, error) {
	rv := new(Store)
	var err error
	rv.db, err = pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return rv, nil
}

func (s *Store) SourceAdd(link string, defaultInterval int) (int, error) {
	var id int
	var sql string = `INSERT INTO sources
		(link, defaultInterval, fetchUpdates) VALUES
		($1, $2, 'on')
		RETURNING id`
	err := s.db.QueryRow(context.Background(), sql, link, defaultInterval).Scan(&id)
	if err == nil && id == 0 {
		return id, errors.New("Inserted, but ID is bad!")
	}
	return id, err
}

func (s *Store) SourceGetActive() ([]Source, error) {
	var rv []Source = make([]Source, 0)
	var sql string = `SELECT
						id,
						link,
						defaultInterval
					FROM sources
					WHERE fetchUpdates = 'on'`
	r, err := s.db.Query(context.Background(), sql)
	if err != nil {
		return rv, err
	}
	for r.Next() {
		var item Source
		err = r.Scan(
			&item.Id,
			&item.Link,
			&item.DefaultInterval,
		)

		if err != nil {
			return rv, err
		}
		rv = append(rv, item)
	}
	r.Close()
	return rv, nil
}

/*
 * Функция добавления поста в БД
 * Здесь такая логика:
 * error != nil - ошибка
 * id > 0 - успешная вставка
 * id == 0 - всё хорошо, но пытались вставить дубль
 */
func (s *Store) PostAdd(p Post) (int, error) {
	var sql string = `INSERT INTO posts
				(source, title, content, pubTime, link, guid)
			SELECT
				$1, $2, $3, $4, $5, $6
			WHERE NOT EXISTS
				(SELECT id FROM posts WHERE source = $1 AND guid = $6)
			RETURNING id`
	var id int
	rows, err := s.db.Query(
		context.Background(),
		sql,
		p.Source,
		p.Title,
		p.Content,
		p.PubTime,
		p.Link,
		p.Guid,
	)
	if err != nil {
		return 0, err
	}
	if rows.Next() {
		rows.Scan(&id)
	}
	rows.Close()
	return id, nil
}

func (s *Store) PostSearch(count int, offset int, searchQuery string) ([]Post, int, error) {
	rv := make([]Post, 0)
	var rvTc int
	var sql string = `SELECT
			id,
			source,
			title,
			content,
			pubTime,
			link,
			guid,
			COUNT(*) OVER() AS total_count
		FROM posts `
	if searchQuery != "" {
		sql += ` WHERE title ILIKE $3 `
	}
	sql += ` ORDER BY pubTime DESC
		LIMIT $2
		OFFSET $1;`
	var r pgx.Rows
	var err error
	if searchQuery != "" {
		query := "%" + searchQuery + "%"
		r, err = s.db.Query(context.Background(), sql, offset, count, query)
	} else {
		r, err = s.db.Query(context.Background(), sql, offset, count)
	}
	if err != nil {
		return rv, rvTc, err
	}
	for r.Next() {
		var p Post
		r.Scan(
			&p.Id,
			&p.Source,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
			&p.Guid,
			&rvTc,
		)
		rv = append(rv, p)
	}
	r.Close()
	return rv, rvTc, nil
}

func (s *Store) PostGetById(id int) (Post, error) {
	var rv Post
	var sql string = `SELECT
			id,
			source,
			title,
			content,
			pubTime,
			link,
			guid
		FROM posts
		WHERE id = $1`
	r, err := s.db.Query(context.Background(), sql, id)
	if err != nil {
		return rv, err
	}
	defer r.Close()
	if !r.Next() {
		return rv, errors.New("Post not found!")
	}
	r.Scan(
		&rv.Id,
		&rv.Source,
		&rv.Title,
		&rv.Content,
		&rv.PubTime,
		&rv.Link,
		&rv.Guid,
	)
	return rv, nil
}
