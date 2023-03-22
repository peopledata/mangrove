package middleware

import (
	controller2 "mangrove/internal/controller"
	"mangrove/pkg/contracts"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

func EthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		network := viper.GetString("nft.network")
		apiKey := viper.GetString("nft.infura_api_key")
		client, err := contracts.Client(network, apiKey)
		if err != nil {
			controller2.ResponseErr(c, controller2.CodeInvalidEthClient)
			c.Abort()
			return
		}
		c.Set(controller2.CtxEthKey, client)
		c.Next()
	}
}
