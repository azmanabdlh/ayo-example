package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	auth "github.com/azmanabdlh/ayo-example/internal/authentication/handler"
	authentication "github.com/azmanabdlh/ayo-example/internal/authentication/service"
	"github.com/azmanabdlh/ayo-example/internal/config"
	league "github.com/azmanabdlh/ayo-example/internal/league-management/handler"
	leagueSvc "github.com/azmanabdlh/ayo-example/internal/league-management/service"
	"github.com/azmanabdlh/ayo-example/internal/team-management/handler"
	teamManagement "github.com/azmanabdlh/ayo-example/internal/team-management/service"
	"github.com/azmanabdlh/ayo-example/pkg/logger"
	"github.com/azmanabdlh/ayo-example/pkg/middleware"
	"github.com/azmanabdlh/ayo-example/pkg/provider"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func startApp(cfg *config.Config) error {

	logger.Init()

	r := gin.New()

	r.Use(
		middleware.RequestContext(),
	)

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	provider := provider.NewJsonWebTokenProvider("123")

	authSvc := authentication.New(provider, db)

	// handler
	l := league.New(
		leagueSvc.New(db),
	)

	a := auth.New(authSvc)

	h := handler.New(
		teamManagement.New(db),
	)

	// init seed
	runSeed(db)

	prefix := r.Group("/api")
	{

		prefix.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Pong",
			})
		})

		// LEAGUE ===

		// RESTFULL api/teams

		// POST api/players
		// PUT api/players/:id
		// DELETE api/players/:id

		// POST api/matches
		// GET api/matches/:id
		// POST api/matches/:id/highlight

		prefix.GET("/players/:id", h.FindPlayer)
		prefix.GET("/teams", h.FindAllTeam)
		prefix.GET("/matches/:id", l.FindMatchHighlight)

		prefix.GET("/teams/:id", h.FindTeam)

		// ADMIN
		// POST api/login
		prefix.POST("/login", a.Login)

		group := prefix.Use(
			middleware.RequiredAuthentication(
				provider,
				authSvc,
			),
		)

		// auth
		group.POST("/teams", h.CreateTeam)
		group.PUT("/teams/:id", h.ModifyTeam)
		group.DELETE("/teams/:id", h.RemoveTeam)

		group.POST("/players", h.AssignPlayerToTeam)

		group.DELETE("/players/:id", h.RemovePlayer)

		group.POST("/matches", l.CreateMatch)
		group.POST("/matches/:id/goals", l.AddRecordGoal)
		group.POST("/matches/:id/finish", l.Finish)

		group.POST("/matches/:id/lineup", l.AssignMatchPlayerLineup)
		group.POST("/matches/:id/substitutions", l.SubstitutePlayer)
	}

	server := &http.Server{
		Addr:    ":" + cfg.Server.Address,
		Handler: r,

		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("http server running at :" + cfg.Server.Address)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	// wait OS signal
	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	sig := <-quit

	log.Printf("signal received: %v", sig)

	// shutdown timeout context
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	// graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown failed: %v", err)
	}

	log.Println("server stopped gracefully")

	return nil
}

func runSeed(db *gorm.DB) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("invalid Getwd: %v", err)
	}

	pathToSQL := filepath.Join(dir, "files", "main.sql")

	query, err := os.ReadFile(pathToSQL)
	if err != nil {
		log.Fatalf("invalid readfile: %v", err)
	}

	err = db.Exec(string(query)).Error
	if err != nil {
		log.Fatalf("invalid exec SQL: %v", err)
	}

}
