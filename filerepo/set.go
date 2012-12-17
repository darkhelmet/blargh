package filerepo

import (
    "github.com/darkhelmet/blargh/post"
)

type postSet struct {
    s map[string]*post.Post
}

func newPostSet() *postSet {
    return &postSet{make(map[string]*post.Post)}
}

func (ps *postSet) add(p *post.Post) {
    ps.s[p.Id] = p
}

func (ps *postSet) values() []*post.Post {
    posts := make([]*post.Post, 0, len(ps.s))
    for _, post := range ps.s {
        posts = append(posts, post)
    }
    return posts
}
