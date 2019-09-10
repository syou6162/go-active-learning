package repository

import (
	"github.com/lib/pq"
	"github.com/syou6162/go-active-learning/lib/model"
)

var hatenaBookmarkNotFoundError = model.NotFoundError("hatenaBookmark")

func (r *repository) UpdateHatenaBookmark(e *model.Example) error {
	if e.HatenaBookmark == nil || e.HatenaBookmark.Count == 0 {
		return nil
	}

	tmp, err := r.FindExampleByUlr(e.Url)
	if err != nil {
		return err
	}
	id := tmp.Id

	e.HatenaBookmark.ExampleId = id
	if _, err = r.db.NamedExec(`
INSERT INTO hatena_bookmark
( example_id,  title,  screenshot,  entry_url,  count,  url,  eid)
VALUES
(:example_id, :title, :screenshot, :entry_url, :count, :url, :eid)
ON CONFLICT (example_id)
DO UPDATE SET
title = :title, count = :count
;`, e.HatenaBookmark); err != nil {
		return err
	}

	hb := model.HatenaBookmark{}
	if err = r.db.Get(&hb, `SELECT id FROM hatena_bookmark WHERE example_id = $1;`, id); err != nil {
		return err
	}

	for _, b := range e.HatenaBookmark.Bookmarks {
		b.HatenaBookmarkId = hb.Id
		if _, err = r.db.NamedExec(`
INSERT INTO bookmark
(hatena_bookmark_id, "user", comment, timestamp, tags)
VALUES
(:hatena_bookmark_id, :user, :comment, :timestamp, :tags)
ON CONFLICT (hatena_bookmark_id, "user") DO NOTHING
;`, b); err != nil {
			return err
		}
	}
	return nil
}

func (r *repository) SearchHatenaBookmarks(examples model.Examples, limitForEachExample int) ([]*model.HatenaBookmark, error) {
	hatenaBookmarks := make([]*model.HatenaBookmark, 0)
	exampleIds := make([]int, 0)
	for _, e := range examples {
		exampleIds = append(exampleIds, e.Id)
	}

	query := `SELECT * FROM hatena_bookmark WHERE example_id = ANY($1);`
	err := r.db.Select(&hatenaBookmarks, query, pq.Array(exampleIds))
	if err != nil {
		return hatenaBookmarks, err
	}

	hatenaBookmarkIds := make([]int, 0)
	for _, hb := range hatenaBookmarks {
		hatenaBookmarkIds = append(hatenaBookmarkIds, hb.Id)
		hb.Bookmarks = make([]*model.Bookmark, 0)
	}
	if limitForEachExample == 0 {
		return hatenaBookmarks, nil
	}

	bookmarks := make([]*model.Bookmark, 0)
	query = `SELECT * FROM bookmark WHERE hatena_bookmark_id = ANY($1) ORDER BY timestamp LIMIT $2;`
	err = r.db.Select(&bookmarks, query, pq.Array(hatenaBookmarkIds), limitForEachExample)
	if err != nil {
		return hatenaBookmarks, err
	}

	bookmarksByHatenaBookmarkId := make(map[int][]*model.Bookmark)
	for _, b := range bookmarks {
		bookmarksByHatenaBookmarkId[b.HatenaBookmarkId] = append(bookmarksByHatenaBookmarkId[b.HatenaBookmarkId], b)
	}

	result := make([]*model.HatenaBookmark, 0)
	for _, hb := range hatenaBookmarks {
		bookmarks := bookmarksByHatenaBookmarkId[hb.Id]
		hb.Bookmarks = bookmarks
		result = append(result, hb)
	}
	return result, nil
}

func (r *repository) FindHatenaBookmark(e *model.Example, limit int) (*model.HatenaBookmark, error) {
	hatenaBookmark := &model.HatenaBookmark{}

	query := `SELECT * FROM hatena_bookmark WHERE example_id = $1;`
	err := r.db.Get(hatenaBookmark, query, e.Id)
	if err != nil {
		return hatenaBookmark, err
	}

	bookmarks := make([]*model.Bookmark, 0)
	if limit == 0 {
		hatenaBookmark.Bookmarks = bookmarks
		return hatenaBookmark, nil
	}

	hatenaBookmarkId := hatenaBookmark.Id
	query = `SELECT * FROM bookmark WHERE hatena_bookmark_id = $1 ORDER BY timestamp LIMIT $2;`
	err = r.db.Select(&bookmarks, query, hatenaBookmarkId, limit)
	if err != nil {
		return hatenaBookmark, err
	}

	hatenaBookmark.Bookmarks = bookmarks
	return hatenaBookmark, nil
}
