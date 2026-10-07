package routes

import (
	"github.com/Wei-Shaw/sub2api/internal/handler"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

// RegisterPublicRoutes registers read-only routes used by standalone public pages.
func RegisterPublicRoutes(v1 *gin.RouterGroup, h *handler.Handlers, settingService *service.SettingService) {
	public := v1.Group("/public")
	{
		models := public.Group("/models")
		{
			models.GET("", h.ModelCatalog.List)
			models.GET("/vendors", h.ModelCatalog.Vendors)
		}

		monitors := public.Group("/channel-monitors")
		{
			monitors.GET("", h.ChannelMonitor.List)
			monitors.GET("/:id/status", h.ChannelMonitor.GetStatus)
		}

		// Public v2 passive views: same mode guard as the console, non-admin
		// redaction, no per-user group restriction, no users-ranking endpoint.
		monitorV2 := public.Group("/channel-monitor-v2")
		monitorV2.Use(channelMonitorModeV2Guard(settingService))
		{
			monitorV2.GET("/dimensions", h.ChannelMonitorV2.PublicDimensions)
			monitorV2.GET("/snapshot", h.ChannelMonitorV2.PublicSnapshot)
			monitorV2.GET("/models", h.ChannelMonitorV2.PublicModels)
			monitorV2.GET("/matrix", h.ChannelMonitorV2.PublicMatrix)
			monitorV2.GET("/errors", h.ChannelMonitorV2.PublicErrors)
		}
	}
}
