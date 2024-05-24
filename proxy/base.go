package proxy

type DurationInfo struct {
	Start int64 `json:"start" bson:"start"`
	Stop  int64 `json:"stop" bson:"stop"`
}

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

type ContactInfo struct {
	Name    string `json:"name" bson:"name"`
	Phone   string `json:"phone" bson:"phone"`
	Address string `json:"address" bson:"address"`
	Remark  string `json:"remark" bson:"remark"`
}

type SheetPage struct {
	UID    string `json:"uid" bson:"uid"`
	Weight uint32 `json:"weight" bson:"weight"`
}

type SlotInfo struct {
	Index    uint32  `json:"index" bson:"index"`
	Type     uint8   `json:"type" bson:"type"`
	Position Vector2 `json:"position" bson:"position"`
	Size     Vector2 `json:"size" bson:"size"`
}

type FolderContent struct {
	Type uint32 `json:"type" bson:"type"` //资源类型0=素材， 1=实体
	UID  string `json:"uid" bson:"uid"`   //资源UID
}

type StyleRelate struct {
	Entity string `json:"entity" bson:"entity"` //关联的实体
	Way    uint8  `json:"way" bson:"way"`       //获取的路径
}

type PageContents struct {
	Slot     uint32   `json:"slot" bson:"slot"`         //构图的位置
	Type     uint32   `json:"type" bson:"type"`         //资源类型0=素材， 1=实体
	Way      uint8    `json:"way" bson:"way"`           //播放方式
	Interval uint32   `json:"interval" bson:"interval"` //播放间隔
	List     []string `json:"list" bson:"list"`
}

type Vector2 struct {
	X int32 `json:"x" bson:"x"`
	Y int32 `json:"y" bson:"y"`
}

type StyleSlot struct {
	Key    string `json:"key" bson:"key"`
	X      int32  `json:"x" bson:"x"`
	Y      int32  `json:"y" bson:"y"`
	Width  uint32 `json:"width" bson:"width"`
	Height uint32 `json:"height" bson:"height"`
	Bold   uint32 `json:"bold" bson:"bold"`
	Size   uint32 `json:"size" bson:"size"`
}
