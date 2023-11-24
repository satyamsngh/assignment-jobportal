package main

import (
	"context"
	"fmt"
	"job-portal-api/config"
	"job-portal-api/internal/auth"
	"job-portal-api/internal/cache"
	"job-portal-api/internal/database"
	"job-portal-api/internal/handlers"
	"job-portal-api/internal/otputil"
	"job-portal-api/internal/redisutil"
	"job-portal-api/internal/repository"

	"net/http"
	"os"
	"os/signal"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"

	"time"
)

func main() {
	err := startApp()
	if err != nil {
		log.Panic().Err(err).Send()
	}
	log.Info().Msg("hello this is our app")
}

func startApp() error {
	Config := config.GetConfig()
	log.Info().Msg("main : Started : Initializing authentication support")
	fmt.Printf("%t=================================", Config.KeyConfig.PrivateKeyPath)
	fmt.Println("===============", Config.KeyConfig.PrivateKeyPath)

	privatePEM := Config.KeyConfig.PrivateKeyPath
	// if err != nil {
	// 	return fmt.Errorf("reading auth private key %w", err)
	// }

	if len(privatePEM) == 0 {
		return fmt.Errorf("error: Private key content is empty")

	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privatePEM))
	if err != nil {
		return fmt.Errorf("parsing auth private key %w", err)
	}

	publicPEM := []byte(Config.KeyConfig.PublicKeyPath)

	if len(privatePEM) == 0 {
		return fmt.Errorf("error: Public key content is empty")

	}
	// if err != nil {
	// 	return fmt.Errorf("reading auth public key %w", err)
	// }

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	if err != nil {
		return fmt.Errorf("parsing auth public key %w", err)
	}

	a, err := auth.NewAuth(privateKey, publicKey)
	if err != nil {
		return fmt.Errorf("constructing auth %w", err)
	}

	log.Info().Msg("main : Started : Initializing db support")
	dataConfig := Config.DataConfig
	db, err := database.Open(dataConfig)
	if err != nil {
		return fmt.Errorf("connecting to db %w", err)
	}
	pg, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w ", err)
	}
	redisConfig := Config.RedisConfig
	rd, err := redisutil.Redis(redisConfig)
	if err != nil {
		return fmt.Errorf("failed to connect redis: %w", err)
	}
	rdb, err := cache.NewCache(rd)
	if err != nil {
		log.Print("err")
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = pg.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("database is not connected: %w ", err)
	}

	repo, err := repository.NewRepository(db)
	if err != nil {
		log.Print("err")
		return err
	}

	err = repo.AutoMigrate()
	if err != nil {
		log.Print(err)
		return err
	}
	verify, err := otputil.NewOtp(Config.OtpGeneratorConfig.Port)
	if err != nil {
		return err
	}

	api := http.Server{
		Addr:         fmt.Sprintf(":%s", Config.AppConfig.Port),
		ReadTimeout:  8000 * time.Second,
		WriteTimeout: 800 * time.Second,
		IdleTimeout:  800 * time.Second,
		Handler:      handlers.API(a, repo, rdb, verify),
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Info().Str("port", api.Addr).Msg("main: API listening")
		serverErrors <- api.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error %w", err)
	case sig := <-shutdown:
		log.Info().Msgf("main: Start shutdown %s", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := api.Shutdown(ctx)
		if err != nil {
			err = api.Close()
			return fmt.Errorf("could not stop server gracefully %w", err)
		}
	}

	return nil
}
