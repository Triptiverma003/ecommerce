package middleware

import (
	token "github.com/Triptiverma003/ecommerce/tokens"
	"net/http"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func (c *gin.Context){
		ClientToken := c.Request.Header.Get("token")
		if ClientToken == ""{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization Header"})
			c.Abort()
			return
		}
		claims, err := token.ValidateToken(ClientToken)
		if err != ""{
			c.JSON(http.StatusInternalServerError , gin.H{"error" : err})
			c.Abort()
			return 
		}

		c.Set("email", claims.Email)
		c.Set("uid" , claims.Uid)
		c.Next()
	}
}