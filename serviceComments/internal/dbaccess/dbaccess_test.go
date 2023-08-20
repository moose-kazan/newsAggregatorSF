package dbaccess

import (
	"context"
	"testing"
)

const TestDSN = "postgresql://postgres:postgres@127.0.0.1:15432/test"

func TestStore_Add(t *testing.T) {
	t.Run("Main", func(t *testing.T) {
		s, err := New(TestDSN)
		if err != nil {
			t.Errorf("Can't connect to DB: %v", err)
			return
		}
		c1, err := s.Add(10, "test comment 1")
		if err != nil {
			t.Errorf("Add returned error: %v", err)
			return
		}
		c2, err := s.Add(10, "test comment 2")
		if err != nil {
			t.Errorf("Add returned error: %v", err)
			return
		}
		idDelta := c2.Id - c1.Id
		if idDelta != 1 {
			t.Errorf("id2 - id1: Expected: 1, Got: %v", idDelta)
		}
	})
}

func TestStore_GetForPost(t *testing.T) {
	t.Run("Main", func(t *testing.T) {
		s, err := New(TestDSN)
		if err != nil {
			t.Errorf("Can't connect to DB: %v", err)
			return
		}
		_, err = s.db.Exec(context.Background(), "DELETE FROM comments;")
		if err != nil {
			t.Errorf("Can't cleanup sources: %v", err)
			return
		}
		var sql string = `INSERT INTO comments
				(id, id_post, content, pubTime, flag_obscene)
			VALUES
				(1, 10, 'comment 1/10', 1692481427, false),
				(2, 20, 'comment 2/20', 1692481425, false),
				(3, 10, 'comment 3/10', 1692481429, false),
				(4, 10, 'comment 4/10', 1692481429, false),
				(5, 20, 'comment 5/20', 1692481439, false),
				(6, 30, 'comment 6/30', 1692481459, false),
				(7, 40, 'comment 7/40', 1692481469, false);`
		_, err = s.db.Exec(context.Background(), sql)
		if err != nil {
			t.Errorf("Can't insert comments: %v", err)
			return
		}
		comments, err := s.GetForPost(10)
		if err != nil {
			t.Errorf("Problem with comments: %v", err)
			return
		}
		if len(comments) != 3 {
			t.Errorf("Expected 3 comments, but found: %v", len(comments))
		}
		comments, err = s.GetForPost(50)
		if err != nil {
			t.Errorf("Problem with comments: %v", err)
			return
		}
		if len(comments) != 0 {
			t.Errorf("Expected 0 comments, but found: %v", len(comments))
		}
	})
}
