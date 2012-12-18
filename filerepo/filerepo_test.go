package filerepo_test

import (
    "github.com/darkhelmet/blargh/filerepo"
    . "launchpad.net/gocheck"
    "testing"
    "time"
)

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct{}

var (
    _       = Suite(&TestSuite{})
    good, _ = filerepo.New("test/repo")
)

func (ts *TestSuite) TestNew(c *C) {
    r, err := filerepo.New("test/repo")
    c.Assert(r, NotNil)
    c.Assert(err, IsNil)
}

func (ts *TestSuite) TestErrorsOnNoDirectory(c *C) {
    r, err := filerepo.New("test/nothing")
    c.Assert(r, IsNil)
    c.Assert(err, NotNil)
}

func (ts *TestSuite) TestOnlyLoadPublished(c *C) {
    c.Assert(good.Len(), Equals, 3)
}

func (ts *TestSuite) TestLatest(c *C) {
    latest, err := good.FindLatest(3)
    c.Assert(err, IsNil)
    c.Assert(len(latest), Equals, 3)
    c.Assert(latest[0].Title, Equals, "My Second Post")
    c.Assert(latest[1].Title, Equals, "My First Post")
}

func (ts *TestSuite) TestShouldSortOnPublishedOn(c *C) {
    r, err := filerepo.New("test/sort")
    c.Assert(err, IsNil)

    latest, err := r.FindLatest(3)

    c.Assert(len(latest), Equals, 3)
    c.Assert(latest[0].Title, Equals, "C")
    c.Assert(latest[1].Title, Equals, "A")
    c.Assert(latest[2].Title, Equals, "B")
}

func (ts *TestSuite) TestFindByTag(c *C) {
    posts, _ := good.FindByTag("meta")
    c.Assert(len(posts), Equals, 2)
    c.Assert(posts[0].Title, Equals, "My Second Post")
    c.Assert(posts[1].Title, Equals, "My First Post")
}

func (ts *TestSuite) BenchmarkFindByTag(c *C) {
    for i := 0; i < c.N; i++ {
        good.FindByTag("meta")
    }
}

func (ts *TestSuite) TestFindByCategory(c *C) {
    posts, err := good.FindByCategory("editorial")
    c.Assert(err, IsNil)
    c.Assert(len(posts), Equals, 1)
    c.Assert(posts[0].Title, Equals, "My First Post")

    posts, err = good.FindByCategory("writing")
    c.Assert(err, IsNil)
    c.Assert(len(posts), Equals, 1)
    c.Assert(posts[0].Title, Equals, "My Second Post")
}

func (ts *TestSuite) BenchmarkFindByCategory(c *C) {
    for i := 0; i < c.N; i++ {
        good.FindByCategory("editorial")
    }
}

func (ts *TestSuite) TestFindByMonth(c *C) {
    posts, err := good.FindByMonth(2012, time.December)
    c.Assert(err, IsNil)
    c.Assert(len(posts), Equals, 2)

    posts, err = good.FindByMonth(2012, time.February)
    c.Assert(err, IsNil)
    c.Assert(len(posts), Equals, 1)
    c.Assert(posts[0].Title, Equals, "An Old Post")
}

func (ts *TestSuite) BenchmarkFindByMonth(c *C) {
    for i := 0; i < c.N; i++ {
        good.FindByMonth(2012, time.December)
    }
}

func (ts *TestSuite) TestSearch(c *C) {
    posts, err := good.Search("first")
    c.Assert(err, IsNil)
    c.Assert(len(posts), Equals, 1)
    c.Assert(posts[0].Title, Equals, "My First Post")

    posts, err = good.Search("post")
    c.Assert(err, IsNil)
    c.Assert(len(posts), Equals, 3)
    c.Assert(posts[0].Title, Equals, "My Second Post")
    c.Assert(posts[1].Title, Equals, "My First Post")
    c.Assert(posts[2].Title, Equals, "An Old Post")

    posts, err = good.Search("old first")
    c.Assert(err, IsNil)
    c.Assert(len(posts), Equals, 2)
    c.Assert(posts[0].Title, Equals, "My First Post")
    c.Assert(posts[1].Title, Equals, "An Old Post")
}

func (ts *TestSuite) BenchmarkSearch(c *C) {
    for i := 0; i < c.N; i++ {
        good.Search("old first")
    }
}

// func (ts *TestSuite) TestLoadVerboseLogging(c *C) {
//     verboselogging, err := filerepo.New("test/verboselogging")
//     c.Assert(err, IsNil)
//     c.Assert(verboselogging.Len(), Equals, 179)
// }

// func (ts *TestSuite) BenchmarkLoadVerboseLogging(c *C) {
//     for i := 0; i < c.N; i++ {
//         filerepo.New("test/verboselogging")
//     }
// }
