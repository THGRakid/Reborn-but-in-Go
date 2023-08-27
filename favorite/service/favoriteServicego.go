package service

import (
	"Reborn-but-in-Go/Favorite/dao"
	vidDao "Reborn-but-in-Go/video/dao"
	"log"
)

// 服务层后端代码

// FavoriteService 服务层
type FavoriteService struct {
	FavoriteDao *dao.FavoriteDao
}

// NewFavoriteService 创建一个新的 FavoriteService 实例
func NewFavoriteService(favoriteDao *dao.FavoriteDao) *FavoriteService {
	return &FavoriteService{
		FavoriteDao: favoriteDao,
	}
}

//token鉴权未考虑

/*
	1.其他模块(video)需要使用的业务方法。
*/
//IsFavorite 根据当前视频id判断是否点赞了该视频。
//FavouriteCount 根据当前视频id获取当前视频点赞数量。
//GetTotalFavouriteVideoCount 根据userId获取这个用户点赞视频数量(点赞数)
//GetTotalFavouriteVideoCount(userId int64) (int64, error) 根据userId获取这个用户点赞视频数量

// IsFavorite 根据当前视频id判断是否点赞了该视频。
func (fs *FavoriteService) IsFavourite(videoId int64, userId int64) (bool, error) {
	// 根据userId查询Favorites表，返回点赞的videoId列表
	videoIdList, err := fs.FavoriteDao.GetFavoriteUserIdList(userId)
	if err != nil {
		log.Printf(err.Error())
		return false, err
	}

	// 判断videoId是否在点赞列表中
	for _, FavoriteVideoId := range videoIdList {
		if FavoriteVideoId == videoId {
			return true, nil
		}
	}

	return false, nil
}

// FavouriteCount 根据当前视频id获取当前视频点赞数量。
func (fs *FavoriteService) FavouriteCount(videoId int64) (int64, error) {
	// 获取点赞用户列表
	userIdList, err := fs.FavoriteDao.GetFavoriteVideoIdList(videoId)
	if err != nil {
		log.Printf(err.Error())
		return 0, err
	}

	// 将点赞用户列表的长度转换为 int64，作为点赞数量
	count := int64(len(userIdList))
	return count, nil
}

// GetTotalFavoritedCount 根据userId获取这个用户总共被点赞数量(获赞数)
/*
	该函数存疑
*/
func (fs *FavoriteService) GetTotalFavoritedCount(userId int64) (int64, error) {
	videoIdList, err := vidDao.NewVideoDaoInstance().GetVideoListByUserId(userId)
	if err != nil {
		log.Printf(err.Error())
		return 0, err
	}
	var sum int64 // 该用户的总被点赞数

	for _, video := range videoIdList {
		count, err := fs.FavouriteCount(video.Id) // 修改这里
		if err != nil {
			log.Printf(err.Error())
			return 0, err
		}
		sum += count
	}

	return sum, nil
}

// 根据userId获取这个用户点赞视频数量
func (fs *FavoriteService) GetTotalFavouriteVideoCount(userId int64) (int64, error) {
	videoIdList, err := fs.FavoriteDao.GetFavoriteVideoIdList(userId)
	if err != nil {
		log.Printf(err.Error())
		return 0, err
	}
	return int64(len(videoIdList)), nil
}

/*
	2.request需要实现的功能
*/
//当前用户对视频的点赞操作 ,并把这个行为更新到favorite表中。
//当前操作行为，1点赞，2取消点赞。

// 点赞状态改变
func (fs *FavoriteService) FavouriteAction(userId int64, videoId int64, actionType int8) error {
	// 维护数据库信息
	err := fs.FavoriteDao.UpdateFavorite(userId, videoId, actionType)
	if err != nil {
		return err
	}
	return nil
}

// GetFavouriteList 函数根据给定的 userId 和 curId 获取用户的点赞视频列表。
/*
	该函数存疑
*/
func (fs *FavoriteService) GetFavouriteList(userId int64, curId int64) ([]int64, error) {
	// 根据 userId 查询用户点赞的视频 id 列表
	videoIdList, err := fs.FavoriteDao.GetFavoriteVideoIdList(userId)
	if err != nil {
		log.Printf("GetFavouriteList 获取用户点赞视频列表错误：%v", err)
		return nil, err
	}

	return videoIdList, nil
}
