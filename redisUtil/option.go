package redisUtil

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Option struct {
	Network     string        `json:"network" yaml:"network" toml:"network"`
	Ip          string        `json:"ip" yaml:"ip" toml:"ip"`
	Port        string        `json:"port" yaml:"port" toml:"port"`
	Pwd         string        `json:"pwd" yaml:"pwd" toml:"pwd"`
	Db          int           `json:"db" yaml:"db" toml:"db"`
	MaxIdle     int           `json:"maxidle" yaml:"maxidle" toml:"maxidle"`
	MaxActive   int           `json:"maxactive" yaml:"maxactive" toml:"maxactive"`
	IdleTimeout time.Duration `json:"idletimeout" yaml:"idletimeout" toml:"idletimeout"`
	KeepAlive   time.Duration `json:"keepalive" yaml:"keepalive" toml:"keepalive"`
	SSLMode     bool          `json:"sslmode" yaml:"sslmode" toml:"sslmode"`
}

func ConnectRedis(opt Option) *redis.Client {

	client := redis.NewClient(&redis.Options{
		//连接信息
		Network:  opt.Network,             //网络类型，tcp or unix，默认tcp
		Addr:     opt.Ip + ":" + opt.Port, //主机名+冒号+端口，默认localhost:6379
		Password: opt.Pwd,                 //密码
		DB:       opt.Db,                  // redis数据库index

		//连接池容量及闲置连接数量
		PoolSize:     15, // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
		MinIdleConns: 10, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。

		//超时
		DialTimeout:  5 * time.Second, //连接建立超时时间，默认5秒。
		ReadTimeout:  3 * time.Second, //读超时，默认3秒， -1表示取消读超时
		WriteTimeout: 3 * time.Second, //写超时，默认等于读超时
		PoolTimeout:  4 * time.Second, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic("Connect rdb error:\n" + err.Error())
	}

	return client
}
