package service

import (
	"Reborn-but-in-Go/Favorite/dao"

	"log"
	"strconv"
	"sync"
)

// 服务层后端代码

// FavoriteService 服务层
type FavoriteService struct {
	FavoriteDao *dao.FavoriteDao
}

//type FavoriteService struct {
//	VideoService
//	UserService
//}

// NewFavoriteService 创建一个新的 FavoriteService 实例
func NewFavoriteService(favoriteDao *dao.FavoriteDao) *FavoriteService {
	return &FavoriteService{
		FavoriteDao: favoriteDao,
	}
}

/*
	1.其他模块(video)需要使用的业务方法。
*/
//IsFavorite 根据当前视频id判断是否点赞了该视频。
//FavouriteCount 根据当前视频id获取当前视频点赞数量。
//GetTotalFavouriteVideoCount 根据userId获取这个用户点赞视频数量(点赞数)
//GetTotalFavouriteVideoCount(userId int64) (int64, error) 根据userId获取这个用户点赞视频数量

// IsFavorite 根据当前视频id判断是否点赞了该视频。
func (*FavoriteService) IsFavourite(videoId int64, userId int64) (bool, error) {
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

// FavouriteCount 根据当前视频id获取当前视频点赞数量。
func (*FavoriteService) FavouriteCount(videoId int64) (int64, error) {
	// 获取点赞用户列表
	userIdList, err := dao.GetFavoriteUserIdList(videoId)
	if err != nil {
		log.Printf(err.Error())
		return 0, err
	}

	// 将点赞用户列表的长度转换为 int64，作为点赞数量
	count := int64(len(userIdList))
	return count, nil
}

// GetTotalFavoritedCount 根据userId获取这个用户总共被点赞数量(获赞数)
// 该函数未完成
func (*FavoriteService) GetTotalFavoritedCount(userId int64) (int64, error) {
	//根据userId获取这个用户的发布视频列表信息
	videoIdList, err := dao.GetVideoIdList(userId)
	if err != nil {
		log.Printf(err.Error())
		return 0, err
	}
	var sum int64 //该用户的总被点赞数
	//提前开辟空间,存取每个视频的点赞数
	videoFavoriteCountList := new([]int64)
	//采用协程并发将对应videoId的点赞数添加到集合中去
	i := len(videoIdList)
	var wg sync.WaitGroup
	wg.Add(i)
	for j := 0; j < i; j++ {
		go favorite.addVideoFavoriteCount(videoIdList[j], videoFavoriteCountList, &wg)
	}
	wg.Wait()
	//遍历累加，求总被点赞数
	for _, count := range *videoFavoriteCountList {
		sum += count
	}
	return sum, nil
}

// 根据userId获取这个用户点赞视频数量
func (f *FavoriteService) GetTotalFavouriteVideoCount(userId int64) (int64, error) {
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

// 点赞状态改变
func (*FavoriteService) FavouriteAction(userId int64, videoId int64, actionType int8) error {
	// 维护数据库信息
	err := dao.UpdateFavorite(userId, videoId, actionType)
	if err != nil {
		return err
	}
	return nil
}

// GetFavouriteList 函数根据给定的 userId 和 curId 获取用户的点赞视频列表。
func (f *FavoriteService) GetFavouriteList(userId int64, curId int64) ([]Video, error) {
	// 根据 userId 查询用户点赞的视频 id 列表
	videoIdList, err := dao.GetFavoriteVideoIdList(userId)
	if err != nil {
		log.Printf("GetFavouriteList 获取用户点赞视频列表错误：%v", err)
		return nil, err
	}

	var favoriteVideoList []Video
	var wg sync.WaitGroup

	for _, videoIdStr := range videoIdList {
		if videoId, err := strconv.ParseInt(videoIdStr, 10, 64); err == nil && videoId != config.DefaultRedisValue {
			wg.Add(1)
			go func(id int64) {
				defer wg.Done()
				// 获取并添加视频到列表中
				video, err := f.GetVideo(id, curId)
				if err == nil {
					favoriteVideoList = append(favoriteVideoList, video)
				} else {
					log.Println("GetFavouriteList: ", err)
				}
			}(videoId)
		}
	}
	wg.Wait()
	return favoriteVideoList, nil
}
