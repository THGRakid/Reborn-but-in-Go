package dao

import (
	"errors"
	"log"
	"sync"
	"time"

	"Reborn-but-in-Go/config"
	"Reborn-but-in-Go/favorite/model"
	userMod "Reborn-but-in-Go/user/model"
	vidMod "Reborn-but-in-Go/video/model"
	"gorm.io/gorm"
)

const (
	ActionTypeLike   = 1
	ActionTypeUnlike = 2
)

const (
	ErrGetFavoriteUserIdList  = "get FavoriteUserIdList failed"
	ErrGetVideoFavoriteCount  = "get video favorite_count failed"
	ErrInsertFavorite         = "insert data fail"
	ErrAddVideoFavoriteCount  = "add video favorite_count fail"
	ErrUnsupportedActionType  = "不支持的操作类型"
	ErrGetFavoriteInfo        = "get FavoriteInfo failed"
	ErrGetFavoriteVideoIdList = "get FavoriteVideoIdList failed"
	ErrGetUserInfo            = "get userInfo failed"
)

type FavoriteDao struct{}

var (
	favoriteDao  *FavoriteDao
	favoriteOnce sync.Once
)

func NewFavoriteDaoInstance() *FavoriteDao {
	favoriteOnce.Do(func() {
		favoriteDao = &FavoriteDao{}
	})
	return favoriteDao
}

func (*FavoriteDao) GetFavoriteUserIdList(videoId int64) ([]int64, error) {
	var FavoriteUserIdList []int64
	err := config.DB.
		Model(model.Favorite{}).
		Where(map[string]interface{}{"video_id": videoId, "status": ActionTypeLike}).
		Pluck("user_id", &FavoriteUserIdList).Error
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New(ErrGetFavoriteUserIdList)
	}
	return FavoriteUserIdList, nil
}

func GetFavoriteCount(videoId int64) (int64, error) {
	var video vidMod.Video
	err := config.DB.Model(vidMod.Video{}).Where(map[string]interface{}{"id": videoId}).First(&video).Error
	if err != nil {
		log.Println(err.Error())
		return 0, errors.New(ErrGetVideoFavoriteCount)
	}
	return video.FavoriteCount, nil
}

func (*FavoriteDao) InsertFavorite(FavoriteData model.Favorite) error {
	err := config.DB.Model(model.Favorite{}).Create(&FavoriteData).Error
	if err != nil {
		log.Println(err.Error())
		return errors.New(ErrInsertFavorite)
	}
	log.Printf("添加点赞记录成功")
	verr := config.DB.Model(vidMod.Video{}).Where(map[string]interface{}{"id": FavoriteData.VideoId}).
		Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
	if verr != nil {
		log.Println(verr.Error())
		return errors.New(ErrAddVideoFavoriteCount)
	}
	log.Printf("点赞数量加一")
	return nil
}

func (*FavoriteDao) UpdateFavorite(userId int64, videoId int64, actionType int8) error {
	if actionType == ActionTypeUnlike {
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
		return errors.New(ErrUnsupportedActionType)
	}
	return nil
}

func (*FavoriteDao) UpdateOrInsertFavorite(userId int64, videoId int64, actionType int8) error {
	existingFavorite, err := favoriteDao.GetFavoriteInfo(userId, videoId)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	if actionType == ActionTypeLike && existingFavorite.UserId != userId {
		newFavorite := model.Favorite{
			UserId:   userId,
			VideoId:  videoId,
			Status:   ActionTypeLike,
			CreateAt: time.Now(),
		}
		err := favoriteDao.InsertFavorite(newFavorite)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	} else {
		err := favoriteDao.UpdateFavorite(userId, videoId, actionType)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	}
	return nil
}

func (*FavoriteDao) GetFavoriteInfo(userId int64, videoId int64) (model.Favorite, error) {
	favorite := model.Favorite{}
	err := config.DB.Model(model.Favorite{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).
		First(&favorite).Error
	if err != nil {
		if "record not found" == err.Error() {
			log.Println("can't find data")
			return model.Favorite{}, nil
		} else {
			log.Println(err.Error())
			return favorite, errors.New(ErrGetFavoriteInfo)
		}
	}
	return favorite, nil
}

func GetFavoriteVideoIdList(userId int64) ([]int64, error) {
	var FavoriteVideoIdList []int64
	err := config.DB.Model(model.Favorite{}).Where(map[string]interface{}{"user_id": userId, "status": ActionTypeLike}).
		Pluck("video_id", &FavoriteVideoIdList).Error
	if err != nil {
		if "record not found" == err.Error() {
			log.Println("there are no FavoriteVideoId")
			return FavoriteVideoIdList, nil
		} else {
			log.Println(err.Error())
			return FavoriteVideoIdList, errors.New(ErrGetFavoriteVideoIdList)
		}
	}
	return FavoriteVideoIdList, nil
}

func GetUserByID(userId int64) (userMod.User, error) {
	user := userMod.User{}
	err := config.DB.Model(userMod.User{}).Where(map[string]interface{}{"id": userId}).
		First(&user).Error
	if err != nil {
		if "record not found" == err.Error() {
			log.Println("can't find data")
			return userMod.User{}, nil
		} else {
			log.Println(err.Error())
			return user, errors.New(ErrGetUserInfo)
		}
	}
	return user, nil
}
