package service

import (
	"Reborn-but-in-Go/config"
	"Reborn-but-in-Go/follow/dao"
	"Reborn-but-in-Go/user/model"
	"fmt"
	"log"
)

//	// AddFollowRelation 当前用户关注目标用户
//	AddFollowRelation(userId int64, targetId int64) (bool, error)
//	// DeleteFollowRelation 当前用户取消对目标用户的关注
//	DeleteFollowRelation(userId int64, targetId int64) (bool, error)

// FollowService Service层
type FollowService struct {
	FollowDao *dao.FollowDao
}

func NewFollowService(followDao *dao.FollowDao) *FollowService {
	return &FollowService{
		FollowDao: followDao,
	}
}

// IsFollowing 根据当前用户id和目标用户id来判断当前用户是否关注了目标用户
func (*FollowService) IsFollowing(userId int64, targetId int64) (bool, error) {
	// SQL中查询
	relation, err := dao.NewFollowDaoInstance().FindRelation(userId, targetId)

	if nil != err || nil == relation {
		return false, err
	} else {
		return true, nil
	}
}

// GetFollowingNum 根据用户id来查询该用户关注数目
func (*FollowService) GetFollowingNum(userId int64) (int64, error) {
	// SQL中查询
	ids, err1 := dao.NewFollowDaoInstance().GetFollowingIds(userId)

	if nil != err1 {
		return 0, err1
	} else {
		newFollowCount := int64(len(ids))
		err2 := config.DB.Model(&model.User{}).
			Where("id = ?", userId).
			Update("follow_count", newFollowCount).Error
		log.Println("更新数据库成功")
		return newFollowCount, err2
	}
}

// GetFollowerNum 根据用户id来查询该用户的粉丝数目
func (*FollowService) GetFollowerNum(userId int64) (int64, error) {
	// SQL中查询
	ids, err1 := dao.NewFollowDaoInstance().GetFollowersIds(userId)

	if nil != err1 {
		return 0, err1
	} else {
		newFollowerCount := int64(len(ids))
		err2 := config.DB.Model(&model.User{}).
			Where("id = ?", userId).
			Update("follower_count", newFollowerCount).Error
		log.Println("更新数据库成功")
		return newFollowerCount, err2
	}
}

// GetFollowing 获取用户关注列表
func (f *FollowService) GetFollowing(targetID int64) ([]model.User, error) {
	// 查询关注的对象的信息
	var followedUsers []model.User
	if err := config.DB.Table("users").
		Joins("INNER JOIN follows ON users.id = follows.follower_id").
		Where("follows.user_id = ? AND follows.status = 0", targetID).
		Select("users.id, users.name, users.follow_count, users.follower_count").
		Find(&followedUsers).Error; err != nil {
		// 处理查询错误
		fmt.Println("查询关注的对象信息时发生错误:", err)
		return nil, err
	} else {
		// 查询成功，followedUsers 包含了目标用户关注的对象的信息
		for i := range followedUsers {
			result, _ := f.IsFollowing(targetID, followedUsers[i].Id)
			followedUsers[i].IsFollow = result
		}
		//fmt.Println("目标用户关注的对象信息:", followedUsers)
		return followedUsers, nil
	}
}

// GetFollowers 获取用户粉丝列表
func (f *FollowService) GetFollowers(targetID int64) ([]model.User, error) {
	// 查询粉丝的信息
	var followers []model.User
	if err := config.DB.Table("users").
		Joins("INNER JOIN follows ON users.id = follows.user_id").
		Where("follows.follower_id = ? AND follows.status = 0", targetID).
		Select("users.id, users.name, users.follow_count, users.follower_count").
		Find(&followers).Error; err != nil {
		// 处理查询错误
		fmt.Println("查询粉丝信息时发生错误:", err)
		return nil, err
	} else {
		// 查询成功，followedUsers 包含了目标用户关注的对象的信息
		for i := range followers {
			result, _ := f.IsFollowing(targetID, followers[i].Id)
			followers[i].IsFollow = result
		}
		//fmt.Println("目标用户粉丝信息:", followers)
		return followers, nil
	}
}
