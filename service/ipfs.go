package service

import (
	"bytes"
	"fmt"
	"io"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type IPFSService struct {
	api    *shell.Shell
	logger *zap.Logger
}

var ipfsInstance *IPFSService

func NewIPFSService(logger *zap.Logger) *IPFSService {
	ipfsURL := viper.GetString("ipfs.url")
	ipfsPort := viper.GetString("ipfs.port")
	if ipfsInstance == nil {
		api := shell.NewShell(ipfsURL + ipfsPort)
		ipfsInstance = &IPFSService{
			api:    api,
			logger: logger,
		}
	}
	return ipfsInstance
}

func (ipfs *IPFSService) GetIPFSInstance() *IPFSService {
	return ipfsInstance
}

// AddContent adds content to the IPFS network and returns the CID.
func (ipfs *IPFSService) AddContent(content []byte) (string, error) {
	cid, err := ipfs.api.Add(bytes.NewReader(content))
	if err != nil {
		return "", fmt.Errorf("failed to add content to IPFS: %w", err)
	}
	return cid, nil
}

// GetContent retrieves content from the IPFS network using the given CID.
func (ipfs *IPFSService) GetContent(cid string) ([]byte, error) {
	reader, err := ipfs.api.Cat(cid)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve content from IPFS: %w", err)
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read content from IPFS: %w", err)
	}
	// fmt.Println(content)
	return content, nil
}
