package controller

import "Reborn-but-in-Go/follow/service"

type FollowController struct {
	FollowService service.FollowService
}

func NewFollowController(followService *service.FollowService) *FollowController {
	return &FollowController{
		FollowService: followService,
	}
}
