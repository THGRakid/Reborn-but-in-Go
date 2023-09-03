package dao

import (
	"Reborn-but-in-Go/config"
	"Reborn-but-in-Go/favorite/model"
	vidMod "Reborn-but-in-Go/video/model"
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

// FavoriteDao 把dao层看成整体，把dao的curd封装在一个结构体中。
type FavoriteDao struct {
}

var (
	favoriteDao  *FavoriteDao //操作该dao层crud的结构体变量。
	favoriteOnce sync.Once    //单例限定，去限定申请一个followDao结构体变量。sync.Once 是一个同步原语（synchronization primitive）用于确保在并发环境下只执行一次特定的操作。它通常用于延迟初始化或只需要在程序的生命周期内执行一次的操作
)

// NewFavoriteDaoInstance 用于获取 FavoriteDao 单例实例的函数
// 传递一个匿名函数（闭包），其作用是创建一个新的 FavoriteDao 实例并将其赋值给 FavoriteDao 变量。
// 这个函数只会在第一次调用 Do 方法时执行
func NewFavoriteDaoInstance() *FavoriteDao {
	favoriteOnce.Do(
		func() {
			favoriteDao = &FavoriteDao{}
		})
	return favoriteDao
}

//以下为favorite对数据库的操作方法

// GetFavoriteUserIdList 根据videoId获取点赞userId,获取视频的所有点赞用户
func (*FavoriteDao) GetFavoriteUserIdList(videoId int64) ([]int64, error) {
	var FavoriteUserIdList []int64 // 存储所有该视频点赞用户id
	// 查询Favorites表对应视频id点赞用户，返回查询结果
	err := config.DB.
		Model(model.Favorite{}).
		Where(map[string]interface{}{"video_id": videoId, "status": 1}).
		Pluck("user_id", &FavoriteUserIdList).Error //匹配查询条件的记录中的 user_id 字段，并将提取的值存储在 &FavoriteUserIdList 变量中。
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("get FavoriteUserIdList failed")
	}
	return FavoriteUserIdList, nil
}

// 根据videoId查询video表中的favorite_count字段
func GetFavoriteCount(videoId int64) (int64, error) {
	var video vidMod.Video
	err := config.DB.Model(vidMod.Video{}).Where(map[string]interface{}{"id": videoId}).First(&video).Error
	if err != nil {
		log.Println(err.Error())
		return 0, errors.New("get video favorite_count failed")
	}
	return video.FavoriteCount, nil
}

// InsertFavorite 插入点赞数据
func (*FavoriteDao) InsertFavorite(FavoriteData model.Favorite) error {
	// 创建点赞数据，默认为点赞，status为1
	err := config.DB.Model(model.Favorite{}).Create(&FavoriteData).Error
	if err != nil {
		log.Println(err.Error())
		return errors.New("insert data fail")
	}
	log.Printf("添加点赞记录成功")
	//同时更改video表中的favorite_count字段
	verr := config.DB.Model(vidMod.Video{}).Where(map[string]interface{}{"id": FavoriteData.VideoId}).
		Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
	if verr != nil {
		log.Println(verr.Error())
		return errors.New("add video favorite_count fail")
	}
	log.Printf("点赞数量加一")
	return nil
}

// UpdateFavorite 取消赞删除记录
func (*FavoriteDao) UpdateFavorite(userId int64, videoId int64, actionType int8) error {
	if actionType == 2 {
		// 删除记录
		err := config.DB.Delete(&model.Favorite{}, "user_id = ? AND video_id = ?", userId, videoId).Error
		if err != nil {
			log.Printf("删除记录失败：%v", err)
			return err
		}
		log.Printf("删除点赞记录成功")
		verr := config.DB.Model(vidMod.Video{}).Where(map[string]interface{}{"id": videoId}).
			Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error
		if verr != nil {
			log.Println(verr.Error())
			return errors.New("subtract video favorite_count fail")
		}
		log.Printf("点赞数量减一")
	} else {
		// 其他操作或错误处理
		return errors.New("不支持的操作类型")
	}

	return nil
}
func (*FavoriteDao) UpdateOrInsertFavorite(userId int64, videoId int64, actionType int8) error {
	// 查询是否已存在该用户对该视频的点赞记录
	existingFavorite, err := favoriteDao.GetFavoriteInfo(userId, videoId)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	if actionType == 1 && existingFavorite.UserId != userId {
		// 如果记录不存在，插入点赞数据
		newFavorite := model.Favorite{
			UserId:   userId,
			VideoId:  videoId,
			Status:   1, // 默认点赞状态为1
			CreateAt: time.Now(),
		}
		err := favoriteDao.InsertFavorite(newFavorite)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	} else {
		// 如果已点赞记录存在，直接删除
		err := favoriteDao.UpdateFavorite(userId, videoId, actionType)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	}

	return nil
}

// GetFavoriteInfo 根据userId,videoId查询点赞信息
func (*FavoriteDao) GetFavoriteInfo(userId int64, videoId int64) (model.Favorite, error) {
	// 根据userid,videoId查询是否有该条信息，如果有，存储在FavoriteInfo，返回查询结果
	favorite := model.Favorite{}
	err := config.DB.Model(model.Favorite{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).
		First(&favorite).Error
	if err != nil {
		if "record not found" == err.Error() { //未找到匹配的记录
			log.Println("can't find data")
			return model.Favorite{}, nil
		} else { //其他错误
			log.Println(err.Error())
			return favorite, errors.New("get FavoriteInfo failed")
		}
	}
	return favorite, nil
}

// GetFavoriteVideoIdList 根据userId查询所属点赞全部videoId
func GetFavoriteVideoIdList(userId int64) ([]int64, error) {
	var FavoriteVideoIdList []int64
	err := config.DB.Model(model.Favorite{}).Where(map[string]interface{}{"user_id": userId, "status": 1}).
		Pluck("video_id", &FavoriteVideoIdList).Error
	if err != nil {
		if "record not found" == err.Error() {
			log.Println("there are no FavoriteVideoId")
			return FavoriteVideoIdList, nil
		} else {
			log.Println(err.Error())
			return FavoriteVideoIdList, errors.New("get FavoriteVideoIdList failed")
		}
	}
	return FavoriteVideoIdList, nil
}
