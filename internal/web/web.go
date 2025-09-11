package web

import (
	"github.com/gin-gonic/gin"
	"github.com/liancccc/goauto/internal/options"
	"github.com/liancccc/goauto/internal/paths"
	"github.com/projectdiscovery/gologger"
)

func StartWebServer(opts *options.Options) {
	r := gin.Default()
	setAuth(r, opts)
	setStatic(r)
	setRouter(r)
	r.Run(opts.WebServer.Addr)
}

func setStatic(r *gin.Engine) {
	r.Static("/workspace", paths.WorkspaceDir)
	r.GET("/static/*filepath", ServeEmbeddedStatic)
	r.GET("/", func(c *gin.Context) {
		ServeEmbeddedFile(c, "/index.html")
	})
	r.GET("/index.html", func(c *gin.Context) {
		ServeEmbeddedFile(c, "/index.html")
	})
	r.GET("/commands.html", func(c *gin.Context) {
		ServeEmbeddedFile(c, "/commands.html")
	})
	r.GET("/task-detail.html", func(c *gin.Context) {
		ServeEmbeddedFile(c, "/task-detail.html")
	})
}

func setAuth(r *gin.Engine, opts *options.Options) {
	if opts.WebServer.User == "" {
		opts.WebServer.User = randString(5)
	}
	if opts.WebServer.Pass == "" {
		opts.WebServer.Pass = randString(9)
	}
	r.Use(gin.BasicAuth(gin.Accounts{
		opts.WebServer.User: opts.WebServer.Pass,
	}))
	gologger.Info().Str("user", opts.WebServer.User).Str("pass", opts.WebServer.Pass).Msgf("auth")
}

func setRouter(r *gin.Engine) {
	r.GET("/exec", execCommandHandler)
	r.GET("/execHelp", execHelpHandler)
	r.GET("/task/list", getALLTaskDBHandler)
	r.DELETE("/task/", deleteTaskDir)
	r.GET("/task/detail", getTaskDetail)
	r.GET("/system/info", getSysInfoHandler)
	r.POST("/upload/targets", uploadTargets)
}
