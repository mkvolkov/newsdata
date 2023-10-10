package fserver

import (
	"context"
	"database/sql"
	"fmt"
	"newsdata/cfg"
	"os"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/mysql"
)

type FServer struct {
	FiberApp *fiber.App
	Addr     string
	DB       *sql.DB
	DReform  *reform.DB
	Auth     map[string]string
	Logger   *zerolog.Logger
}

func NewServer(cfg *cfg.Cfg, db *sql.DB) *FServer {
	fCfg := fiber.Config{
		Prefork:     false,
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	}

	logger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel).With().Timestamp().Logger()

	dReform := reform.NewDB(db, mysql.Dialect, reform.NewPrintfLogger(logger.Printf))
	if dReform == nil {
		return nil
	}

	filePass, err := os.Open("login.txt")
	if err != nil {
		logger.WithLevel(zerolog.FatalLevel).Msgf("Fatal error: %v", err)
	}

	var user string
	var pass string

	var mpAuth = make(map[string]string)

	for {
		_, err := fmt.Fscanf(filePass, "%s %s\n", &user, &pass)
		if err != nil {
			break
		}

		mpAuth[user] = pass
	}

	address := cfg.Fserver.Host + ":" + cfg.Fserver.Port

	return &FServer{
		FiberApp: fiber.New(fCfg),
		Addr:     address,
		DB:       db,
		DReform:  dReform,
		Auth:     mpAuth,
		Logger:   &logger,
	}
}

func (s *FServer) Run(ctx context.Context) error {
	fHandlers := &RBase{
		dreform: s.DReform,
		logger:  s.Logger,
	}

	s.MapHandlers(fHandlers)

	if err := s.FiberApp.Listen(s.Addr); err != nil {
		s.Logger.Fatal().Msgf("Fiber Failed: Listen(), %v\n", err)
	}

	return nil
}
