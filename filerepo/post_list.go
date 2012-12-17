package filerepo

import (
    "github.com/darkhelmet/blargh/post"
)

type postList []*post.Post

func (pl postList) Len() int {
    return len(pl)
}

func (pl postList) Less(i, j int) bool {
    return pl[i].PublishedOn.After(pl[j].PublishedOn.Time)
}

func (pl postList) Swap(i, j int) {
    pl[i], pl[j] = pl[j], pl[i]
}
