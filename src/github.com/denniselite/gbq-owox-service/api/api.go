package api

import (
    "github.com/denniselite/toolkit/conn"
    "github.com/kataras/iris"
    "github.com/satori/go.uuid"
    "github.com/iris-contrib/middleware/logger"
	"github.com/iris-contrib/middleware/recovery"
)

type Context struct {
    Rmq         conn.RmqInt
    RmqString   string
    QueuePrefix string
}

func (c *Context) LogPayload (ctx *iris.Context) {
    uuidVal := uuid.NewV4().String()
    ctx.Set("uuid", uuidVal)
    w := ctx.Recorder()
    ctx.Next()
    ctx.Log("%s Request: IRIS PROBLEM Response: %s\n", uuidVal, w.Body())
}

func (c *Context) NewRouter() *iris.Framework {
    router := iris.New(iris.Configuration{ DisableBanner: true })

	router.UseFunc(recovery.Handler)
    router.UseFunc(c.LogPayload)
    router.Use(logger.New(logger.Config{ Status: true, IP: true, Method: true, Path: true }))

    router.Post("/data", c.Data)

    router.Get("/ping", func (ctx *iris.Context) {
        ctx.Text(iris.StatusOK, "pong")
    })

    return router
}