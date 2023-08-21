package service

import (
	"Reborn-but-in-Go/follow/dao"
)

//	// AddFollowRelation 当前用户关注目标用户
//	AddFollowRelation(userId int64, targetId int64) (bool, error)
//	// DeleteFollowRelation 当前用户取消对目标用户的关注
//	DeleteFollowRelation(userId int64, targetId int64) (bool, error)
//	// GetFollowing 获取当前用户的关注列表
//	//GetFollowing(userId int64) ([]User, error)
//	// GetFollowers 获取当前用户的粉丝列表
//	//GetFollowers(userId int64) ([]User, error)

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

	if nil != err {
		return false, err
	}

	if nil == relation {
		return false, nil
	}

	return true, nil
}

// GetFollowingNum 根据用户id来查询该用户关注数目
func (*FollowService) GetFollowingNum(userId int64) (int64, error) {
	// SQL中查询
	ids, err := dao.NewFollowDaoInstance().GetFollowingIds(userId)

	if nil != err {
		return 0, err
	}

	return int64(len(ids)), err
}

// GetFollowerNum 根据用户id来查询该用户的粉丝数目
func (*FollowService) GetFollowerNum(userId int64) (int64, error) {
	// SQL中查询
	ids, err := dao.NewFollowDaoInstance().GetFollowersIds(userId)

	if nil != err {
		return 0, err
	}

	return int64(len(ids)), err
}
