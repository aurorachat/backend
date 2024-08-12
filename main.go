package main

import (
	"github.com/aurorachat/auth/auth"
	"github.com/aurorachat/backend/chat"
	"github.com/aurorachat/backend/config"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	byeSignal := make(chan os.Signal)
	signal.Notify(byeSignal, syscall.SIGINT, syscall.SIGTERM)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	router := gin.New()

	router.Use(sloggin.New(logger))
	router.Use(gin.Recovery())

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  config.GetPostgresDSN(),
		PreferSimpleProtocol: true, // according to the docs gorm already does all the things
	}), &gorm.Config{})

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	authEngine, err := auth.NewEngine(auth.NewOptions(authHandler(), db, router, []byte(config.GetAuthSecretKey())))

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	_ = authEngine.SetUserActivated(1, true)

	chat.Initialize(router)

	go func() {
		err = router.Run(config.GetListenHost())

		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}()

	sig := <-byeSignal
	logger.Info("Shutting the application down", "signal", sig.String())
}

func authHandler() auth.Handler {
	return func(ctx *auth.ActionContext) {

	}
}
