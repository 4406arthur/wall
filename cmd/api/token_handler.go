package apis

import (
	"net/http"
	"wall/pkg/token"
	"wall/utils/logger"

	"github.com/gin-gonic/gin"
)

// Token handler
type TokenHandler struct {
	Log      logger.Logger
	TService token.Service
}

type getTokenRQ struct {
	UserID string `json:"userID"`
}

type getTokenRP struct {
	Token string `json:"token"`
}

type revokeTokenRQ struct {
	Tid string `json:"tid"`
}

func TokenHandlerInit(log logger.Logger, s token.Service) *TokenHandler {
	return &TokenHandler{
		Log:      log,
		TService: s,
	}
}

func (tokenHandler *TokenHandler) IssueToken(c *gin.Context) {
	var rq *getTokenRQ
	err := c.BindJSON(&rq)
	if err != nil {
		logger.Log(tokenHandler.Log, "Error", "NA", logger.BuildLogInfo(c), err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	devID, _ := c.Get("devID")
	scopes, _ := c.Get("scopes")

	token, err := tokenHandler.TService.Issue(rq.UserID, devID.(string), scopes.([]string))
	if err != nil {
		logger.Log(tokenHandler.Log, "Error", "NA", logger.BuildLogInfo(c), err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &getTokenRP{
		Token: token,
	})
}

func (tokenHandler *TokenHandler) RevokeToken(c *gin.Context) {
	var rq *revokeTokenRQ
	err := c.BindJSON(&rq)
	if err != nil {
		logger.Log(tokenHandler.Log, "Error", "NA", logger.BuildLogInfo(c), err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	devID, _ := c.Get("devID")
	err = tokenHandler.TService.Revoke(rq.Tid, devID.(string))
	if err != nil {
		logger.Log(tokenHandler.Log, "Error", "NA", logger.BuildLogInfo(c), err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}
