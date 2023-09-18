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
	r.PUT("/:id/packaging", ph.updateProductPackaging)
}

//	@Summary		Get product packaging
//	@Description	Calculates number of packets based on product configuration
//	@Tags			packaging, products
//	@Produce		json
//	@Param			id	path		int64	true	"ID of the product"
//	@Param			qty	query		int64	true	"Order quantity for product"
//	@Success		200	{object}	[]shipping.PackConfig
//	@Failure		400	{object}	object{error=string}
//	@Failure		404
//	@Failure		500
//	@Router			/v1/products/{id}/packaging [get]
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

//	@Summary		Update product packaging configuration
//	@Description	Updates configuration for the specified product
//	@Tags			packaging, products
//	@Accept			json
//	@Produce		json
//	@Param			id			path	int64		true	"ID of the product"
//	@Param			pack_sizes	body	[]uint64	true	"The list of supported pack sizes"
//	@Success		204
//	@Failure		400	{object}	object{error=string}
//	@Failure		404
//	@Failure		500
//	@Router			/v1/products/{id}/packaging [put]
func (ph *productHandler) updateProductPackaging(c *gin.Context) {
	productID := c.Param("id")
	id, err := strconv.ParseUint(productID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}
	var req []uint64
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	err = ph.ps.UpdatePacksConfiguration(c.Request.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, shipping.InternalServerErr):
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error occurred"})
			return
		case errors.Is(err, shipping.ErrNotFound):
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "product id not found"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	c.Status(http.StatusNoContent)
}
