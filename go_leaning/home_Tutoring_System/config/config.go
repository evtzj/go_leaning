package config

var JWTSecret = []byte("evtzj的家教管理系统") // JWT 签名密钥

const (
	DBPath = "jiajiao.db" // SQLite 数据库文件名
	Port   = ":8080"      // 服务端口
)
