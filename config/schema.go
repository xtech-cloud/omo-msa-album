package config

type ServiceConfig struct {
	TTL      int64  `json:"ttl"`
	Interval int64  `json:"interval"`
	Address  string `json:"address"`
}

type LoggerConfig struct {
	Level string `json:"level"`
	File  string `json:"file"`
	Std   bool   `json:"std"`
}

type DBConfig struct {
	Type     string `json:"type"`
	User     string `json:"user"`
	Password string `json:"password"`
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Name     string `json:"name"`
}

type BasicConfig struct {
	SynonymMax int `ini:"synonym"`
	TagMax     int `ini:"tag"`
}

type AlbumLimit struct {
	// 一个相册最大照片个数
	MaxCount uint16 `json:"count"`
	// 一个相册最大容量（kb）
	MaxSize uint64 `json:"size"`
}

type AlbumConfig struct {
	Person AlbumLimit `json:"person"`
	Group  AlbumLimit `json:"group"`
}

type SchemaConfig struct {
	Service  ServiceConfig `json:"service"`
	Logger   LoggerConfig  `json:"logger"`
	Database DBConfig      `json:"database"`
	Basic    BasicConfig   `json:"basic"`
	Album    AlbumConfig   `json:"album"`
}
