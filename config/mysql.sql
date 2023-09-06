create table comments
(
    comment_id bigint auto_increment comment '评论id'
        primary key,
    user_id    bigint                  not null comment '用户id',
    video_id   bigint                  not null comment '视频id',
    content    varchar(255) default '' null comment '评论内容',
    create_at  datetime                not null comment '评论发布时间',
    status     tinyint(1)   default 0  null comment '评论状态(''0''为不存在''1''为存在，默认为''0'')',
    constraint comments_id_uindex
        unique (comment_id)
)
    comment '评论模块';

create table favorites
(
    id        bigint auto_increment comment '点赞id
'
        primary key,
    user_id   bigint               not null comment '用户id',
    video_id  bigint               not null comment '视频id',
    create_at datetime             not null comment '点赞时间',
    status    tinyint(1) default 0 not null comment '点赞状态(''0''为不存在''1''为存在，默认为''0'')',
    constraint favorites_id_uindex
        unique (id)
)
    comment '喜欢模块';

create index favorite_user_id_uindex
    on favorites (user_id);

create index favorite_video_id_uindex
    on favorites (video_id);

create table follows
(
    id          bigint auto_increment comment '关注信息表自增id'
        primary key,
    user_id     bigint               not null comment '当前用户',
    follower_id bigint               not null comment '关注的人',
    create_at   datetime             not null comment '关注时间',
    status      tinyint(1) default 0 not null comment '关注状态(''0''为不存在''1''为存在，默认为''0'')'
)
    comment '关注模块';

create table messages
(
    id         bigint auto_increment comment '序号'
        primary key,
    user_id    bigint       not null comment '发送人id',
    to_user_id bigint       not null comment '接收人id',
    content    varchar(255) null comment '消息内容',
    create_at  datetime     not null comment '消息时间',
    constraint messages_id_uindex
        unique (id)
)
    comment '聊天模块';

create table users
(
    id               bigint auto_increment comment '用户id'
        primary key,
    name             varchar(10)          not null comment '用户名称（不重复）',
    follow_count     bigint     default 0 null comment '关注人数',
    follower_count   bigint     default 0 null comment '粉丝人数',
    is_follow        tinyint(1) default 0 null comment '是否关注',
    password         varchar(20)          not null comment '密码',
    avatar           varchar(255)         null comment '用户头像',
    background_image varchar(255)         null comment '背景图像',
    signature        varchar(255)         null comment '个人简介',
    total_favorited  bigint     default 0 null comment '获赞数',
    work_count       bigint     default 0 null comment '作品数',
    favorite_count   bigint     default 0 not null comment '点赞数',
    constraint user_name_uindex
        unique (name),
    constraint user_user_id_uindex
        unique (id)
)
    comment '用户模块';

create table videos
(
    id             bigint auto_increment comment '视频id'
        primary key,
    user_id        bigint               not null comment '用户id',
    video_path     varchar(255)         not null comment '视频地址',
    cover_path     varchar(255)         null comment '视频封面地址',
    favorite_count bigint     default 0 null comment '点赞数',
    comment_count  bigint     default 0 null comment '评论数',
    title          varchar(50)          null comment '视频标题',
    create_at      datetime             null comment '发布时间',
    status         tinyint(1) default 0 not null,
    constraint video_video_id_uindex
        unique (id)
)
    comment '视频模块';

create index video_user_id_uindex
    on videos (user_id);

