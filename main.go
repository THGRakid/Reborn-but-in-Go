package main

import (
	"Reborn-but-in-Go/comment"
	"Reborn-but-in-Go/config"
	"Reborn-but-in-Go/follow"
	"Reborn-but-in-Go/message"
	"Reborn-but-in-Go/submission"
	"Reborn-but-in-Go/user"
	"Reborn-but-in-Go/video"
	"github.com/gin-gonic/gin"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	// 初始化 Redis 客户端
	wg.Add(1)
	go func() {
		defer wg.Done()
		config.InitRedis()
	}()

	// 等待 Redis 客户端初始化完成
	wg.Wait()

	r := gin.Default()
	r.Static("/static", "./public") // 创建一个默认的 Gin 路由引擎

	// 调用各模块接口
	video.InitVideoRouter(r)
	submission.InitSubmissionRouter(r)
	user.InitUserRouter(r)

	//	favorite.InitFavoriteRouter(r)
	comment.InitCommentRouter(r)
	message.InitMessageRouter(r)
	follow.InitFollowRouter(r)

	//输出语句（成功打开服务端）
	config.Begin()

	// 启动服务器并监听在 :8090 端口上
	if err := r.Run(":8090"); err != nil {
		panic("Failed to run server: " + err.Error())
	}

}
