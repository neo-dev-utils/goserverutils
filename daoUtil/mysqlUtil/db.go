package mysqlUtil

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Option struct {
	Host     string `json:"host" yaml:"host" toml:"host"`
	Port     string `json:"port" yaml:"port" toml:"port"`
	Dbname   string `json:"dbname" yaml:"dbname" toml:"dbname"`
	UserName string `json:"username" yaml:"username" toml:"username"`
	Password string `json:"password" yaml:"password" toml:"password"`
	Charset  string `json:"charset" yaml:"charset" toml:"charset"`
	Timezone string `json:"timezone" yaml:"timezone" toml:"timezone"`
	Sslmode  string `json:"sslmode" yaml:"sslmode" toml:"sslmode"`
}

func ConnectDB(ctx context.Context, o Option) (db *gorm.DB, err error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		o.UserName,
		o.Password,
		o.Host,
		o.Port,
		o.Dbname,
		o.Charset,
		o.Timezone,
	)

	instance, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Connect db error:\n" + err.Error())
	}

	pool, _ := instance.DB()
	pool.SetMaxIdleConns(2)
	pool.SetMaxOpenConns(4)
	pool.SetConnMaxLifetime(3 * time.Hour)
	db = instance

	return db, nil
}
