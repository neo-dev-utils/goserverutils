package cfgUtil

type LogOpt struct {
	IsConsole  bool   `json:"isConsole" yaml:"isConsole"`
	IsFile     bool   `json:"isFile" yaml:"isFile"`
	SavePath   string `json:"savePath" yaml:"savePath"`     // 保存路径
	MaxSize    int64  `json:"maxSize" yaml:"maxSize"`       // 文件大小限制,单位MB
	MaxBackups int64  `json:"maxBackups" yaml:"maxBackups"` // 最大保留日志文件数量
	MaxAge     int64  `json:"maxAge" yaml:"maxAge"`         // 日志文件保留天数
	IsCompress bool   `json:"isCompress" yaml:"isCompress"` // 是否压缩处理
}
