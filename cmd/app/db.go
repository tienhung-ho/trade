package main

import (
	"log"
	"os"
	"sync"
	"time"

	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db   *gorm.DB
	once sync.Once
)

func redisConnection() *redis.Client {
	url := os.Getenv("REDIS_URL")
	if url != "" {
		url = "redis://localhost:6379"
	}

	opts, err := redis.ParseURL(url)
	if err != nil {
		log.Fatal("Error parsing URL of Redis: ", err)
	}

	rdb := redis.NewClient(opts)
	_, err = rdb.Ping(context.Background()).Result()

	if err != nil {
		log.Fatal("Error connecting to Reids: ", err)
	}

	return rdb
}

func mysqlConnection() (*gorm.DB, error) {
	once.Do(func() {
		var err error
		dsn := os.Getenv("DB_URL")

		if dsn == "" {
			log.Fatal("Environment variable DSN was empty!!")
		}

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err != nil {
			log.Fatal("Failed to connect MySQL", err)
		}

		mySQL, err := db.DB()

		if err != nil {
			log.Fatal("Failed to get database for set up pool connection", err)
		}

		mySQL.SetMaxOpenConns(100)
		mySQL.SetMaxIdleConns(10)
		mySQL.SetConnMaxLifetime(40 * time.Minute)
	})

	return db, nil
}
