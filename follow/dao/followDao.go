package dao

import (
	"Reborn-but-in-Go/config"
	"Reborn-but-in-Go/follow/model"
	"log"
	"sync"
	"time"
)

// FollowDao 把dao层看成整体，把dao的curd封装在一个结构体中。
type FollowDao struct {
}

var (
	followDao  *FollowDao //操作该dao层crud的结构体变量。
	followOnce sync.Once  //单例限定，去限定申请一个followDao结构体变量。
)

// NewFollowDaoInstance 生成并返回followDao的单例对象。
func NewFollowDaoInstance() *FollowDao {
	followOnce.Do(
		func() {
			followDao = &FollowDao{}
		})
	return followDao
}

// 以下为follow层相关方法

// InsertFollowRelation 给定用户和关注对象id，建立其关系
func (*FollowDao) InsertFollowRelation(userId int64, targetId int64) (bool, error) {
	// 生成需要插入的关系结构体。
	date := time.Now()
	follow := model.Follow{
		UserID:     userId,
		FollowerID: targetId,
		CreateAt:   date,
	}
	// 插入失败，返回err.
	if err := config.DB.Select("UserId", "FollowerId", "CreateAt").Create(&follow).Error; nil != err {
		log.Println(err.Error())
		return false, err
	}
	// 插入成功
	return true, nil
}

// FindRelation 给定当前用户和目标用户id，查询follow表中相应的记录。
func (*FollowDao) FindRelation(userId int64, targetId int64) (*model.Follow, error) {
	// follow变量用于后续存储数据库查出来的用户关系。
	follow := model.Follow{}
	//当查询出现错误时，日志打印err msg，并return err.
	if err := config.DB.
		Where("UserID = ?", targetId).
		Where("FollowerID = ?", userId).
		Take(&follow).Error; nil != err {
		// 当没查到数据时，gorm也会报错。
		if "record not found" == err.Error() {
			return nil, nil
		}
		log.Println(err.Error())
		return nil, err
	}
	//正常情况，返回取到的值和空err.
	return &follow, nil
}

// GetFollowerNum 给定当前用户id，查询follow表中该用户的粉丝数。
func (*FollowDao) GetFollowerNum(userId int64) (int64, error) {
	// 用于存储当前用户粉丝数的变量
	var num int64
	// 当查询出现错误的情况，日志打印err msg，并返回err.
	if err := config.DB.
		Model(model.Follow{}).
		Where("UserID = ?", userId).
		Count(&num).Error; nil != err {
		log.Println(err.Error())
		return 0, err
	}
	// 正常情况，返回取到的粉丝数。
	return num, nil
}

// GetFollowingCnt 给定当前用户id，查询follow表中该用户关注了人数。
func (*FollowDao) GetFollowingCnt(userId int64) (int64, error) {
	// 用于存储当前用户关注了多少人。
	var cnt int64
	// 查询出错，日志打印err msg，并return err
	if err := config.DB.Model(model.Follow{}).
		Where("FollowerID = ?", userId).
		Count(&cnt).Error; nil != err {
		log.Println(err.Error())
		return 0, err
	}
	// 查询成功，返回人数。
	return cnt, nil
}

// GetFollowingIds 给定用户id，查询他关注的人的id。
func (*FollowDao) GetFollowingIds(userId int64) ([]int64, error) {
	var ids []int64
	if err := config.DB.
		Model(model.Follow{}).
		Where("follower_id = ?", userId).
		Pluck("user_id", &ids).Error; nil != err {
		// 没有关注任何人，但是不能算错。
		if "record not found" == err.Error() {
			return nil, nil
		}
		// 查询出错。
		log.Println(err.Error())
		return nil, err
	}
	// 查询成功。
	return ids, nil
}

// GetFollowersIds 给定用户id，查询他的粉丝的id。
func (*FollowDao) GetFollowersIds(userId int64) ([]int64, error) {
	var ids []int64
	if err := config.DB.
		Model(model.Follow{}).
		Where("user_id = ?", userId).
		Pluck("follower_id", &ids).Error; nil != err {
		// 没有粉丝，但是不能算错。
		if "record not found" == err.Error() {
			return nil, nil
		}
		// 查询出错。
		log.Println(err.Error())
		return nil, err
	}
	// 查询成功。
	return ids, nil
}
