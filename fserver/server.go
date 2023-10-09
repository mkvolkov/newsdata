package fserver

import (
	"context"
	"database/sql"
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

	address := cfg.Fserver.Host + ":" + cfg.Fserver.Port

	return &FServer{
		FiberApp: fiber.New(fCfg),
		Addr:     address,
		DB:       db,
		DReform:  dReform,
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
