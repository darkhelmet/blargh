package filerepo

import (
    "fmt"
    "github.com/darkhelmet/blargh/post"
    "github.com/darkhelmet/nltk"
    "github.com/darkhelmet/nltk/filter"
    "github.com/darkhelmet/nltk/tokenizer"
    "os"
    "path/filepath"
    "sort"
    "time"
)

type FileRepo struct {
    dir           string
    posts         []*post.Post
    categoryIndex map[string][]*post.Post
    tagIndex      map[string][]*post.Post
    searchIndex   map[string]*postSet
}

func New(dir string) (*FileRepo, error) {
    stat, err := os.Stat(dir)
    if err != nil {
        return nil, err
    }
    if !stat.Mode().IsDir() {
        return nil, fmt.Errorf("filerepo: %#v is not a directory", dir)
    }
    repo := &FileRepo{
        dir:           dir,
        posts:         make([]*post.Post, 0, 10),
        categoryIndex: make(map[string][]*post.Post),
        tagIndex:      make(map[string][]*post.Post),
        searchIndex:   make(map[string]*postSet),
    }
    err = repo.preload()
    if err != nil {
        return nil, fmt.Errorf("filerepo: failed preloading: %s", err)
    }
    sort.Sort(postList(repo.posts))
    repo.index()
    return repo, nil
}

func (fr *FileRepo) preload() error {
    return filepath.Walk(fr.dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }
        post, err := post.FromFile(path)
        if err != nil {
            return err
        }
        if post.Published {
            fr.posts = append(fr.posts, post)
        }
        return nil
    })
}

func tokens(bits ...string) nltk.TokenChan {
    tokens := tokenizer.Simple(bits...)
    tokens = filter.Superstrip(tokens)
    tokens = filter.SnowballStemmer(tokens)
    tokens = filter.DoubleMetaphone(tokens)
    return tokens
}

func tokensFor(p *post.Post) nltk.TokenChan {
    bits := []string{p.Title, p.Description, p.Clean()}
    bits = append(bits, p.Tags...)
    return tokens(bits...)
}

func (fr *FileRepo) index() {
    for _, post := range fr.posts {
        cp := fr.categoryIndex[post.Category]
        fr.categoryIndex[post.Category] = append(cp, post)

        for _, tag := range post.Tags {
            tp := fr.tagIndex[tag]
            fr.tagIndex[tag] = append(tp, post)
        }

        for tok := range tokensFor(post) {
            s := tok.String()
            ps, ok := fr.searchIndex[s]
            if !ok {
                ps = newPostSet()
                fr.searchIndex[s] = ps
            }
            ps.add(post)
        }
    }
}

func (fr *FileRepo) Len() int {
    return len(fr.posts)
}

func (fr *FileRepo) FindByTag(tag string) ([]*post.Post, error) {
    return fr.tagIndex[tag], nil
}

func (fr *FileRepo) FindByCategory(category string) ([]*post.Post, error) {
    return fr.categoryIndex[category], nil
}

func (fr *FileRepo) FindLatest(limit int) ([]*post.Post, error) {
    if fr.Len() < limit {
        limit = fr.Len()
    }
    return fr.posts[0:limit], nil
}

func (fr *FileRepo) FindByMonth(year int, month time.Month) ([]*post.Post, error) {
    posts := make([]*post.Post, 0)
    for _, post := range fr.posts {
        if post.InYear(year) && post.InMonth(month) {
            posts = append(posts, post)
        }
    }
    return posts, nil
}

func (fr *FileRepo) Search(query string) ([]*post.Post, error) {
    ps := newPostSet()
    for tok := range tokens(query) {
        if index, ok := fr.searchIndex[tok.String()]; ok {
            for _, post := range index.s {
                ps.add(post)
            }
        }
    }
    posts := ps.values()
    sort.Sort(postList(posts))
    return posts, nil
}

func (fr *FileRepo) FindByPermalink(year int, month time.Month, day int, slug string) (*post.Post, error) {
    for _, post := range fr.posts {
        if post.InYear(year) && post.InMonth(month) && post.HasSlug(slug) {
            return post, nil
        }
    }
    return nil, nil
}
