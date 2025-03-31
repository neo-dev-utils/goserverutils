package pgUtil

import (
	"context"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Option struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Dbname   string `toml:"dbname"`
	Sslmode  string `toml:"sslmode"`
}

func ConnectDB(ctx context.Context, o Option) (db *gorm.DB, err error) {

	dsn := fmt.Sprintf("host=%s port=%s user=%s  password=%s dbname=%s sslmode=%s",
		o.Host,
		o.Port,
		o.User,
		o.Password,
		o.Dbname,
		o.Sslmode,
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		// panic(fmt.Errorf("connect db fail: %w", err))
		return nil, err
	}

	return db, nil

}
