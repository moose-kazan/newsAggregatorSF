package dbaccess

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

type Comment struct {
	Id          int
	IdPost      int
	Content     string
	PubTime     int64
	FlagObscene bool
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

func (s *Store) Add(id_post int, content string) (Comment, error) {
	var rv Comment
	var sql string = `INSERT INTO comments
		(id_post, content, pubTime, flag_obscene) VALUES
		($1, $2, EXTRACT(EPOCH FROM NOW()) * 1000, false)
		RETURNING id, id_post, content, pubTime, flag_obscene`
	err := s.db.QueryRow(
		context.Background(),
		sql,
		id_post,
		content,
	).Scan(
		&rv.Id,
		&rv.IdPost,
		&rv.Content,
		&rv.PubTime,
		&rv.FlagObscene,
	)
	if err != nil {
		return rv, err
	}
	return rv, nil
}

func (s *Store) GetForPost(idPost int) ([]Comment, error) {
	var rv []Comment = make([]Comment, 0)
	var sql string = `SELECT
						id,
						id_post,
						content,
						pubTime,
						flag_obscene
					FROM comments
					WHERE id_post = $1
					ORDER BY pubTime`
	r, err := s.db.Query(context.Background(), sql, idPost)
	if err != nil {
		return rv, err
	}
	for r.Next() {
		var item Comment
		err = r.Scan(
			&item.Id,
			&item.IdPost,
			&item.Content,
			&item.PubTime,
			&item.FlagObscene,
		)

		if err != nil {
			return rv, err
		}
		rv = append(rv, item)
	}
	r.Close()
	return rv, nil
}

func (s *Store) Update(c Comment) error {
	return nil
}
