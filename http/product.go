package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/silvan-talos/shipping"
	"github.com/silvan-talos/shipping/product"
)

type productHandler struct {
	ps product.Service
}

func (ph *productHandler) addRoutes(r *gin.RouterGroup) {
	r.GET("/:id/packaging", ph.getProductPackaging)
}

func (ph *productHandler) getProductPackaging(c *gin.Context) {
	productID := c.Param("id")
	id, err := strconv.ParseUint(productID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}
	quantity := c.Query("qty")
	if quantity == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "qty query param missing"})
		return
	}
	qty, err := strconv.ParseUint(quantity, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid qty"})
		return
	}
	resp, err := ph.ps.CalculatePacksConfiguration(c.Request.Context(), id, qty)
	if err != nil {
		switch {
		case errors.Is(err, shipping.InternalServerErr):
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error occurred"})
			return
		case errors.Is(err, shipping.ErrNotFound):
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "no configuration found for the specified product"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, resp)
}
