package filerepo

import (
    "fmt"
    . "github.com/darkhelmet/blargh/errors"
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
    posts         post.PostList
    categoryIndex map[string]post.PostList
    tagIndex      map[string]post.PostList
    searchIndex   map[string]*post.PostSet
    slugIndex     map[string]*post.Post
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
        posts:         make(post.PostList, 0, 10),
        categoryIndex: make(map[string]post.PostList),
        tagIndex:      make(map[string]post.PostList),
        searchIndex:   make(map[string]*post.PostSet),
        slugIndex:     make(map[string]*post.Post),
    }
    err = repo.preload()
    if err != nil {
        return nil, fmt.Errorf("filerepo: failed preloading: %s", err)
    }
    sort.Sort(repo.posts)
    err = repo.index()
    if err != nil {
        return nil, fmt.Errorf("filerepo: failed indexing", err)
    }
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

func (fr *FileRepo) index() error {
    for _, p := range fr.posts {
        cp := fr.categoryIndex[p.Category]
        fr.categoryIndex[p.Category] = append(cp, p)

        for _, tag := range p.Tags {
            tp := fr.tagIndex[tag]
            fr.tagIndex[tag] = append(tp, p)
        }

        for tok := range tokensFor(p) {
            s := tok.String()
            ps, ok := fr.searchIndex[s]
            if !ok {
                ps = post.NewPostSet()
                fr.searchIndex[s] = ps
            }
            ps.Add(p)
        }

        slug := p.Slug()
        _, found := fr.slugIndex[slug]
        if found {
            return fmt.Errorf("slug %#v already found", slug)
        }
        fr.slugIndex[slug] = p
    }
    return nil
}

func (fr *FileRepo) Len() int {
    return len(fr.posts)
}

func (fr *FileRepo) FindByTag(tag string) (post.PostList, error) {
    posts := fr.tagIndex[tag].PublishedBefore(time.Now())
    if len(posts) == 0 {
        return posts, NotFound(fmt.Sprintf("No posts found for tag %#v", tag))
    }
    return posts, nil
}

func (fr *FileRepo) FindByCategory(category string) (post.PostList, error) {
    posts := fr.categoryIndex[category].PublishedBefore(time.Now())
    if len(posts) == 0 {
        return posts, NotFound(fmt.Sprintf("No posts found for category %#v", category))
    }
    return posts, nil
}

func (fr *FileRepo) FindLatest(limit int) (post.PostList, error) {
    if fr.Len() < limit {
        limit = fr.Len()
    }
    return fr.posts.PublishedBefore(time.Now())[0:limit], nil
}

func (fr *FileRepo) FindByMonth(year int, month time.Month) (post.PostList, error) {
    posts := make(post.PostList, 0)
    for _, post := range fr.posts.PublishedBefore(time.Now()) {
        if post.InYear(year) && post.InMonth(month) {
            posts = append(posts, post)
        }
    }
    return posts.PublishedBefore(time.Now()), nil
}

func (fr *FileRepo) Search(query string) (post.PostList, error) {
    ps := post.NewPostSet()
    for tok := range tokens(query) {
        if index, ok := fr.searchIndex[tok.String()]; ok {
            ps.AddSet(index)
        }
    }
    posts := ps.Values()
    sort.Sort(posts)
    posts = posts.PublishedBefore(time.Now())
    if len(posts) == 0 {
        return posts, NotFound(fmt.Sprintf("No posts found for query %#v", query))
    }
    return posts, nil
}

func (fr *FileRepo) FindBySlug(slug string) (*post.Post, error) {
    post, ok := fr.slugIndex[slug]
    if ok {
        return post, nil
    }
    return nil, NotFound(fmt.Sprintf("Post not found with slug %#v", slug))
}
