package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rshby/go-event-ticketing/cacher"
	"github.com/rshby/go-event-ticketing/internal/controller"
	"github.com/rshby/go-event-ticketing/internal/middleware"
	"github.com/rshby/go-event-ticketing/internal/repository"
	"github.com/rshby/go-event-ticketing/internal/service"
	"gorm.io/gorm"
)

// SetupRouter setups router
func SetupRouter(app *gin.RouterGroup, db *gorm.DB, cache cacher.CacheManager) {
	eventRepository := repository.NewEventRepository(db, cache)
	eventService := service.NewEventService(eventRepository)
	eventController := controller.NewEventController(eventService)

	app.Use(middleware.TraceMiddleware())

	v1Group := app.Group("v1")
	{
		eventGroup := v1Group.Group("event")
		{
			eventGroup.POST("", eventController.CreateEvent)
			eventGroup.GET("", eventController.GetListEvents)
			eventGroup.GET(":id", eventController.GetEventByID)
		}
	}
}
