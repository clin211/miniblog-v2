// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package apiserver

import (
	"context"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	appHandler "github.com/clin211/miniblog-v2/internal/apiserver/handler/http/app"
	systemHandler "github.com/clin211/miniblog-v2/internal/apiserver/handler/http/system"
	"github.com/clin211/miniblog-v2/internal/pkg/core"
	"github.com/clin211/miniblog-v2/internal/pkg/errno"
	mw "github.com/clin211/miniblog-v2/internal/pkg/middleware/gin"
	"github.com/clin211/miniblog-v2/pkg/server"
)

// ginServer 定义一个使用 Gin 框架开发的 HTTP 服务器.
type ginServer struct {
	srv server.Server
}

// 确保 *ginServer 实现了 server.Server 接口.
var _ server.Server = (*ginServer)(nil)

// NewGinServer 初始化一个新的 Gin 服务器实例.
func (c *ServerConfig) NewGinServer() server.Server {
	// 创建 Gin 引擎
	engine := gin.New()

	// 注册全局中间件，用于恢复 panic、设置 HTTP 头、添加请求 ID 等
	engine.Use(gin.Recovery(), mw.RequestIDMiddleware(), mw.AccessLogger(), mw.NoCache, mw.Cors, mw.Secure)

	// 注册 REST API 路由
	c.InstallRESTAPI(engine)

	httpsrv := server.NewHTTPServer(c.cfg.HTTPOptions, c.cfg.TLSOptions, engine)

	return &ginServer{srv: httpsrv}
}

// 安装/注册 API 路由。路由的路径和 HTTP 方法，严格遵循 REST 规范.
func (c *ServerConfig) InstallRESTAPI(engine *gin.Engine) {
	// 注册业务无关的 API 接口
	InstallGenericAPI(engine)

	// 创建核心业务处理器
	sys := systemHandler.NewHandler(c.biz, c.val)
	app := appHandler.NewHandler(c.biz, c.val)

	// 注册健康检查接口
	engine.GET("/healthz", sys.Healthz)

	authMiddlewares := []gin.HandlerFunc{mw.AuthnMiddleware(c.retriever), mw.AuthzMiddleware(c.authz)}

	// 注册 v1 版本 API 路由分组
	sysv1 := engine.Group("/v1/system")
	{

		// auth
		authentication := sysv1.Group("/auth")
		{
			authentication.POST("/login", sys.Login)
			authentication.PUT("/refresh-token", mw.AuthnMiddleware(c.retriever), sys.RefreshToken)
		}

		// 用户相关路由
		user := sysv1.Group("/users")
		{
			// 创建用户。这里要注意：创建用户是不用进行认证和授权的
			user.POST("", sys.CreateUser)
			user.Use(authMiddlewares...)
			user.PUT(":userID/change-password", sys.ChangePassword) // 修改用户密码
			user.PUT(":userID", sys.UpdateUser)                     // 更新用户信息
			user.DELETE(":userID", sys.DeleteUser)                  // 删除用户
			user.GET(":userID", sys.GetUser)                        // 查询用户详情
			user.GET("", sys.ListUser)                              // 查询用户列表.
		}

		// 博客相关路由
		post := sysv1.Group("/posts", authMiddlewares...)
		{
			post.POST("", sys.CreatePost)       // 创建博客
			post.PUT(":postID", sys.UpdatePost) // 更新博客
			post.DELETE("", sys.DeletePost)     // 删除博客
			post.GET(":postID", sys.GetPost)    // 查询博客详情
			post.GET("", sys.ListPost)          // 查询博客列表
		}

		// 标签相关路由
		tag := sysv1.Group("/post-tags", authMiddlewares...)
		{
			tag.POST("", sys.CreateTag)      // 创建标签
			tag.PUT(":id", sys.UpdateTag)    // 更新标签
			tag.DELETE(":id", sys.DeleteTag) // 删除标签
			tag.GET(":id", sys.GetTag)       // 查询标签详情
			tag.GET("", sys.ListTag)         // 查询标签列表
		}

		// 分类相关路由
		category := sysv1.Group("/categories", authMiddlewares...)
		{
			category.POST("", sys.CreateCategory)      // 创建分类
			category.PUT(":id", sys.UpdateCategory)    // 更新分类
			category.DELETE(":id", sys.DeleteCategory) // 删除分类
			category.GET(":id", sys.GetCategory)       // 查询分类详情
			category.GET("", sys.ListCategory)         // 查询分类列表
		}

		// 设备相关路由
		device := sysv1.Group("/devices")
		{
			device.POST("", sys.CreateDevice)      // 创建设备
			device.PUT(":id", sys.UpdateDevice)    // 更新设备
			device.DELETE(":id", sys.DeleteDevice) // 删除设备
			device.GET(":id", sys.GetDevice)       // 获取单个设备
			device.GET("", sys.ListDevices)        // 查询设备列表
		}
	}

	appv1 := engine.Group("/v1/app")
	{
		post := appv1.Group("/posts")
		{
			post.GET("", app.ListPost)       // 查询所有文章
			post.GET(":postID", app.GetPost) // 查询单篇文章
		}

		category := appv1.Group("/categories")
		{
			category.GET("", app.ListCategories) // 查询所有分类
			category.GET(":id", app.GetCategory) // 查询单条分类
		}
	}
}

// InstallGenericAPI 安装业务无关的路由，例如 pprof、404 处理等.
func InstallGenericAPI(engine *gin.Engine) {
	// 注册 pprof 路由
	pprof.Register(engine)

	// 注册 404 路由处理
	engine.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errno.ErrPageNotFound, nil)
	})
}

// RunOrDie 启动 Gin 服务器，出错则程序崩溃退出.
func (s *ginServer) RunOrDie() {
	s.srv.RunOrDie()
}

// GracefulStop 优雅停止服务器.
func (s *ginServer) GracefulStop(ctx context.Context) {
	s.srv.GracefulStop(ctx)
}
