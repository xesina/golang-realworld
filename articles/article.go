package articles

import (
	"github.com/xesina/golang-realworld/pkg/types"
	"time"
)

type Article struct {
	ID          uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   types.NullTime `sql:"index"`
	Slug        string         `gorm:"unique_index"`
	Title       string
	Description string `gorm:"size:2048"`
	Body        string `gorm:"size:2048"`
	AuthorID    uint
}

type ArticleRepository interface {
	Find(id uint) (*Article, error)
	Create(user *Article) error
	Update(user *Article) error
	Delete(id uint) error
}

type ArticleInteractor interface {
	Find(id uint) (*Article, error)
	Create(req ArticleRequest) (*Article, error)
	Update(user *Article) error
	Delete(id uint) error
}

type interactor struct {
	repo ArticleRepository
}

func NewArticleInteractor(r ArticleRepository) ArticleInteractor {
	return &interactor{
		repo: r,
	}
}

func (i *interactor) Find(id uint) (*Article, error) {
	article, err := i.repo.Find(id)
	if err != nil {
		return nil, err
	}
	return article, nil
}

type ArticleRequest struct {
	Title       string
	Description string
	Body        string
	Tags        []string
}

func (i *interactor) Create(req ArticleRequest) (*Article, error) {
	article := new(Article)

	article.Title = req.Title
	// TODO: Slug
	// TODO: Author ID
	article.Description = req.Description
	article.Body = req.Body

	err := i.repo.Create(article)

	// TODO: Return appropriate error (maybe wrapped)
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (i *interactor) Update(article *Article) error {
	err := i.repo.Update(article)
	if err != nil {
		return err
	}
	return nil
}

func (i *interactor) Delete(id uint) error {
	err := i.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
