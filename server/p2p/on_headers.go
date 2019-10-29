package p2p

import (
	"github.com/daglabs/btcd/peer"
	"github.com/daglabs/btcd/wire"
)

// OnHeaders is invoked when a peer receives a headers bitcoin
// message.  The message is passed down to the sync manager.
func (sp *Peer) OnHeaders(_ *peer.Peer, msg *wire.MsgHeaders) {
	sp.server.SyncManager.QueueHeaders(msg, sp.Peer)
}