package articles_test

import (
	"log"
	"os"
	"testing"
	articleMock "github.com/xesina/golang-realworld/mock/articles"
	"github.com/xesina/golang-realworld/articles"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

var (
	articleRepository *articleMock.ArticleRepository
)

func init() {
}
func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}
func setup() {
	log.Println("Running Setup ...")
	// make sure that mocks satisfies the interfaces
	var _ articles.ArticleRepository = (*articleMock.ArticleRepository)(nil)
	articleRepository = &articleMock.ArticleRepository{}
}
func tearDown() {
}

func TestFindArticle(t *testing.T) {
	expected := &articles.Article{
		ID:       1,
		Slug: "hello-world",
		Title:    "Hello World",
		Description: "test",
		Body: "test",
	}
	articleRepository.On("Find", expected.ID).Return(expected, nil)
	i := articles.NewArticleInteractor(articleRepository)
	actual, err := i.Find(expected.ID)
	require.NoError(t, err)
	assert.Equal(t, actual, expected)
	articleRepository.AssertExpectations(t)
}
