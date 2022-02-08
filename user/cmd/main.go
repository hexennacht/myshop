package main

import (
	"context"
	"fmt"
	"github.com/hexennacht/myshop/user/config"
	repo "github.com/hexennacht/myshop/user/repository"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/hexennacht/myshop/user/handler"
	"github.com/hexennacht/myshop/user/module/user"
	userRepository "github.com/hexennacht/myshop/user/repository/user"
	"github.com/hexennacht/myshop/user/routes"
)

func main() {
	cfg := config.Read()

	db, err := initDB(cfg)
	if err != nil {
		panic(err)
	}

	if err = db.AutoMigrate(&repo.User{}); err != nil {
		panic(err)
	}

	cache := initCache(cfg)
	if err := cache.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	userRep := userRepository.NewRepository(db)

	userMod := user.NewModule(userRep, cfg.TokenLifeTime, cfg.SecretJWT)

	userHandler := handler.NewUserHandler(userMod)

	app := echo.New()
	routes.
		NewRoutes(cfg.SecretJWT, app, userHandler).
		RegisterMiddleware().
		PublicRoute().
		AuthorizedRoute()
	app.Logger.Fatal(app.Start(fmt.Sprintf(":%s", cfg.Port)))
}

func initDB(cfg *config.Config) (*gorm.DB, error) {
	if cfg.Environment == "development" {
		return gorm.Open(sqlite.Open("./database.db"), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   logger.Default.LogMode(logger.Info),
		})
	}

	return gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName)), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Info),
	})
}

func initCache(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
