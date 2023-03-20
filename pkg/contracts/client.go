package contracts

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
)

func Client(network, apiKey string) (*ethclient.Client, error) {
	//alchemyApiKey := viper.GetString("nft.alchemy_api_key")
	//return ethclient.Dial("http://localhost:8545")
	//https://mainnet.infura.io/v3/9b648f95861e4d719432d38d1aa5d05d
	//https://goerli.infura.io/v3/9b648f95861e4d719432d38d1aa5d05d
	//#  wss://goerli.infura.io/ws/v3/a62e439c8c1048b6a1f983e5d9a0e72d
	return ethclient.Dial(fmt.Sprintf("wss://%s.infura.io/ws/v3/%s", network, apiKey))
	//return ethclient.Dial(fmt.Sprintf("wss://eth-goerli.g.alchemy.com/v2/%s", alchemyApiKey))
}
