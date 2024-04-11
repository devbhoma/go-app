package httpserver

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
)

type Base struct {
	Router *gin.Engine
}

type Reader interface {
	Get() *Base
	GetEngine() *gin.Engine
	Run(port string) error
	NotFound(c *gin.Context, env string)
	Recovery(c *gin.Context, err interface{})
	CORSMiddleware() gin.HandlerFunc
	AuthHandler(cb func(s bool, ctx *gin.Context)) gin.HandlerFunc
}

func NewServer(env string) Reader {
	if env == "local" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.Logger())
	eg := &Base{
		Router: router,
	}

	router.NoRoute(func(c *gin.Context) {
		eg.NotFound(c, env)
	})
	router.Use(gin.RecoveryWithWriter(nil, eg.Recovery))
	router.Use(eg.CORSMiddleware())

	if err := eg.Router.SetTrustedProxies([]string{}); err != nil {
		fmt.Printf("error in setting proxies", "err", err)
		return &Base{}
	}

	return &Base{
		Router: router,
	}
}

func (e *Base) GetEngine() *gin.Engine {
	return e.Router
}

func (e *Base) Get() *Base {
	return e
}

func (e *Base) Run(port string) error {
	if port == "" {
		panic("Port not be null")
	}
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: e.Router.Handler(),
	}
	fmt.Printf("Listening and serving HTTP on :%s\n", port)
	return srv.ListenAndServe()
}

func (e *Base) NotFound(c *gin.Context, env string) {
	c.JSON(404, gin.H{"code": "NOT_FOUND", "message": "not found"})
}

func (e *Base) Recovery(c *gin.Context, err interface{}) {
	fmt.Printf("recovery err: ", err)
	msg := "Internal server error"
	if er, ok := err.(*runtime.TypeAssertionError); ok {
		fmt.Printf("server recovery: ", er.Error())
	}
	c.JSON(http.StatusInternalServerError, &map[string]interface{}{
		"status":  false,
		"code":    "InternalServerError",
		"message": msg,
	})
	c.Abort()
}

func (e *Base) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (e *Base) AuthHandler(cb func(s bool, ctx *gin.Context)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cb(true, ctx)
	}
}
