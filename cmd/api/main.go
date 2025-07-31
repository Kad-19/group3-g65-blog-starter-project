package main

import (
	"g3-g65-bsp/delivery/controller"
	"g3-g65-bsp/delivery/route"
	config "g3-g65-bsp/infrastructure/database"
	"g3-g65-bsp/repository"
	"g3-g65-bsp/usecase"
)

func main() {

    db := config.InitMongoDB().Database("blogdb")
    blogCollection := db.Collection("blogs")

    // Initialize repository, usecase, controller
    blogRepo := repository.NewBlogRepository(blogCollection)
    blogUsecase := usecase.NewBlogUsecase(blogRepo)
    blogController := controller.NewBlogController(blogUsecase)

    // Initialize router
    r := route.NewRouter()
	route.BlogRouter(r, blogController)


   // Start the server on port 8080
	if err := r.Run("localhost:8080"); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
