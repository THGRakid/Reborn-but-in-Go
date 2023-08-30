package model

import "time"

type Video struct {
	Id            int64     //视频id
	UserId        int64     //作者id
	VideoPath     string    //视频地址
	CoverPath     string    //视频封面地址
	FavoriteCount int64     //点赞数
	CommentCount  int64     //评论数
	Title         string    //视频标题
	Time          time.Time //投稿时间
	Status        int32     //待审核0 发布成功1 审核失败2 下架删除3
}

func (Video) TableName() string {
	return "videos"
}
