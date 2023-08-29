package service

import (
	"Reborn-but-in-Go/config"
	"Reborn-but-in-Go/follow/dao"
	"Reborn-but-in-Go/user/model"
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

// GetFollowing 获取用户关注列表
func (*FollowService) GetFollowing(userId int64) ([]model.User, error) {
	users := make([]model.User, 1)
	// 查询出错。
	if err := config.DB.Raw("select id,`name`,"+
		"\ncount(if(tag = 'follower' and cancel is not null,1,null)) follower_count,"+
		"\ncount(if(tag = 'follow' and cancel is not null,1,null)) follow_count,"+
		"\n 'true' is_follow\nfrom\n("+
		"\nselect f1.follower_id fid,u.id,`name`,f2.cancel,'follower' tag"+
		"\nfrom follows f1 join users u on f1.user_id = u.id and f1.cancel = 0"+
		"\nleft join follows f2 on u.id = f2.user_id and f2.cancel = 0\n\tunion all"+
		"\nselect f1.follower_id fid,u.id,`name`,f2.cancel,'follow' tag"+
		"\nfrom follows f1 join users u on f1.user_id = u.id and f1.cancel = 0"+
		"\nleft join follows f2 on u.id = f2.follower_id and f2.cancel = 0\n) T"+
		"\nwhere fid = ? group by fid,id,`name`", userId).Scan(&users).Error; nil != err {
		return nil, err
	}
	// 返回关注对象列表。
	return users, nil
}

func (*FollowService) GetFollowers(userId int64) ([]model.User, error) {
	users := make([]model.User, 1)
	if err := config.DB.Raw("select T.id,T.name,T.follow_cnt follow_count,T.follower_cnt follower_count,if(f.cancel is null,'false','true') is_follow"+
		"\nfrom follows f right join"+
		"\n(select fid,id,`name`,"+
		"\ncount(if(tag = 'follower' and cancel is not null,1,null)) follower_cnt,"+
		"\ncount(if(tag = 'follow' and cancel is not null,1,null)) follow_cnt"+
		"\nfrom("+
		"\nselect f1.user_id fid,u.id,`name`,f2.cancel,'follower' tag"+
		"\nfrom follows f1 join users u on f1.follower_id = u.id and f1.cancel = 0"+
		"\nleft join follows f2 on u.id = f2.user_id and f2.cancel = 0"+
		"\nunion all"+
		"\nselect f1.user_id fid,u.id,`name`,f2.cancel,'follow' tag"+
		"\nfrom follows f1 join users u on f1.follower_id = u.id and f1.cancel = 0"+
		"\nleft join follows f2 on u.id = f2.follower_id and f2.cancel = 0"+
		"\n) T group by fid,id,`name`"+
		"\n) T on f.user_id = T.id and f.follower_id = T.fid and f.cancel = 0 where fid = ?", userId).
		Scan(&users).Error; nil != err {
		// 查询出错。
		return nil, err
	}
	// 查询成功。
	return users, nil
}
