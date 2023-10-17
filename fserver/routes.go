package fserver

import (
	"newsdata/model"
	"newsdata/storage"
	"strconv"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/rs/zerolog"
	"gopkg.in/reform.v1"
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
		r.logger.Info().Msg("Get News request")
		ListData, err := r.urepo.GetNews()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}

		return c.Status(fiber.StatusOK).JSON(ListData)
	}
}

func (r *RBase) EditNews() fiber.Handler {
	return func(c *fiber.Ctx) error {
		r.logger.Info().Msgf("Edit News: %s", string(c.Body()))

		ID := c.Params("Id")

		var inputNews model.Article

		err := json.Unmarshal(c.Body(), &inputNews)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("JSON unmarshal error")
		}

		// простая валидация:
		// нужно проверить, что ID в пути равен ID в запросе JSON
		id, err := strconv.Atoi(ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("strconv error")
		}

		if id != int(inputNews.ID) {
			return c.Status(fiber.StatusBadRequest).SendString("Validation error")
		}

		article, err := r.urepo.EditNews(inputNews)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error in EditNews")
		}

		return c.Status(fiber.StatusOK).JSON(article)
	}
}

func (s *FServer) MapHandlers(rs Routes) {
	s.FiberApp.Use(basicauth.New(basicauth.Config{
		Users: s.Auth,
		Realm: "Forbidden",
		Authorizer: func(user, pass string) bool {
			passAuth, _ := s.Auth[user]

			return passAuth == pass
		},

		Unauthorized: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusUnauthorized).SendString("Authorizer failed")
		},
	}))

	s.FiberApp.Get("/list", rs.GetNews())
	s.FiberApp.Post("/edit/:Id", rs.EditNews())
}
