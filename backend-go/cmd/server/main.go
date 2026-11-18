package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"whoami/backend-go/internal/auth"
	"whoami/backend-go/internal/config"
	"whoami/backend-go/internal/httpapi"
	"whoami/backend-go/internal/service"
	"whoami/backend-go/internal/store/gormdb"
	"whoami/backend-go/internal/store/localfs"
)

const storeTokenCleanupInterval = time.Hour

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config", "err", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(cfg.ProfileDataDir, 0o755); err != nil {
		slog.Error("create profile data dir", "err", err)
		os.Exit(1)
	}
	profileStore, err := localfs.NewProfileStore(cfg.ProfileDataDir)
	if err != nil {
		slog.Error("init profile store", "err", err)
		os.Exit(1)
	}

	db, err := gormdb.Open(cfg.DatabaseURL)
	if err != nil {
		slog.Error("init database", "err", err)
		os.Exit(1)
	}
	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("sql db", "err", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	userRepo := gormdb.NewUserRepo(db)
	keyRepo := gormdb.NewApiKeyRepo(db)
	tokenRepo := gormdb.NewStoreTokenRepo(db)

	tokenMgr := auth.NewTokenManager(cfg.JWTSecret, cfg.JWTExpireDays)
	stateStore := auth.NewStateStore(10*time.Minute, 10000)
	githubOAuth := auth.NewGitHubOAuth(cfg.GitHubClientID, cfg.GitHubClientSecret, cfg.GitHubRedirectURI)
	authSvc := service.NewAuthService(userRepo)
	keySvc := service.NewKeyService(keyRepo, userRepo, cfg.APIKeyPrefix)
	storeTokenSvc := service.NewStoreTokenService(tokenRepo)
	profileSvc := service.NewProfileService(profileStore)

	cleanupCtx, cleanupCancel := context.WithCancel(context.Background())
	defer cleanupCancel()
	go runStoreTokenCleanup(cleanupCtx, storeTokenSvc)

	router := httpapi.NewRouter(httpapi.Deps{
		Config:      cfg,
		Tokens:      tokenMgr,
		States:      stateStore,
		GitHub:      githubOAuth,
		Auth:        authSvc,
		Users:       userRepo,
		Keys:        keySvc,
		StoreTokens: storeTokenSvc,
		Profiles:    profileSvc,
	})
	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		slog.Info("whoami-go listening",
			"addr", cfg.HTTPAddr,
			"env", cfg.Env,
			"db", cfg.DatabaseURL,
			"profiles", cfg.ProfileDataDir,
		)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "err", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	cleanupCancel()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("shutdown", "err", err)
		os.Exit(1)
	}
	slog.Info("whoami-go stopped")
}

func runStoreTokenCleanup(ctx context.Context, svc *service.StoreTokenService) {
	ticker := time.NewTicker(storeTokenCleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			n, err := svc.CleanupExpired(ctx)
			if err != nil {
				slog.Error("store token cleanup failed", "err", err)
				continue
			}
			if n > 0 {
				slog.Info("cleaned expired store tokens", "count", n)
			}
		}
	}
}
