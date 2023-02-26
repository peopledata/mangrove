package contracts

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
)

func Client(alchemyApiKey string) (*ethclient.Client, error) {
	//alchemyApiKey := viper.GetString("nft.alchemy_api_key")
	// http://localhost:8545
	return ethclient.Dial(fmt.Sprintf("https://eth-goerli.g.alchemy.com/v2/%s", alchemyApiKey))
}
