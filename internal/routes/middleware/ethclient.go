package middleware

import (
	controller2 "patronus/internal/controller"
	"patronus/pkg/contracts"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

func EthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		alchemyApiKey := viper.GetString("nft.alchemy_api_key")
		// http://localhost:8545
		client, err := contracts.Client(alchemyApiKey)
		if err != nil {
			controller2.ResponseErr(c, controller2.CodeInvalidEthClient)
			c.Abort()
			return
		}
		c.Set(controller2.CtxEthKey, client)
		c.Next()
	}
}
