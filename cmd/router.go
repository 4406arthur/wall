package core

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"
	apis "wall/cmd/api"
	"wall/pkg/developer"
	"wall/pkg/middleware"
	"wall/pkg/token"
	"wall/utils/logger"
	"wall/utils/throttle"

	ginlogrus "github.com/4406arthur/gin-logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	ginprometheus "github.com/zsais/go-gin-prometheus"

	queue "wall/pkg/queue/mocks"
	"wall/pkg/token/mocks"
)

//InitRouter used config api endpoint and auth middleware
func InitRouter(log logger.Logger, config *viper.Viper) {

	r := gin.Default()
	//set up prometheus exporter
	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)

	//set up logger
	r.Use(ginlogrus.Logger(log.GetLogger()))

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	// developerRepo := developer.NewMongoRepository(config.GetString("mongo_config.endpoint"))
	developerRepo := developer.NewConfigmapRepository(config.GetString("developer_repo.file_path"))
	//tokenRepo := token.NewMongoRepository(config.GetString("mongodb_config.conn_string"))
	tokenRepo := &mocks.Repository{}
	// mQueue := queue.NewQueue(log, config.GetString("mq_config.endpoint"), config.GetString("mq_config.queue_name"))
	mockQueue := &queue.Service{}

	r.POST("/api/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	devGroup := r.Group("/developer")
	//Token bucket: 20 tickets withun 10 sec
	devGroup.Use(throttle.Throttle(10, 20))
	devGroup.Use(RequestLogger(log))
	devGroup.Use(middleware.BasicAuth(log, developerRepo))
	{
		signBytes, _ := ioutil.ReadFile(config.GetString("server_config.jwt"))
		signKey, _ := jwt.ParseRSAPrivateKeyFromPEM(signBytes)

		// tokenService := token.NewTokenService(tokenRepo, developerRepo, signKey, mQueue)
		tokenService := token.NewTokenService(tokenRepo, signKey, mockQueue)
		TokenHandler := apis.TokenHandlerInit(log, tokenService)
		devGroup.POST("/getToken", TokenHandler.IssueToken)
		devGroup.POST("/revokeToken", TokenHandler.RevokeToken)
	}

	if config.IsSet("server_config.cert") && config.IsSet("server_config.key") {
		r.RunTLS(
			config.GetString("server_config.listen_addr"),
			config.GetString("server_config.cert"),
			config.GetString("server_config.key"),
		)
	}

	s := &http.Server{
		Addr:         config.GetString("server_config.listen_addr"),
		Handler:      r,
		ReadTimeout:  4 * time.Second,
		WriteTimeout: 4 * time.Second,
		// 1Mb
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}

func RequestLogger(esLog logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, _ := ioutil.ReadAll(c.Request.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.

		logger.Log(esLog, "Info", "RQ", logger.BuildLogInfo(c), readBody(rdr1))

		c.Request.Body = rdr2
		c.Next()
	}
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}
