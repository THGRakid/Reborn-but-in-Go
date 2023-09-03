package dao

import (
	"Reborn-but-in-Go/config"
	"Reborn-but-in-Go/follow/model"
	"fmt"
	"gorm.io/gorm"
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

// InsertFollowRelation 给定用户和关注对象id，建立其关系
func (*FollowDao) InsertFollowRelation(userId int64, targetId int64) (bool, error) {
	// 相关信息预处理
	currentTime := time.Now()
	var existingFollow model.Follow
	// 尝试查询是否已经存在关注信息
	if err := config.DB.Where("user_id = ? AND follower_id = ?", userId, targetId).
		First(&existingFollow).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果记录不存在，插入新的关注信息
			newFollow := model.Follow{
				UserId:     userId,
				FollowerId: targetId,
				CreateAt:   currentTime,
				Status:     0,
			}
			config.DB.Create(&newFollow)
			log.Println("首次关注成功")
			return true, nil
		} else {
			// 处理其他查询错误
			fmt.Println("查询关注信息时发生错误:", err)
			return false, err
		}
	} else {
		// 如果存在关注信息，更新状态和时间
		existingFollow.Status = 0
		existingFollow.CreateAt = currentTime
		config.DB.Save(&existingFollow)
		log.Println("重新关注成功")
		return true, nil
	}
}

// DeleteFollowRelation 给定用户和取消关注对象id，取消其关系
func (*FollowDao) DeleteFollowRelation(userId int64, targetId int64) (bool, error) {
	// 定义要更新的数据
	updateData := map[string]interface{}{
		"Status":   1,          // 更新状态为1
		"CreateAt": time.Now(), // 更新创建时间为当前时间
	}
	// 具体取关操作
	if err := config.DB.Model(&model.Follow{}).
		Where("user_id = ? AND follower_id = ?", userId, targetId).
		Updates(updateData).Error; err != nil {
		// 更新失败
		log.Println("取关失败")
		return false, nil
	} else {
		// 更新成功
		log.Println("取关成功")
		return true, nil
	}
}

// FindRelation 给定当前用户和目标用户id，查询follow表中相应的记录。
func (*FollowDao) FindRelation(userId int64, targetId int64) (*model.Follow, error) {
	// follow变量用于后续存储数据库查出来的用户关系。
	follow := model.Follow{}
	result := config.DB.Where("user_id = ? AND follower_id = ? AND status = 0", userId, targetId).
		First(&follow)
	if result.Error != nil {
		// 处理查询错误
		if result.Error == gorm.ErrRecordNotFound {
			// 如果记录不存在，说明当前用户没有关注目标用户
			fmt.Println("当前用户没有关注目标用户")
			return nil, nil
		} else {
			fmt.Println("查询关注信息时发生错误:", result.Error)
			return nil, nil
		}
	} else {
		// 如果查询成功，说明当前用户关注了目标用户
		fmt.Println("当前用户已关注目标用户")
		return &follow, nil
	}
}

// GetFollowingIds 给定用户id，查询他关注的人的id。
func (*FollowDao) GetFollowingIds(userId int64) ([]int64, error) {
	var ids []int64
	if err := config.DB.
		Model(model.Follow{}).
		Where("user_id = ? AND status = 0", userId).
		Pluck("follower_id", &ids).Error; nil != err {
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
		Where("follower_id = ? AND status = 0", userId).
		Pluck("user_id", &ids).Error; nil != err {
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
