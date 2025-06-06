package routes

import (
	"github.com/Triptiverma003/ecommerce/controllers"
	"github.com/gin-gonic/gin"
)


func userRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/user/signup" , controllers.Signup())
	incomingRoutes.POST("/user/login" , controllers.Login())
	incomingRoutes.POST("/admin/addproduct" , controllers.ProductViewerAdmin())
	incomingRoutes.GET("/user/productview" , controllers.SearchProduct())
	incomingRoutes.GET("/user/search" , SearchProductByQuery())
}
