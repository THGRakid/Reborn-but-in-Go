create table follow
(
id          bigint auto_increment comment '关注信息表自增id'
primary key,
user_id     bigint               not null comment '当前用户',
follower_id bigint               not null comment '关注的人',
creat_at    datetime             not null comment '关注时间',
status      tinyint(1) default 0 not null comment '关注状态(''0''为不存在''1''为存在，默认为''0'')'
)
comment '关注模块';

create table user
(
id              bigint auto_increment comment '用户id'
primary key,
name            varchar(10)          not null comment '用户名称（不重复）',
follower        bigint     default 0 null comment '粉丝人数',
following       bigint     default 0 null comment '关注人数',
password        varchar(20)          not null comment '密码',
avatar          varchar(255)         null comment '用户头像',
background      varchar(255)         null comment '背景图像',
introduce       varchar(255)         null comment '个人简介',
favorited_count bigint     default 0 null comment '获赞数',
work_count      bigint     default 0 null comment '作品数',
favorite_count  bigint     default 0 not null comment '点赞数',
time            datetime             not null comment '用户创建时间',
status          tinyint(1) default 0 not null comment '用户状态',
constraint user_name_uindex
unique (name),
constraint user_user_id_uindex
unique (id)
)
comment '用户模块';

create table messages
(
id         bigint auto_increment comment '序号'
primary key,
user_id    bigint               not null comment '发送人id',
to_user_id bigint               not null comment '接收人id',
content    varchar(255)         null comment '消息内容',
create_at  datetime             not null comment '消息时间',
status     tinyint(1) default 0 not null comment '消息状态(''0''为不存在''1''为存在，默认为''0'')',
constraint message_receive_id_uindex
unique (to_user_id),
constraint message_send_id_uindex
unique (user_id),
constraint message_user_id_fk
foreign key (user_id) references user (id),
constraint message_user_id_fk_2
foreign key (to_user_id) references user (id)
)
comment '聊天模块';

create table video
(
id             bigint auto_increment comment '视频id'
primary key,
user_id        bigint               not null comment '用户id',
video_path     varchar(255)         not null comment '视频地址',
cover_path     varchar(255)         null comment '视频封面地址',
favorite_count bigint     default 0 null comment '点赞数',
comment_count  bigint     default 0 null comment '评论数',
title          varchar(50)          null comment '视频标题',
time           datetime             null comment '发布时间',
status         tinyint(1) default 0 not null,
constraint video_user_id_uindex
unique (user_id),
constraint video_video_id_uindex
unique (id),
constraint video_user_id_fk
foreign key (user_id) references user (id)
)
comment '视频模块';

create table comment
(
user_id  bigint                  not null comment '用户id'
primary key,
video_id bigint                  not null comment '视频id',
content  varchar(255) default '' null comment '评论内容',
time     datetime                not null comment '评论发布时间',
status   tinyint(1)   default 0  null comment '评论状态(''0''为不存在''1''为存在，默认为''0'')',
constraint comment_user_id_fk
foreign key (user_id) references user (id),
constraint comment_video_id_fk
foreign key (video_id) references video (id)
)
comment '评论模块';

create table favorite
(
user_id  bigint               not null comment '用户id'
primary key,
video_id bigint               not null comment '视频id',
time     datetime             not null comment '点赞时间',
status   tinyint(1) default 0 not null comment '点赞状态(''0''为不存在''1''为存在，默认为''0'')',
constraint favorite_user_id_uindex
unique (user_id),
constraint favorite_video_id_uindex
unique (video_id),
constraint favorite_user_id_fk
foreign key (user_id) references user (id),
constraint favorite_video_id_fk
foreign key (video_id) references video (id)
)
comment '喜欢模块';

