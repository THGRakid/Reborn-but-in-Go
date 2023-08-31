package model

import "time"

type Video struct {
	Id            int64     `json:"id,omitempty"`   //视频id
	UserId        int64     `json:"author"`         //作者id
	VideoPath     string    `json:"play_url"`       //视频地址
	CoverPath     string    `json:"cover_url"`      //视频封面地址
	FavoriteCount int64     `json:"favorite_count"` //点赞数
	CommentCount  int64     `json:"comment_count"`  //评论数
	Title         string    `json:"title"`          //视频标题
	CreateAt      time.Time //投稿时间
	Status        int32     //待审核0 发布成功1 审核失败2 下架删除3
}

func (Video) TableName() string {
	return "videos"
}
