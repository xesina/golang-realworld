package gorm

import (
	"github.com/jinzhu/gorm"
	"github.com/xesina/golang-realworld/articles"
)

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) articles.ArticleRepository {
	return &articleRepository{
		db: db,
	}
}

func (r *articleRepository) Find(id uint) (*articles.Article, error) {
	article := new(articles.Article)
	r.db.First(article).RecordNotFound()
	if err := r.db.First(article).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return article, nil
}

func (r *articleRepository) Create(article *articles.Article) error {
	if err := r.db.Create(article).Error; err != nil {
		return err
	}
	return nil
}

func (r *articleRepository) Update(article *articles.Article) error {
	if err := r.db.Update(article).Error; err != nil {
		return err
	}
	return nil
}

func (r *articleRepository) Delete(id uint) error {
	return r.db.Delete(articles.Article{}, "id = ?", id).Error
}
