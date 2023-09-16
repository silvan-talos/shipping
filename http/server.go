package http

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/silvan-talos/shipping"
	"github.com/silvan-talos/shipping/product"
)

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
