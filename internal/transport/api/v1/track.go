package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initTrackRoutes(api *gin.RouterGroup) {
	items := api.Group("/track")
	{
		items.GET("/list", h.trackList)
	}
}

func (h *Handler) trackList(c *gin.Context) {
	h.services.Track.List()
	c.String(http.StatusOK, "track list.......")
}
