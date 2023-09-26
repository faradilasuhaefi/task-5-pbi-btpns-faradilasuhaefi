package router

import (
	"final-project-pbi-btpns/controllers"
	"final-project-pbi-btpns/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Public routes (without authentication middleware)
	public := r.Group("/api")
	{
		public.POST("/users/register", controllers.Register)
		public.POST("/users/login", controllers.Login)
	}

	// Protected routes (with authentication middleware)
	protected := r.Group("/api")
	protected.Use(middlewares.Authenticate())
	{
		// Routes for the resource: user
		protected.GET("/users/:userId", controllers.GetUserByID)
		protected.PUT("/users/:userId", controllers.UpdateUser)
		protected.DELETE("/users/:userId", controllers.DeleteUser)

		// Routes for the resource: photo
		protected.POST("/photos", controllers.CreatePhoto)
		protected.GET("/photos", controllers.GetPhotoList)
		protected.PUT("/photos/:photoId", controllers.UpdatePhoto)
		protected.GET("/photos/:photoId", controllers.GetPhotoByID)
		protected.DELETE("/photos/:photoId", controllers.DeletePhoto)
	}

	return r
}
