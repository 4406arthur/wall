package middleware

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"wall/pkg/developer"
	"wall/utils/logger"

	"github.com/gin-gonic/gin"
)

func BasicAuth(log logger.Logger, r developer.Repository) gin.HandlerFunc {

	return func(c *gin.Context) {
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
		logger.Log(log, "Debug", "NA", logger.BuildLogInfo(c), toString(auth))

		if len(auth) != 2 || auth[0] != "Basic" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)
		username := pair[0]
		plainPassword := pair[1]

		// var dev *entity.Developer
		dev, err := r.Find(username)
		if err != nil {
			logger.Log(log, "Error", "NA", logger.Trace(), err.Error())
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if len(pair) != 2 || !authenticateUser(plainPassword, dev.Password) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("devID", dev.Name)
		c.Set("scopes", dev.Grants)
		c.Next()
	}
}

func authenticateUser(plainPassword, ans string) bool {

	afterHash := hex.EncodeToString(newSHA256([]byte(plainPassword)))

	return afterHash == ans
}

func newSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func toString(data interface{}) string {
	str := fmt.Sprintf("%#v", data)
	return str
}
