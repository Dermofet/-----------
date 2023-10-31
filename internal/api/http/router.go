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
	userHandlers handlers.UserHandlers
	authHandlers handlers.AuthHandlers
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
	userRepository := repository.NewUserRepository(pgSource)
	userInteractor := usecase.NewUserInteractor(userRepository)
	userPresenter := presenter.NewUserPresenter()
	tokenPresenter := presenter.NewTokenPresenter()
	r.handlers.authHandlers = handlers.NewAuthHandlers(userInteractor, tokenPresenter)

	authGroup := basePath.Group("/auth")
	authGroup.POST("/signup", r.handlers.authHandlers.SignUp)
	authGroup.POST("/signin", r.handlers.authHandlers.SignIn)

	forAllUserGroup := basePath.Group("/users")
	{
		forAllUserGroup.Use(middlewares.NewAuthMiddleware())
		forAllUserGroup.Use(middlewares.NewCheckRoleMiddleware([]string{entity.UserRole, entity.AdminRole}, userInteractor))

		r.handlers.userHandlers = handlers.NewUserHandlers(userInteractor, userPresenter)
		forAllUserGroup.GET("/me", r.handlers.userHandlers.GetMeHandler)
		forAllUserGroup.PUT("/me", r.handlers.userHandlers.UpdateMeHandler)
		forAllUserGroup.DELETE("/me", r.handlers.userHandlers.DeleteMeHandler)
		forAllUserGroup.GET("/id/:id", r.handlers.userHandlers.GetByIdHandler)
		forAllUserGroup.GET("/username/:username", r.handlers.userHandlers.GetByUsernameHandler)
		forAllUserGroup.PUT("/id/:id", r.handlers.userHandlers.UpdateHandler)
		forAllUserGroup.DELETE("/id/:id", r.handlers.userHandlers.DeleteHandler)
	}

	return nil
}
