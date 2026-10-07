package main

import (
	"log"
	"time"

	"groundTurn/src/config"
	"groundTurn/src/repositories"
	"groundTurn/src/routes"
	"groundTurn/src/seed"

	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	// 等待数据库 healthy（compose 已保证，本地直连时重试兜底）。
	db, err := connectWithRetry(cfg, 30, 2*time.Second)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	if err := seed.Run(db, cfg); err != nil {
		log.Printf("seed skipped: %v", err)
	}

	r := routes.Start(":"+cfg.Port, cfg, db)
	log.Printf("ground-turn backend listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

func connectWithRetry(cfg config.Config, attempts int, interval time.Duration) (*gorm.DB, error) {
	var last error
	for i := 0; i < attempts; i++ {
		db, err := repositories.Connect(cfg)
		if err == nil {
			return db, nil
		}
		last = err
		time.Sleep(interval)
	}
	return nil, last
}
