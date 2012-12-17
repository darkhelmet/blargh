package post

import (
    "bytes"
    "fmt"
    "github.com/darkhelmet/blargh/html"
    "github.com/james4k/fmatter"
    "github.com/nu7hatch/gouuid"
    md "github.com/russross/blackfriday"
    T "html/template"
    "regexp"
    "strings"
    "time"
)

var ws = regexp.MustCompile(`\s+`)

type Post struct {
    Id                                 string
    Title, Category, Description, Body string
    Published                          bool
    Slugs, Terms, Tags                 []string
    PublishedOn                        *Time
    Images                             map[string]map[string]string
}

func (p *Post) Slug() string {
    return p.Slugs[0]
}

func (p *Post) HTML() T.HTML {
    extensions := md.EXTENSION_SPACE_HEADERS | md.EXTENSION_TABLES | md.EXTENSION_NO_INTRA_EMPHASIS | md.EXTENSION_STRIKETHROUGH
    renderer := md.HtmlRenderer(0, "", "")
    return T.HTML(md.Markdown([]byte(p.Body), renderer, extensions))
}

func (p *Post) InYear(year int) bool {
    return p.PublishedOn.Year() == year
}

func (p *Post) InMonth(month time.Month) bool {
    return p.PublishedOn.Month() == month
}

func (p *Post) HasSlug(slug string) bool {
    for _, s := range p.Slugs {
        if s == slug {
            return true
        }
    }
    return false
}

func (p *Post) Clean() string {
    z := html.NewTokenizer(strings.NewReader(string(p.HTML())))
    var buffer bytes.Buffer
loop:
    for {
        switch tt := z.Next(); tt {
        case html.ErrorToken:
            break loop
        case html.TextToken:
            buffer.Write(z.Text())
        }
    }
    return string(bytes.TrimSpace(ws.ReplaceAll(buffer.Bytes(), []byte{' '})))
}

func FromFile(path string) (*Post, error) {
    var post Post
    content, err := fmatter.ReadFile(path, &post)
    if err != nil {
        return nil, fmt.Errorf("post: failed reading post: %s", err)
    }
    post.Body = string(content)
    if post.Id == "" {
        guid, err := uuid.NewV4()
        if err != nil {
            return nil, fmt.Errorf("post: failed generating UUID: %s", err)
        }
        post.Id = guid.String()
    }
    return &post, nil
}
