package service

import (
	"fmt"
	"log"

	"Reborn-but-in-Go/favorite/dao"
	userMod "Reborn-but-in-Go/user/model"
	vidDao "Reborn-but-in-Go/video/dao"
	vidMod "Reborn-but-in-Go/video/model"
)

type FavoriteService struct {
	FavoriteDao *dao.FavoriteDao
}

// NewFavoriteService 创建一个新的 FavoriteService 实例
func NewFavoriteService(favoriteDao *dao.FavoriteDao) *FavoriteService {
	return &FavoriteService{
		FavoriteDao: favoriteDao,
	}
}

/*
	1.其他模块(video)需要使用的业务方法。
*/

// IsFavorite 根据当前视频id判断是否点赞了该视频。
func IsFavorite(videoId int64, userId int64) (bool, error) {
	// 根据userId查询Favorites表，返回点赞的videoId列表
	videoIdList, err := dao.GetFavoriteVideoIdList(userId)
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

// FavoriteCount 根据当前视频id获取当前视频点赞数量。
func (fs *FavoriteService) FavoriteCount(videoId int64) (int64, error) {
	// 获取点赞用户列表
	userIdList, err := fs.FavoriteDao.GetFavoriteUserIdList(videoId)
	if err != nil {
		log.Printf(err.Error())
		return 0, err
	}

	// 将点赞用户列表的长度转换为 int64，作为点赞数量
	count := int64(len(userIdList))
	return count, nil
}

// GetTotalFavoriteCount 根据userId获取这个用户总共被点赞数量(获赞数)
func (fs *FavoriteService) GetTotalFavoriteCount(userId int64) (int64, error) {
	videoIdList, err := vidDao.NewVideoDaoInstance().GetVideoListByUserId(userId)
	if err != nil {
		log.Printf(err.Error())
		return 0, err
	}
	var sum int64 // 该用户的总被点赞数

	for _, video := range videoIdList {
		count, err := dao.GetFavoriteCount(video.Id) // 调用 GetFavoriteCount 函数获取 favorite_count
		if err != nil {
			log.Printf(err.Error())
			return 0, err
		}
		sum += count
	}

	return sum, nil

}

// GetTotalFavoriteVideoCount 根据userId获取这个用户点赞视频数量
func (fs *FavoriteService) GetTotalFavoriteVideoCount(userId int64) (int64, error) {
	videoIdList, err := dao.GetFavoriteVideoIdList(userId)
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

// FavoriteAction 点赞状态改变
func (fs *FavoriteService) FavoriteAction(userId int64, videoId int64, actionType int8) error {
	// 维护数据库信息
	err := fs.FavoriteDao.UpdateOrInsertFavorite(userId, videoId, actionType)
	if err != nil {
		return err
	}
	return nil
}

/*
	2.request需要实现的功能：获取喜欢列表的相关函数。
*/
// GetFavoriteList 函数根据给定的 userId 和 curId 获取用户的点赞视频列表。(返回值是videoid)
func (fs *FavoriteService) GetFavoriteList(userId int64) ([]int64, error) {
	// 根据 userId 查询用户点赞的视频 id 列表
	videoIdList, err := dao.GetFavoriteVideoIdList(userId)
	if err != nil {
		log.Printf("获取用户点赞视频列表错误：%v", err)
		return nil, err
	}

	return videoIdList, nil
}

// GetVideosByVideoIDs 函数根据给定的视频id切片获取用户的点赞视频实例列表。(返回值是video结构体的切片)
func (fs *FavoriteService) GetVideosByVideoIDs(videoIDs []int64) ([]vidMod.Video, error) {
	videos := make([]vidMod.Video, len(videoIDs))
	errors := make(chan error, len(videoIDs))

	// 使用并发获取视频信息
	for i, videoID := range videoIDs {
		go func(i int, videoID int64) {
			video, err := vidDao.NewVideoDaoInstance().GetVideoById(videoID)
			if err != nil {
				errors <- err
				return
			}
			videos[i] = *video
			errors <- nil
		}(i, videoID)
	}

	// 等待所有并发任务完成
	for range videoIDs {
		err := <-errors
		if err != nil {
			fmt.Printf("获取视频信息失败：%v\n", err)
		}
	}

	return videos, nil
}

func (fs *FavoriteService) GetUserByID(userID int64) (userMod.User, error) {
	user, err := dao.GetUserByID(userID)
	if err != nil {
		return userMod.User{}, err
	}
	return user, nil
}
