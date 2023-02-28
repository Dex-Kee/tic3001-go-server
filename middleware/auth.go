package middleware

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"tic3001-go-server/common/constant"
	"tic3001-go-server/common/dto"
	"tic3001-go-server/service"
)

type authFilter struct{}

var AuthFileter = new(authFilter)

func (filter authFilter) ValidateResource(c *gin.Context) {
	resourcePath := c.FullPath()

	// login api, pass
	if resourcePath == "/api/auth/login" {
		c.Next()
		return
	}

	// for other api, check the role & resource
	// check token
	token := c.GetHeader("token")
	claims, err := service.AuthService.ParserToken(token)
	if err != nil {
		log.Error("error when parse token: ", err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.NewResponseDto(constant.RespCodeUnauthorized,
			constant.RespMsgUnauthenticated, ""))
		return
	}

	// log.Infof("parsed claims: %##v", claims)

	// if role is admin, no resource check is required
	role := claims.Auth.Role
	if role == "admin" {
		c.Next()
		return
	}

	// for other roles, check accessible resource
	resourceMap := service.AuthService.FindAccessibleResourceByRole(role)
	if ok, _ := resourceMap[resourcePath]; !ok {
		log.Infof("request uri: [%s] is denial by role of user [%s]", resourcePath, role)
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			dto.NewResponseDto(constant.RespCodeInvalidResourceAccess, constant.RespMsgInvalidResourceAccess, ""))
		return
	}

	c.Next()
}
