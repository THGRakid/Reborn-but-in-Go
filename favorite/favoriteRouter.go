package favorite

import (
	"Reborn-but-in-Go/favorite/controller"
	"Reborn-but-in-Go/favorite/dao"
	"Reborn-but-in-Go/favorite/service"
	"github.com/gin-gonic/gin"
)

func InitFavoriteRouter(r *gin.Engine) {
	// 创建数据访问层（DAO）的单例实例
	favoriteDao := dao.NewFavoriteDaoInstance()

	// 创建服务层（Service）的实例，传递数据访问层实例
	favoriteService := service.NewFavoriteService(favoriteDao)

	// 创建表现层（Controller）的实例，传递服务层实例
	favoriteController := controller.NewFavoriteController(favoriteService)

	// 注册 GET 路由，处理点赞，使用表现层中的 FavoriteAction 函数
	r.POST("/douyin/favorite/action/", favoriteController.FavoriteAction)

	// 注册 POST 路由，处理点赞列表，使用表现层中的 GetFavoriteList 函数
	r.GET("/douyin/favorite/list/", favoriteController.GetFavoriteList)

}
