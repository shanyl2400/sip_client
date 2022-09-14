package api

import (
	"net/http"
	"sipsimclient/model"
	"sipsimclient/tools"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	contextUserNameKey = "username"
)

func (s *Server) mustAdmin(c *gin.Context) {
	token, valid := s.isTokenValid(c)
	if !valid {
		c.JSON(http.StatusUnauthorized, model.HttpErrorResponse{Message: "unauthorized"})
		c.Abort()
		return
	}
	if token.Role != model.RoleAdmin {
		c.Abort()
		c.JSON(http.StatusUnauthorized, model.HttpErrorResponse{Message: "no permission"})
		return
	}

	c.Next()
}

func (s *Server) mustLogin(c *gin.Context) {
	_, valid := s.isTokenValid(c)
	if !valid {
		c.JSON(http.StatusUnauthorized, model.HttpErrorResponse{Message: "unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}

func (s *Server) isTokenValid(c *gin.Context) (*tools.TokenInfo, bool) {
	auth := c.Request.Header.Get("Authorization")
	if len(auth) == 0 {
		return nil, false
	}
	auth = strings.Fields(auth)[1]
	// 校验token
	token, err := tools.VerifyToken(auth)
	if err != nil {
		return nil, false
	}
	c.Set(contextUserNameKey, token.UserName)

	return token, true
}
