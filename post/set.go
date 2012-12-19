package post

type PostSet struct {
    s map[string]*Post
}

func NewPostSet() *PostSet {
    return &PostSet{make(map[string]*Post)}
}

func (ps *PostSet) Add(p *Post) {
    ps.s[p.Id] = p
}

func (ps *PostSet) AddSet(other *PostSet) {
    for _, post := range other.s {
        ps.Add(post)
    }
}

func (ps *PostSet) Values() PostList {
    posts := make(PostList, 0, len(ps.s))
    for _, post := range ps.s {
        posts = append(posts, post)
    }
    return posts
}
