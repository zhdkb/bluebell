package models

// 定义请求的参数结构体

const (
	OrderTime = "time"
	OrderScore = "score"
)

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	Username	string	`json:"username" binding:"required"`
	Password	string	`json:"password" binding:"required"`
	RePassword	string	`json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	Username	string	`json:"username" binding:"required"`
	Password	string	`json:"password" binding:"required"`
}

// ParamVoteData 投票数据
type ParamVoteData struct {
	// userID 从请求中获取当前的用户
	PostID		string 	`json:"post_id" binding:"required"`	// 帖子id
	Direction 	int8 	`json:"direction,string" binding:"oneof=1 0 -1"` // 赞成票还是反对票还是取消投票
}

// ParamPostList 获取帖子列表的query string参数
type ParamPostList struct {
	CommunityID		int64		`json:"community_id" form:"community_id"`
	Page			int64		`json:"page" form:"page"`
	Size			int64		`json:"size" form:"size"`
	Order			string		`json:"order" form:"order"`
}
