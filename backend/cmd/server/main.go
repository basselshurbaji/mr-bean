package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/basselshurbaji/mr_bean/backend/config"
	"github.com/basselshurbaji/mr_bean/backend/internal/auth"
	"github.com/basselshurbaji/mr_bean/backend/internal/gear"
	"github.com/basselshurbaji/mr_bean/backend/internal/health"
	"github.com/basselshurbaji/mr_bean/backend/internal/mailer"
	appmiddleware "github.com/basselshurbaji/mr_bean/backend/internal/middleware"
	"github.com/basselshurbaji/mr_bean/backend/internal/router"
	"github.com/basselshurbaji/mr_bean/backend/internal/user"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DB.DSN())
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("close db: %v", err)
		}
	}()

	if err := db.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	userRepo := user.NewPgUserRepo(db)
	gearRepo := gear.NewPgGearRepo(db)
	tokenSvc := auth.NewTokenService(cfg.Auth.JWTSecret, cfg.Auth.JWTExpiry, cfg.Auth.RefreshExpiry)
	mailerSvc := mailer.NewSMTPMailer(cfg.Mailer.Host, cfg.Mailer.Port, cfg.Mailer.Username, cfg.Mailer.Password, cfg.Mailer.From)
	authSvc := auth.NewAuthService(userRepo, tokenSvc, mailerSvc)
	userSvc := user.NewUserService(userRepo)
	gearSvc := gear.NewGearService(gearRepo)

	appmiddleware.Register(appmiddleware.TagAuthenticated, auth.Middleware(tokenSvc))

	r := router.NewRouter()

	for _, route := range []router.Route{
		router.Adapt(health.NewHandler()),
		router.Adapt(auth.NewLoginHandler(authSvc)),
		router.Adapt(auth.NewRefreshHandler(authSvc)),
		router.Adapt(auth.NewRegisterHandler(authSvc)),
		router.Adapt(user.NewMeHandler(userSvc)),
		router.Adapt(user.NewUpdateHandler(userSvc)),
		router.Adapt(user.NewChangePasswordHandler(userSvc)),
		router.Adapt(gear.NewListGearHandler(gearSvc)),
		router.Adapt(gear.NewCreateGearHandler(gearSvc)),
		router.Adapt(gear.NewGetGearHandler(gearSvc)),
		router.Adapt(gear.NewUpdateGearHandler(gearSvc)),
		router.Adapt(gear.NewDeleteGearHandler(gearSvc)),
		router.Adapt(gear.NewListStationsHandler(gearSvc)),
		router.Adapt(gear.NewCreateStationHandler(gearSvc)),
		router.Adapt(gear.NewUpdateStationHandler(gearSvc)),
		router.Adapt(gear.NewDeleteStationHandler(gearSvc)),
	} {
		router.Register(r, route)
	}

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("server listening on %s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
