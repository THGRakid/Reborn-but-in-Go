package service

// 未完成
// FollowService 包括用户关系接口、用户关系中的方法
type FollowService interface {

	// IsFollowing 根据当前用户id和目标用户id来判断当前用户是否关注了目标用户
	IsFollowing(userId int64, targetId int64) (bool, error)
	// GetFollowerCnt 根据用户id来查询用户被多少其他用户关注
	GetFollowerCnt(userId int64) (int64, error)
	// GetFollowingCnt 根据用户id来查询用户关注了多少其它用户
	GetFollowingCnt(userId int64) (int64, error)

	// AddFollowRelation 当前用户关注目标用户
	AddFollowRelation(userId int64, targetId int64) (bool, error)
	// DeleteFollowRelation 当前用户取消对目标用户的关注
	DeleteFollowRelation(userId int64, targetId int64) (bool, error)
	// GetFollowing 获取当前用户的关注列表
	GetFollowing(userId int64) ([]User, error)
	// GetFollowers 获取当前用户的粉丝列表
	GetFollowers(userId int64) ([]User, error)
}

type User struct {
	Id             int64  `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	FollowCount    int64  `json:"follow_count"`
	FollowerCount  int64  `json:"follower_count"`
	IsFollow       bool   `json:"is_follow"`
	TotalFavorited int64  `json:"total_favorited,omitempty"`
	FavoriteCount  int64  `json:"favorite_count,omitempty"`
}
