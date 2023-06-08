package relayer

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	omniClient "github.com/synapsecns/sanguine/services/omnirpc/client"
)

// HandleSendRequest wraps handleSendRequest for testing.
func (c CCTPRelayer) HandleSendRequest(parentCtx context.Context, txhash common.Hash, originChain uint32) (err error) {
	return c.handleSendRequest(parentCtx, txhash, originChain)
}

// FetchAttestation wraps fetchAttestation for testing.
func (c CCTPRelayer) FetchAttestation(parentCtx context.Context, chainID uint32, msg *UsdcMessage) {
	c.fetchAttestation(parentCtx, chainID, msg)
}

// SetOmnirpcClient sets the omnirpc client for testing.
func (c *CCTPRelayer) SetOmnirpcClient(client omniClient.RPCClient) {
	c.omnirpcClient = client
}

// RecvUsdcMsg receives a usdc message from the given chain.
func (c *CCTPRelayer) GetUsdcMsgRecvChan(chainID uint32) chan *UsdcMessage {
	return c.chainRelayers[chainID].usdcMsgRecvChan
}

// SendUsdcMsg receives a usdc message from the given chain.
func (c *CCTPRelayer) GetUsdcMsgSendChan(chainID uint32) chan *UsdcMessage {
	return c.chainRelayers[chainID].usdcMsgSendChan
}
