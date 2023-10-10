package storage

import (
	"newsdata/model"

	"gopkg.in/reform.v1"
)

type userRepo struct {
	DReform *reform.DB
}

type UserRepo interface {
	GetNews() (model.ListNews, error)
	EditNews(model.News) (model.Article, error)
}

func NewUserRepo(refDB *reform.DB) userRepo {
	return userRepo{DReform: refDB}
}

func (u userRepo) GetNews() (data model.ListNews, err error) {
	news, err := u.DReform.SelectAllFrom(model.ArticleTable, "")
	if err != nil {
		return data, err
	}

	var ListData model.ListNews

	var AllNews []model.News

	categories, err := u.DReform.SelectAllFrom(model.CategoryView, "")
	if err != nil {
		return
	}

	for i := 0; i < len(news); i++ {
		var article model.News
		id := news[i].(*model.Article).ID

		article.ID = id
		article.Title = news[i].(*model.Article).Title
		article.Content = news[i].(*model.Article).Content

		for k := 0; k < len(categories); k++ {
			if id == categories[k].(*model.Category).ID {
				article.Categories = append(article.Categories, int(categories[k].(*model.Category).CatID))
			}
		}

		AllNews = append(AllNews, article)
	}

	ListData.Success = true
	ListData.AllNews = AllNews

	return ListData, nil
}

func (u userRepo) EditNews(inputNews model.News) (model.Article, error) {
	var inputArticle model.Article
	inputArticle.ID = inputNews.ID
	inputArticle.Title = inputNews.Title
	inputArticle.Content = inputNews.Content

	if inputArticle.Title == "" || inputArticle.Content == "" {
		oldArticle, err := u.DReform.FindByPrimaryKeyFrom(model.ArticleTable, inputNews.ID)
		if err != nil {
			return inputArticle, err
		}

		if inputArticle.Title == "" {
			inputArticle.Title = oldArticle.(*model.Article).Title
		}

		if inputArticle.Content == "" {
			inputArticle.Content = oldArticle.(*model.Article).Content
		}
	}

	err := u.DReform.Save(&inputArticle)
	if err != nil {
		return inputArticle, err
	}

	if len(inputNews.Categories) > 0 {
		// удалить старые по NewsId
		u.DReform.DeleteFrom(model.CategoryView, "WHERE NewsId = ?", inputNews.ID)

		// вставить новые
		for _, val := range inputNews.Categories {
			u.DReform.Insert(&model.Category{ID: inputNews.ID, CatID: int32(val)})
		}
	}

	return inputArticle, nil
}
