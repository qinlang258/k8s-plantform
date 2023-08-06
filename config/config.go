package config

import "time"

const (
	ListenAddr     = "0.0.0.0:9090"
	PodLogTailLine = 1000

	DbType = "mysql"
	DbHost = "127.0.0.1"
	DbPort = 3306
	DbName = "k8s_dashboard"
	DbUser = "root"
	DbPass = "ql2252528"
	// 打印mysql debug的sql日志
	LogMode = false
	// 连接池配置
	MaxIdleConns = 10               // 最大空闲连接
	MaxOpenConns = 100              // 最大连接数
	MaxLifeTime  = 30 * time.Second // 会话时间
)
