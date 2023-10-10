package fserver

import (
	"newsdata/model"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/rs/zerolog"
	"gopkg.in/reform.v1"
)

const (
	StatusOK     int = 200
	StatusBadReq int = 400
	StatusIntErr int = 500
)

type Routes interface {
	GetNews() fiber.Handler
	PostNews() fiber.Handler
}

type RBase struct {
	dreform *reform.DB
	logger  *zerolog.Logger
}

func (r *RBase) GetNews() fiber.Handler {
	return func(c *fiber.Ctx) error {
		news, err := r.dreform.SelectAllFrom(model.ArticleTable, "")
		if err != nil {
			r.logger.Error().Msg("Error in Reform")
		}

		var ListData model.ListNews

		var AllNews []model.News

		categories, err := r.dreform.SelectAllFrom(model.CategoryView, "")
		if err != nil {
			r.logger.Error().Msg("Error in Reform")
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

		return c.Status(StatusOK).JSON(ListData)
	}
}

func (r *RBase) PostNews() fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()
		mPath := strings.Trim(path, "/")
		pathParts := strings.Split(mPath, "/")

		if len(pathParts) < 2 {
			return c.Status(StatusBadReq).SendString("Expect /edit/<id> in task handler")
		}

		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			return c.Status(StatusIntErr).SendString("Internal server error")
		}

		var inputNews model.News

		err = json.Unmarshal(c.Body(), &inputNews)
		if err != nil {
			return c.Status(StatusIntErr).SendString("JSON unmarshal error")
		}

		var inputArticle model.Article
		inputArticle.ID = inputNews.ID
		inputArticle.Title = inputNews.Title
		inputArticle.Content = inputNews.Content

		if inputArticle.Title == "" || inputArticle.Content == "" {
			oldArticle, err := r.dreform.FindByPrimaryKeyFrom(model.ArticleTable, inputNews.ID)
			if err != nil {
				return c.Status(StatusIntErr).SendString("Find old article error")
			}

			if inputArticle.Title == "" {
				inputArticle.Title = oldArticle.(*model.Article).Title
			}

			if inputArticle.Content == "" {
				inputArticle.Content = oldArticle.(*model.Article).Content
			}
		}

		err = r.dreform.Save(&inputArticle)
		if err != nil {
			return c.Status(StatusIntErr).SendString("reform error")
		}

		if len(inputNews.Categories) > 0 {
			// удалить старые по NewsId
			r.dreform.DeleteFrom(model.CategoryView, "WHERE NewsId = ?", inputNews.ID)

			// вставить новые
			for _, val := range inputNews.Categories {
				r.dreform.Insert(&model.Category{ID: inputNews.ID, CatID: int32(val)})
			}
		}

		return c.Status(StatusOK).SendString(strconv.Itoa(id))
	}
}

func (s *FServer) MapHandlers(rs Routes) {
	s.FiberApp.Use(basicauth.New(basicauth.Config{
		Users: s.Auth,
		Realm: "Forbidden",
		Authorizer: func(user, pass string) bool {
			passAuth, ok := s.Auth[user]
			if ok {
				if passAuth == pass {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		},

		Unauthorized: func(c *fiber.Ctx) error {
			return c.SendStatus(StatusBadReq)
		},
	}))

	s.FiberApp.Get("/list", rs.GetNews())
	s.FiberApp.Post("/edit/:Id", rs.PostNews())
}
