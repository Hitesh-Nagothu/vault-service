package ipfs

import (
	"bytes"
	"fmt"
	"io"

	shell "github.com/ipfs/go-ipfs-api"
)

type IPFS struct {
	api *shell.Shell
}

var ipfsInstance *IPFS

func GetIPFSInstance() *IPFS {
	if ipfsInstance == nil {
		api := shell.NewShell("/ip4/127.0.0.1/tcp/5001")
		ipfsInstance = &IPFS{
			api: api,
		}
	}
	return ipfsInstance
}

// AddContent adds content to the IPFS network and returns the CID.
func (ipfs *IPFS) AddContent(content []byte) (string, error) {
	cid, err := ipfs.api.Add(bytes.NewReader(content))
	if err != nil {
		return "", fmt.Errorf("failed to add content to IPFS: %w", err)
	}
	return cid, nil
}

// GetContent retrieves content from the IPFS network using the given CID.
func (ipfs *IPFS) GetContent(cid string) ([]byte, error) {
	reader, err := ipfs.api.Cat(cid)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve content from IPFS: %w", err)
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read content from IPFS: %w", err)
	}
	fmt.Println(content)
	return content, nil
}
