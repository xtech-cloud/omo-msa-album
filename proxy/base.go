package proxy

type PhotocopySlot struct {
	// 第几页
	Page uint8 `json:"page" bson:"page"`
	// 第几个
	Index uint8 `json:"index" bson:"index"`
	// 每页可编辑角色,母版创建者或者克隆者
	Role       uint8  `json:"role" bson:"role"`
	Background string `json:"background" bson:"background"`
	Asset      string `json:"asset" bson:"asset"`
	Remark     string `json:"remark" bson:"remark"`
}
