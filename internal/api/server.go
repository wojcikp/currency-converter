package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wojcikp/currency-converter/internal/types"
)

type GinServer struct {
	router    *gin.Engine
	server    *http.Server
	converter types.Converter
}

func NewGinServer(serverPort string, converter types.Converter) *GinServer {
	r := gin.Default()
	return &GinServer{
		router:    r,
		server:    &http.Server{Addr: fmt.Sprintf(":%s", serverPort), Handler: r},
		converter: converter,
	}
}

func (s *GinServer) RegisterRoutes() {
	s.router.GET("/rates", s.GetRates)
	s.router.GET("/exchange", s.ExchangeCryptoCurrencies)
}

func (s *GinServer) Run() error {
	return s.server.ListenAndServe()
}

func (s *GinServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
