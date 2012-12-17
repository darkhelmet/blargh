package post_test

import (
    "github.com/darkhelmet/blargh/post"
    T "html/template"
    . "launchpad.net/gocheck"
    "testing"
    "time"
)

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct{}

func mustParseTime(s string) time.Time {
    t, err := time.Parse(time.RFC822, s)
    if err != nil {
        panic(err)
    }
    return t
}

var (
    _      = Suite(&TestSuite{})
    monday = mustParseTime("10 Dec 12 10:00 MST")
)

func (ts *TestSuite) TestLoadPost(c *C) {
    post, err := post.FromFile("test/my-first-post.md")
    c.Assert(err, IsNil)
    c.Assert(post.Id, Equals, "foobar")
    c.Assert(post.Category, Equals, "editorial")
    c.Assert(post.Description, Equals, "Just a first post to get things going.")
    c.Assert(post.Body, Equals, `This is my first post.

## How are ya'll?

* Foo
* Bar
* Baz`)
    c.Assert(post.Published, Equals, true)
    c.Assert(post.PublishedOn.Time, Equals, monday)
    c.Assert(post.Slugs, DeepEquals, []string{"my-first-post", "my-fist-post"})
    c.Assert(post.Slug(), Equals, "my-first-post")
    c.Assert(post.Tags, DeepEquals, []string{"meta", "foo", "bar"})
}

func (ts *TestSuite) TestSetsIdIfNoneSet(c *C) {
    post, err := post.FromFile("test/my-second-post.md")
    c.Assert(err, IsNil)
    c.Assert(post.Id, Not(Equals), "")
}

func (ts *TestSuite) TestMarkdown(c *C) {
    p := &post.Post{Body: `Hello, World!

# How is everybody today?`}

    html := p.HTML()
    c.Assert(html, Equals, T.HTML("<p>Hello, World!</p>\n\n<h1>How is everybody today?</h1>\n"))
}

func (ts *TestSuite) TestInlineHTML(c *C) {
    p := &post.Post{Body: `Hello, World!

# How is everybody today?

<img class="round bbottom bleft" src="/foo.jpg" />`}

    html := p.HTML()
    c.Assert(html, Equals, T.HTML("<p>Hello, World!</p>\n\n<h1>How is everybody today?</h1>\n\n<p><img class=\"round bbottom bleft\" src=\"/foo.jpg\" /></p>\n"))
}

func (ts *TestSuite) TestClean(c *C) {
    p := &post.Post{Body: `Hello, World!

# How is everybody today?

[Google](http://google.com/)

![an image](http://google.com/logo.jpg)

* Foo
* Bar
* Baz`}
    clean := p.Clean()
    c.Assert(clean, Equals, "Hello, World! How is everybody today? Google Foo Bar Baz")
}
