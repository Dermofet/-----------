package handlers

import "github.com/gin-gonic/gin"

//go:generate mockgen -source=./interfaces.go -destination=./handlers_mock.go -package=handlers

type UserHandlers interface {
	GetMeHandler(c *gin.Context)
	UpdateMeHandler(c *gin.Context)
	DeleteMeHandler(c *gin.Context)
	GetByIdHandler(c *gin.Context)
	GetByUsernameHandler(c *gin.Context)
	UpdateHandler(c *gin.Context)
	DeleteHandler(c *gin.Context)
	LikeTrack(c *gin.Context)
	DislikeTrack(c *gin.Context)
	ShowLikedTracks(c *gin.Context)
}

type AuthHandlers interface {
	SignUp(c *gin.Context)
	SignIn(c *gin.Context)
}

type MusicHandlers interface {
	GetAll(c *gin.Context)
	GetAndSortByPopular(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}
