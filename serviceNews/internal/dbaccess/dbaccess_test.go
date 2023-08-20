package dbaccess

import (
	"context"
	"testing"
	"time"
)

const TestDSN = "postgresql://postgres:postgres@127.0.0.1:15432/test"

func TestStore_SourceAdd(t *testing.T) {
	t.Run("SourceAdd", func(t *testing.T) {
		s, err := New(TestDSN)
		if err != nil {
			t.Errorf("Can't connect to DB: %v", err)
			return
		}
		_, err = s.db.Exec(context.Background(), "DELETE FROM sources;")
		if err != nil {
			t.Errorf("Can't cleanup sources: %v", err)
			return
		}
		id, err := s.SourceAdd("https://example.com/feed", 30)
		if err != nil {
			t.Errorf("SourceAdd returned error: %v", err)
			return
		}
		// На чистой базе первая запись должна иметь id=1
		if id != 1 {
			t.Errorf("SourceAdd expected 1 got: %v", id)
		}
	})

}

func TestStore_SourceGetActive(t *testing.T) {
	t.Run("SourceGetActive", func(t *testing.T) {
		s, err := New(TestDSN)
		if err != nil {
			t.Errorf("Can't connect to DB: %v", err)
			return
		}
		_, err = s.db.Exec(context.Background(), "DELETE FROM sources;")
		if err != nil {
			t.Errorf("Can't cleanup sources: %v", err)
			return
		}
		var sql string = `INSERT INTO sources
				(link, fetchUpdates, defaultInterval)
			VALUES
				('https://example.com/feed', 'on', 5),
				('https://example.net/feed', 'off', 15),
				('https://example.org/feed', 'on', 10);`
		_, err = s.db.Exec(context.Background(), sql)
		if err != nil {
			t.Errorf("Can't insert sources: %v", err)
			return
		}
		sources, err := s.SourceGetActive()
		if err != nil {
			t.Errorf("Can't get sources: %v", err)
			return
		}
		if len(sources) != 2 {
			t.Errorf("Expected 2 sources, but %v selected!", len(sources))
		}
	})
}

func TestStore_PostAdd(t *testing.T) {
	t.Run("Duplicate", func(t *testing.T) {
		s, err := New(TestDSN)
		if err != nil {
			t.Errorf("Can't connect to DB: %v", err)
			return
		}
		_, err = s.db.Exec(context.Background(), "DELETE FROM posts;")
		if err != nil {
			t.Errorf("Can't cleanup posts: %v", err)
			return
		}
		_, err = s.db.Exec(context.Background(), "DELETE FROM sources;")
		if err != nil {
			t.Errorf("Can't cleanup sources: %v", err)
			return
		}
		var sql string = `INSERT INTO sources
				(id, link, fetchUpdates, defaultInterval)
			VALUES
				(1, 'https://example.com/feed', 'on', 5);`
		_, err = s.db.Exec(context.Background(), sql)
		if err != nil {
			t.Errorf("Can't insert sources: %v", err)
			return
		}
		var post Post = Post{
			Source:  1,
			Title:   "Title 1",
			Content: "Content 1",
			PubTime: time.Now().Unix(),
			Link:    "https://example.com/item/1",
			Guid:    "xxxx-yyyy",
		}
		id, err := s.PostAdd(post)
		if err != nil {
			t.Errorf("Can't add posts: %v", err)
			return
		}
		if id < 1 {
			t.Error("Id must be 1 or grater!")
			return
		}
		id, err = s.PostAdd(post)
		if err != nil {
			t.Errorf("Can't add posts (duplicate): %v", err)
			return
		}
		if id != 0 {
			t.Errorf("Id: expected 0, but %v found!", id)
			return
		}
	})
}

func TestStore_PostSearch(t *testing.T) {
	t.Run("Main", func(t *testing.T) {
		s, err := New(TestDSN)
		if err != nil {
			t.Errorf("Can't connect to DB: %v", err)
			return
		}
		_, err = s.db.Exec(context.Background(), "DELETE FROM posts;")
		if err != nil {
			t.Errorf("Can't cleanup posts: %v", err)
			return
		}
		_, err = s.db.Exec(context.Background(), "DELETE FROM sources;")
		if err != nil {
			t.Errorf("Can't cleanup sources: %v", err)
			return
		}
		var sql string = `INSERT INTO sources
				(id, link, fetchUpdates, defaultInterval)
			VALUES
				(1, 'https://example.com/feed', 'on', 5);`
		_, err = s.db.Exec(context.Background(), sql)
		if err != nil {
			t.Errorf("Can't insert sources: %v", err)
			return
		}
		sql = `INSERT INTO posts
				(id, source, title, content, pubTime, link, guid)
			VALUES
				(1, 1, 'title 1', 'content 1', 1688924657, 'https://example.com/n1', 'aaaa'),
				(2, 1, 'title 2', 'content 2', 1688924667, 'https://example.com/n2', 'bbbb'),
				(3, 1, 'title 3', 'content 3', 1688924647, 'https://example.com/n3', 'cccc'),
				(4, 1, 'title 4', 'content 4', 1688924627, 'https://example.com/n4', 'dddd');`
		_, err = s.db.Exec(context.Background(), sql)
		if err != nil {
			t.Errorf("Can't insert posts: %v", err)
			return
		}
		posts, err := s.PostSearch(10, 0, "")
		if err != nil {
			t.Errorf("Can't gat posts: %v", err)
			return
		}
		if len(posts) != 4 {
			t.Errorf("Expected 4 posts, but found: %v", len(posts))
			return
		}
		if posts[0].Id != 2 {
			t.Errorf("First id: expected 2, found: %v", posts[0].Id)
		}
		if posts[1].Id != 1 {
			t.Errorf("First id: expected 1, found: %v", posts[0].Id)
		}
	})
}

func TestStore_PostGetById(t *testing.T) {
	t.Run("Main", func(t *testing.T) {
		s, err := New(TestDSN)
		if err != nil {
			t.Errorf("Can't connect to DB: %v", err)
			return
		}
		_, err = s.db.Exec(context.Background(), "DELETE FROM posts;")
		if err != nil {
			t.Errorf("Can't cleanup posts: %v", err)
			return
		}
		_, err = s.db.Exec(context.Background(), "DELETE FROM sources;")
		if err != nil {
			t.Errorf("Can't cleanup sources: %v", err)
			return
		}
		var sql string = `INSERT INTO sources
				(id, link, fetchUpdates, defaultInterval)
			VALUES
				(1, 'https://example.com/feed', 'on', 5);`
		_, err = s.db.Exec(context.Background(), sql)
		if err != nil {
			t.Errorf("Can't insert sources: %v", err)
			return
		}
		sql = `INSERT INTO posts
				(id, source, title, content, pubTime, link, guid)
			VALUES
				(1, 1, 'title 1', 'content 1', 1688924657, 'https://example.com/n1', 'aaaa'),
				(2, 1, 'title 2', 'content 2', 1688924667, 'https://example.com/n2', 'bbbb'),
				(3, 1, 'title 3', 'content 3', 1688924647, 'https://example.com/n3', 'cccc'),
				(4, 1, 'title 4', 'content 4', 1688924627, 'https://example.com/n4', 'dddd');`
		_, err = s.db.Exec(context.Background(), sql)
		if err != nil {
			t.Errorf("Can't insert posts: %v", err)
			return
		}

		for id := 4; id > 0; id-- {
			p, err := s.PostGetById(id)
			if err != nil {
				t.Errorf("Can't get post %v: %v", id, err)
				return
			}
			if p.Id != id {
				t.Errorf("Try to get post %v, but found %v", id, p.Id)
				return
			}
		}
	})

}
