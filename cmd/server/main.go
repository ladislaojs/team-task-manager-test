package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	app "github.com/ladislaojs/team-task-manager-test/internal/http"
	"github.com/ladislaojs/team-task-manager-test/internal/http/handler"
	mysqlrepo "github.com/ladislaojs/team-task-manager-test/internal/repository/mysql"
	"github.com/ladislaojs/team-task-manager-test/internal/service"
	"github.com/ladislaojs/team-task-manager-test/pkg/cache"
	"github.com/ladislaojs/team-task-manager-test/pkg/mysql"
)

func main() {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET env variable is required")
	}

	db, err := mysql.Connect(mysql.Config{
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     os.Getenv("MYSQL_PORT"),
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Database: os.Getenv("MYSQL_DATABASE"),
	})
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()

	userRepository := mysqlrepo.NewUserRepository(db)
	teamRepository := mysqlrepo.NewTeamRepository(db)
	taskRepository := mysqlrepo.NewTaskRepository(db)

	userService := service.NewUserService(userRepository, jwtSecret)
	teamService := service.NewTeamService(teamRepository, userRepository)
	taskService := service.NewTaskService(taskRepository)

	userHandler := handler.NewUserHandler(userService)
	teamHandler := handler.NewTeamHandler(teamService)
	taskHandler := handler.NewTaskHandler(taskService)

	redisClient, err := cache.NewRedisClient(cache.Config{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       os.Getenv("REDIS_DB"),
	})
	if err != nil {
		log.Fatalf("redis connection failed: %v", err)
	}

	maxRequestsPerMinute, err := strconv.Atoi(os.Getenv("MAX_REQUESTS_PER_MINUTE"))
	if err != nil {
		maxRequestsPerMinute = 0
	}

	router := app.NewRouter(
		jwtSecret,
		redisClient,
		maxRequestsPerMinute,
		userHandler,
		teamHandler,
		taskHandler,
	)
	port := os.Getenv("APP_PORT")

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
