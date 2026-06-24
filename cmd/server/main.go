package main

import (
	"log"
	"net/http"
	"os"

	app "github.com/ladislaojs/team-task-manager-test/internal/http"
	"github.com/ladislaojs/team-task-manager-test/internal/http/handler"
	mysqlrepo "github.com/ladislaojs/team-task-manager-test/internal/repository/mysql"
	"github.com/ladislaojs/team-task-manager-test/internal/service"
	"github.com/ladislaojs/team-task-manager-test/pkg/mysql"
)

func main() {
	// TODO: JWT Secret

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

	userService := service.NewUserService(userRepository)
	teamService := service.NewTeamService(teamRepository)
	taskService := service.NewTaskService(taskRepository)

	userHandler := handler.NewUserHandler(userService)
	teamHandler := handler.NewTeamHandler(teamService)
	taskHandler := handler.NewTaskHandler(taskService)

	router := app.NewRouter(
		userHandler,
		teamHandler,
		taskHandler,
	)
	port := os.Getenv("APP_PORT")

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
