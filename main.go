package main

import (
	"Reborn-but-in-Go/comment"
	"Reborn-but-in-Go/favorite"
	"Reborn-but-in-Go/follow"
	"Reborn-but-in-Go/message"
	"Reborn-but-in-Go/submission"
	"Reborn-but-in-Go/user"
	"Reborn-but-in-Go/video"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	// 调用各模块接口
	video.InitVideoRouter(r)
	submission.InitSubmissionRouter(r)
	user.InitUserRouter(r)

	favorite.InitFavoriteRouter(r)
	comment.InitCommentRouter(r)
	message.InitMessageRouter(r)
	follow.InitFollowRouter(r)

	// 启动服务器并监听在 :8090 端口上
	if err := r.Run(":8090"); err != nil {
		panic("Failed to run server: " + err.Error())
	}

}

/*	// 获取服务器地址
	serverAddress := config.GetServerAddress()

	// 拼接服务器地址和接口路径
	videoListPath := "/douyin/"
	fullURL := serverAddress + videoListPath
	fmt.Print(fullURL) //防报错，忽视这条代码（会删）

	// 发起网络请求
	// ... 使用 fullURL 发起请求，处理响应等
}

// 处理设置服务器地址的请求
func SetServerAddressHandler(w http.ResponseWriter, r *http.Request) {
	//如果请求的方法不是 POST，则返回不允许的状态码
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//使用 JSON 解码器解析请求的 JSON 数据，存储在 configReq 变量中
	decoder := json.NewDecoder(r.Body)
	var configReq struct {
		ServerAddress string `json:"serverAddress"`
	}
	//如果解码出错，返回无效配置的状态码
	err := decoder.Decode(&configReq)
	if err != nil {
		http.Error(w, "Invalid config", http.StatusBadRequest)
		return
	}

	//调用 config.SetServerAddress 来设置服务器地址为传入的地址值
	config.SetServerAddress(configReq.ServerAddress)

	// 返回成功响应给前端
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server address configured"))
}
*/
