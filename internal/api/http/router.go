package http

import (
	"fmt"
	"music-backend-test/docs"
	"music-backend-test/internal/api/http/handlers"
	"music-backend-test/internal/api/http/middlewares"
	"music-backend-test/internal/api/http/presenter"
	"music-backend-test/internal/db"
	"music-backend-test/internal/entity"
	"music-backend-test/internal/repository"
	"music-backend-test/internal/usecase"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"github.com/gin-contrib/cors"
)

type routerHandlers struct {
	userHandlers  handlers.UserHandlers
	authHandlers  handlers.AuthHandlers
	musicHandlers handlers.MusicHandlers
}

type router struct {
	router   *gin.Engine
	db       *sqlx.DB
	handlers routerHandlers
	logger   *zap.Logger
}

func NewRouter(db *sqlx.DB, logger *zap.Logger) *router {
	return &router{
		router: gin.New(),
		db:     db,
		logger: logger,
	}
}

func (r *router) Init() error {
	r.router.Use(
		gin.Logger(),
		gin.CustomRecovery(r.recovery),
	)
	err := r.registerRoutes()
	if err != nil {
		return fmt.Errorf("can't init router: %w", err)
	}

	return nil
}

func (r *router) recovery(c *gin.Context, recovered any) {
	defer func() {
		if e := recover(); e != nil {
			r.logger.Fatal("http server panic", zap.Error(fmt.Errorf("%s", recovered)))
		}
	}()
	c.AbortWithStatus(http.StatusInternalServerError)
}

func (r *router) registerRoutes() error {
	r.router.NoMethod(handlers.NotImplementedHandler)
	r.router.NoRoute(handlers.NotImplementedHandler)

	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})

	r.router.Use(corsMiddleware)

	basePath := r.router.Group(docs.SwaggerInfo.BasePath)

	basePath.GET("/swagger/swagger.json", func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache")
		c.File("docs/swagger.json")
	})
	basePath.GET("/swagger/swagger.yaml", func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache")
		c.File("docs/swagger.yaml")
	})
	basePath.GET("/docs/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("http://"+docs.SwaggerInfo.Host+docs.SwaggerInfo.BasePath+"/swagger/swagger.json"),
	))

	pgSource := db.NewSource(r.db)
	userSource := db.NewUserSourсe(pgSource)
	musicSource := db.NewMusicSource(pgSource)

	userRepository := repository.NewUserRepository(userSource)
	musicRepository := repository.NewMusicRepository(musicSource)

	userInteractor := usecase.NewUserInteractor(userRepository)
	musicInteractor := usecase.NewMusicInteractor(musicRepository)

	presenter := presenter.NewPresenter()

	r.handlers.authHandlers = handlers.NewAuthHandlers(userInteractor, presenter)

	authGroup := basePath.Group("/auth")
	authGroup.POST("/signup", r.handlers.authHandlers.SignUp)
	authGroup.POST("/signin", r.handlers.authHandlers.SignIn)

	userGroup := basePath.Group("/users")
	{
		userGroup.Use(middlewares.NewAuthMiddleware())

		r.handlers.userHandlers = handlers.NewUserHandlers(userInteractor, presenter)
		userGroup.GET("/me", r.handlers.userHandlers.GetMeHandler)
		userGroup.PUT("/me", r.handlers.userHandlers.UpdateMeHandler)
		userGroup.DELETE("/me", r.handlers.userHandlers.DeleteMeHandler)
		userGroup.GET(
			"/id/:id",
			middlewares.NewCheckRoleMiddleware([]string{entity.AdminRole}, userInteractor),
			r.handlers.userHandlers.GetByIdHandler,
		)
		userGroup.GET("/username/:username", r.handlers.userHandlers.GetByUsernameHandler)
		userGroup.PUT(
			"/:id",
			middlewares.NewCheckRoleMiddleware([]string{entity.AdminRole}, userInteractor),
			r.handlers.userHandlers.UpdateHandler,
		)
		userGroup.DELETE(
			"/:id",
			middlewares.NewCheckRoleMiddleware([]string{entity.AdminRole}, userInteractor),
			r.handlers.userHandlers.DeleteHandler,
		)
		userGroup.POST(
			"/add-track/:id",
			middlewares.NewCheckRoleMiddleware([]string{entity.AdminRole}, userInteractor),
			r.handlers.userHandlers.LikeTrack,
		)
		userGroup.DELETE(
			"/remove-track/:id",
			middlewares.NewCheckRoleMiddleware([]string{entity.AdminRole}, userInteractor),
			r.handlers.userHandlers.DislikeTrack,
		)
		userGroup.GET("/get-liked-tracks", r.handlers.userHandlers.ShowLikedTracks)
	}

	r.handlers.musicHandlers = handlers.NewMusicHandlers(musicInteractor, presenter)
	musicGroup := basePath.Group("/music")
	{
		musicGroup.Use(middlewares.NewAuthMiddleware())

		musicGroup.GET("/catalog", r.handlers.musicHandlers.GetAll)
		musicGroup.GET("/download/:id", r.handlers.musicHandlers.Get)
		musicGroup.GET("/release", r.handlers.musicHandlers.GetAllSortByTime)
		musicGroup.GET("/popular", r.handlers.musicHandlers.GetAndSortByPopular)
		musicGroup.POST(
			"/new",
			middlewares.NewCheckRoleMiddleware([]string{entity.AdminRole}, userInteractor),
			r.handlers.musicHandlers.Create,
		)
		musicGroup.PUT(
			"/:id",
			middlewares.NewCheckRoleMiddleware([]string{entity.AdminRole}, userInteractor),
			r.handlers.musicHandlers.Update,
		)
		musicGroup.DELETE(
			"/:id",
			middlewares.NewCheckRoleMiddleware([]string{entity.AdminRole}, userInteractor),
			r.handlers.musicHandlers.Delete,
		)
	}

	return nil
}
