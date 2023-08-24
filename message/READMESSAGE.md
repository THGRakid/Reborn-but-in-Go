# Reborn-but-in-Go\Message

## 消息模块实现逻辑

<br />

### 实现目标：

1、打开消息栏，可以看到和每个人的聊天记录，只保留最新一条

2、点进单个用户，获取与他的全部聊天记录

3、在界面里发送消息并实时更新消息记录

<br />

### 方法说明：


<br />


#### 1、QueryMessage 根据用户ID和上次消息时间获取聊天消息记录

<br />

**对应接口**

`/douyin/message/chat/ `  `GET`

当前登录用户和其他指定用户的聊天消息记录

<br />

**接收数据**

token string  // 用户鉴权token

to_user_id int64  // 对方用户id

pre_msg_time int64  //上次最新消息的时间

<br />

**返回数据**

status_code int32  // 状态码，0-成功，其他值-失败

status_msg string  // 返回状态描述

message_list Message  // 消息列表


<br />


**Message 结构体需包含**

id int64  // 消息id

rto_user_id int64  // 该消息接收者的id

from_user_id int64  // 该消息发送者的id

content string   // 消息内容

create_time string  // 消息创建时间
}


<br />



#### 2、SendMessage 发送消息

<br />


**对应接口**

`/douyin/message/action `  `POST`

登录用户对消息的相关操作，目前只支持消息发送

<br />


**接收数据**

token string  // 用户鉴权token

to_user_id int64  // 对方用户id

action_type int32  // 1-发送消息

content string  // 消息内容

<br />


**返回数据**

status_code int32  // 状态码，0-成功，其他值-失败

status_msg string  // 返回状态描述

