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
	"github.com/basselshurbaji/mr_bean/backend/internal/bean"
	"github.com/basselshurbaji/mr_bean/backend/internal/extraction"
	"github.com/basselshurbaji/mr_bean/backend/internal/gear"
	"github.com/basselshurbaji/mr_bean/backend/internal/health"
	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
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

	userRepo       := user.NewPgUserRepo(db)
	gearRepo       := gear.NewPgGearRepo(db)
	beanRepo       := bean.NewPgBeanRepo(db)
	extractionRepo := extraction.NewPgExtractionRepo(db)
	appTokenRepo   := auth.NewPgAppTokenRepo(db)
	tokenSvc := auth.NewTokenService(cfg.Auth.JWTSecret, cfg.Auth.JWTExpiry, cfg.Auth.RefreshExpiry)
	authSvc      := auth.NewAuthService(userRepo, tokenSvc)
	appTokenSvc  := auth.NewAppTokenService(appTokenRepo, tokenSvc)
	userSvc      := user.NewUserService(userRepo)
	gearSvc      := gear.NewGearService(gearRepo)
	beanSvc      := bean.NewBeanService(beanRepo)
	extractionSvc := extraction.NewExtractionService(extractionRepo)

	middleware.Register(middleware.TagAuthenticated, auth.Middleware(tokenSvc))
	middleware.Register(middleware.TagAnyAuthenticated, middleware.Or(
		auth.Middleware(tokenSvc),
		auth.AppMiddleware(appTokenSvc),
	))

	r := router.NewRouter()

	router.Register(r, health.NewHandler())
	router.Register(r, auth.NewLoginHandler(authSvc))
	router.Register(r, auth.NewRefreshHandler(authSvc))
	router.Register(r, auth.NewRegisterHandler(authSvc))
	router.Register(r, user.NewMeHandler(userSvc))
	router.Register(r, user.NewUpdateHandler(userSvc))
	router.Register(r, user.NewChangePasswordHandler(userSvc))
	router.Register(r, gear.NewListGearHandler(gearSvc))
	router.Register(r, gear.NewCreateGearHandler(gearSvc))
	router.Register(r, gear.NewGetGearHandler(gearSvc))
	router.Register(r, gear.NewUpdateGearHandler(gearSvc))
	router.Register(r, gear.NewDeleteGearHandler(gearSvc))
	router.Register(r, gear.NewListStationsHandler(gearSvc))
	router.Register(r, gear.NewCreateStationHandler(gearSvc))
	router.Register(r, gear.NewUpdateStationHandler(gearSvc))
	router.Register(r, gear.NewDeleteStationHandler(gearSvc))
	router.Register(r, bean.NewListBeansHandler(beanSvc))
	router.Register(r, bean.NewCreateBeanHandler(beanSvc))
	router.Register(r, bean.NewUpdateBeanHandler(beanSvc))
	router.Register(r, bean.NewDeleteBeanHandler(beanSvc))
	router.Register(r, extraction.NewListExtractionsHandler(extractionSvc))
	router.Register(r, extraction.NewCreateExtractionHandler(extractionSvc))
	router.Register(r, extraction.NewGetExtractionHandler(extractionSvc))
	router.Register(r, extraction.NewUpdateExtractionHandler(extractionSvc))
	router.Register(r, extraction.NewDeleteExtractionHandler(extractionSvc))
	router.Register(r, auth.NewCreateAppTokenHandler(appTokenSvc))
	router.Register(r, auth.NewRevokeAppTokenHandler(appTokenSvc))

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("server listening on %s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
