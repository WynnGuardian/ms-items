package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/wynnguardian/common/handlerfunc"
	middleware "github.com/wynnguardian/common/middlewares"
	"github.com/wynnguardian/common/response"
	"github.com/wynnguardian/ms-items/internal/domain/config"
	"github.com/wynnguardian/ms-items/internal/infra/http/handlers"
)

type RouterEntry struct {
	MustBeMod bool
	Handler   handlerfunc.HandlerFunc
	Path      string
	Method    string
}

var (
	entries = []RouterEntry{
		{Path: "/itemWeigh", MustBeMod: false, Method: "POST", Handler: handlers.WeightItem},
		{Path: "/itemAuthenticate", MustBeMod: true, Method: "POST", Handler: handlers.AuthItem},
		{Path: "/rankUpdate", MustBeMod: true, Method: "POST", Handler: handlers.UpdateRank},
		{Path: "/getRank", MustBeMod: false, Method: "POST", Handler: handlers.GetRank},
		{Path: "/createCriteria", MustBeMod: true, Method: "POST", Handler: handlers.CreateCriteria},
		{Path: "/deleteCriteria", MustBeMod: true, Method: "POST", Handler: handlers.DeleteCriteria},
		{Path: "/getCriteria", MustBeMod: false, Method: "POST", Handler: handlers.FindCriteria},
		{Path: "/getCriteriaByName", MustBeMod: false, Method: "POST", Handler: handlers.FindCriteria},
		{Path: "/updateCriteria", MustBeMod: true, Method: "POST", Handler: handlers.UpdateCriteria},
		{Path: "/findItem", MustBeMod: true, Method: "POST", Handler: handlers.FindItem},
	}
)

func post(engine *gin.Engine, path string, handler handlerfunc.HandlerFunc) {
	engine.POST(path, middleware.Parse(handler))
}

func Build() *gin.Engine {
	engine := gin.Default()
	engine.Use(func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Header("Access-Control-Allow-Origin", "http://guardian_proxy:8090")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	for _, entry := range entries {
		switch entry.Method {
		case "POST":
			if entry.MustBeMod {
				post(engine, entry.Path, middleware.RateLimit(middleware.NewRateLimiter(), config.MainConfig.Private.Tokens.RLWhitelist, Authorize(entry.Handler)))
			} else {
				post(engine, entry.Path, middleware.RateLimit(middleware.NewRateLimiter(), config.MainConfig.Private.Tokens.RLWhitelist, Authorize(entry.Handler)))
			}
		}
	}
	return engine
}

func Authorize(hf handlerfunc.HandlerFunc) handlerfunc.HandlerFunc {
	return func(ctx *gin.Context) response.WGResponse {
		wl := config.MainConfig.Private.Tokens.Whitelist
		fmt.Println(ctx.GetHeader("Authorization"))
		for _, w := range wl {
			if w == ctx.GetHeader("Authorization") {
				return hf(ctx)
			}
		}
		return response.ErrUnauthorized
	}
}
