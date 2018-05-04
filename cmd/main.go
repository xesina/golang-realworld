package main

import (
	"log"
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/xesina/golang-realworld/pkg/cmd"
	"github.com/xesina/golang-realworld/pkg/router"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"fmt"
	"github.com/xesina/golang-realworld/users"
	gorm2 "github.com/xesina/golang-realworld/db/gorm"
	"github.com/xesina/golang-realworld/pkg/types"
)

var (
	db *gorm.DB
)

func main() {
	app := cmd.NewCmd(
		cmd.Name("golang-realword"),
		cmd.Description("Real world app built with golang https://realworld.io"),
	)
	err := app.Init()
	if err != nil {
		log.Fatal(err)
	}

	opts := app.Options()
	fmt.Println(opts.DatabaseURI)
	db, err = gorm.Open("postgres", opts.DatabaseURI)
	db.DB()
	db.DB().Ping()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.LogMode(true)

	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&users.User{})
	db.Create(&users.User{Username: "test", Email: "test@test.com", Password: "test"})

	r := router.Engine(opts.Env)
	ur := gorm2.NewUserRepository(db)
	ui := users.NewUserInteractor(ur)
	r.GET("/users/1", func(c *gin.Context) {
		u, err := ui.Find(1)
		if err != nil {
			panic(err)
		}
		c.JSON(200, gin.H{"user": struct {
			Username string           `json:"username"`
			Email    string           `json:"email"`
			Bio      types.NullString `json:"bio"`
			Image    types.NullString `json:"image"`
		}{
			u.Username,
			u.Email,
			u.Bio,
			u.Image,
		}})
	})

	srv := &http.Server{
		Addr:    opts.ServerAddress,
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
