package blargh

import (
    "errors"
    "github.com/darkhelmet/blargh/filerepo"
    . "github.com/darkhelmet/blargh/post"
    "time"
)

func NewFileRepo(dir string) (Repo, error) {
    return filerepo.New(dir)
}

type Repo interface {
    Len() int
    FindByTag(string) ([]*Post, error)
    FindByCategory(string) ([]*Post, error)
    FindLatest(limit int) ([]*Post, error)
    FindByMonth(year int, month time.Month) ([]*Post, error)
    Search(string) ([]*Post, error)
    FindByPermalink(year int, month time.Month, day int, slug string) (*Post, error)
}
