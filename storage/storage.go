package storage

import (
	"newsdata/model"
	"sync"

	"gopkg.in/reform.v1"
)

type userRepo struct {
	Mtx     *sync.Mutex
	DReform *reform.DB
}

type UserRepo interface {
	GetNews() (model.ListNews, error)
	EditNews(model.Article) (model.Article, error)
}

func NewUserRepo(refDB *reform.DB) userRepo {
	return userRepo{Mtx: &sync.Mutex{}, DReform: refDB}
}

func (u userRepo) GetNews() (data model.ListNews, err error) {
	news, err := u.DReform.SelectAllFrom(model.ArticleTable, "")
	if err != nil {
		return data, err
	}

	categories, err := u.DReform.SelectAllFrom(model.CategoryView, "")
	if err != nil {
		return
	}

	var ListData model.ListNews

	var AllNews []*model.Article

	for i := 0; i < len(news); i++ {
		AllNews = append(AllNews, news[i].(*model.Article))
		for k := 0; k < len(categories); k++ {
			if AllNews[i].ID == categories[k].(*model.Category).ID {
				AllNews[i].Categories = append(AllNews[i].Categories, int(categories[k].(*model.Category).CatID))
			}
		}
	}

	ListData.Success = true
	ListData.AllNews = AllNews

	return ListData, nil
}

func (u userRepo) EditNews(inputArticle model.Article) (model.Article, error) {
	u.Mtx.Lock()
	defer u.Mtx.Unlock()

	if inputArticle.Title == "" || inputArticle.Content == "" {
		oldArticle, err := u.DReform.FindByPrimaryKeyFrom(model.ArticleTable, inputArticle.ID)
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

	nCat := len(inputArticle.Categories)

	if nCat > 0 {
		// удалить старые по NewsId
		u.DReform.DeleteFrom(model.CategoryView, "WHERE NewsId = ?", inputArticle.ID)

		// вставить новые
		var batch []reform.Struct

		for _, val := range inputArticle.Categories {
			batch = append(batch, &model.Category{ID: inputArticle.ID, CatID: int32(val)})
		}

		u.DReform.InsertMulti(batch...)
	}

	return inputArticle, nil
}
