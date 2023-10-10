package fserver

import (
	"newsdata/model"
	"newsdata/storage"
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
	EditNews() fiber.Handler
}

type RBase struct {
	dreform *reform.DB
	logger  *zerolog.Logger
	urepo   storage.UserRepo
}

func (r *RBase) GetNews() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ListData, err := r.urepo.GetNews()
		if err != nil {
			return c.SendStatus(StatusIntErr)
		}

		return c.Status(StatusOK).JSON(ListData)
	}
}

func (r *RBase) EditNews() fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()
		mPath := strings.Trim(path, "/")
		pathParts := strings.Split(mPath, "/")

		if len(pathParts) < 2 {
			return c.Status(StatusBadReq).SendString("Expect /edit/<id> in task handler")
		}

		var inputNews model.News

		err := json.Unmarshal(c.Body(), &inputNews)
		if err != nil {
			return c.Status(StatusIntErr).SendString("JSON unmarshal error")
		}

		// простая валидация:
		// нужно проверить, что ID в пути равен ID в запросе JSON
		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			return c.SendStatus(StatusIntErr)
		}

		if id != int(inputNews.ID) {
			return c.SendStatus(StatusBadReq)
		}

		article, err := r.urepo.EditNews(inputNews)
		if err != nil {
			return c.SendStatus(StatusIntErr)
		}

		return c.Status(StatusOK).JSON(article)
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
	s.FiberApp.Post("/edit/:Id", rs.EditNews())
}
