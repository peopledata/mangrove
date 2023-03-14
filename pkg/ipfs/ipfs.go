package ipfs

import (
	"io/ioutil"

	ipfsapi "github.com/ipfs/go-ipfs-api"
)

func Read(shell *ipfsapi.Shell, ipfsPath string) ([]byte, error) {
	// Fetch the IPFS data from the URL
	ipfsData, err := shell.Cat(ipfsPath)
	if err != nil {
		return nil, err
	}
	// Convert the byte stream to a string
	return ioutil.ReadAll(ipfsData)
}
