package memory

import (
	"Mengbot/config"
	"Mengbot/internal/core/model"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitPgsql() error {
	dsn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s time_zone=%s",
		config.Conf.Pgsql.Host,
		config.Conf.Pgsql.Port,
		config.Conf.Pgsql.Database,
		config.Conf.Pgsql.User,
		config.Conf.Pgsql.Password,
		config.Conf.Pgsql.TimeZone,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	//向量扩展
	db.Exec("CREATE EXTENSION IF NOT EXISTS vector")
	db.AutoMigrate(&model.DiaryMessage{})
	return nil
}
