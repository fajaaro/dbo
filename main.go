package main

import (
	"log"

	"github.com/fajaaro/dbo/app"
	"github.com/fajaaro/dbo/app/controllers"
	"github.com/fajaaro/dbo/app/migrations"
	"github.com/fajaaro/dbo/app/routers"
)

func main() {
	db := app.InitDb()
	err := migrations.AutoMigrate(db)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Migration completed successfully.")

	r := routers.SetupRouter(*controllers.AuthController(), *controllers.OrderController(), *controllers.CustomerController())
	_ = r.Run(":8080")
}
