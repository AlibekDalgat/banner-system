package delivery

import (
	"banner-system/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	services *service.Service
	cache    *redis.Client
}

func NewHandler(s *service.Service, rdb *redis.Client) *Handler {
	return &Handler{services: s, cache: rdb}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-in", h.signIn)
	}
	api := router.Group("/api", h.userIdentity)
	{
		api.GET("/user_banner", h.GetUserBanner)
		api.GET("/banner", h.GetAllBanners)
		api.POST("/banner", h.CreateBanner)
		api.PATCH("/banner/:id", h.UpdateBanner)
		api.DELETE("/banner/:id", h.DeleteBanner)
	}
	return router
}
