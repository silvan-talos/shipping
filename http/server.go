package http

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/silvan-talos/shipping"
	_ "github.com/silvan-talos/shipping/docs"
	"github.com/silvan-talos/shipping/product"
)

//	@title			Shipping API docs
//	@description	Shipping is a small API that calculates packaging configuration for a certain amount of ordered product quantity.
//	@version		1.0.0
//	@host			cbhbw91cn7.execute-api.eu-west-1.amazonaws.com
//	@schemes		https
type Server struct {
	router *gin.Engine
}

func NewServer(args ServerArgs) *Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	setLogsFormat(r)
	r.Use(gin.Recovery())

	err := shipping.Validate.Struct(args)
	if err != nil {
		log.Fatal("http server failed to start, args missing, err:", err)
	}

	v1 := r.Group("/v1")
	{
		productRoutes := v1.Group("/products")
		{
			h := productHandler{
				ps: args.ProductService,
			}
			h.addRoutes(productRoutes)
		}
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive"})
	})
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.InstanceName("ShippingAPI")))

	return &Server{
		router: r,
	}
}

type ServerArgs struct {
	ProductService product.Service `validate:"required"`
}

func (s *Server) Serve(lis net.Listener) error {
	log.Println("Starting http server, address:", lis.Addr().String())
	return s.router.RunListener(lis)
}

func setLogsFormat(r *gin.Engine) {
	r.Use(gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		return fmt.Sprintf("[GIN]: %s - client=%s\tlatency=%s\tstatus=%d\tpath=\"%s %s\"\n",
			params.TimeStamp,
			params.ClientIP,
			params.Latency,
			params.StatusCode,
			params.Method,
			params.Path,
		)
	}))
}
