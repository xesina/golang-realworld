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
	"github.com/xesina/golang-realworld/models"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"fmt"
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

	db.AutoMigrate(&models.User{})

	r := router.Engine(opts.Env)
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
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
