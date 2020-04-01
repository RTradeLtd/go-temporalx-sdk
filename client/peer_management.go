package client

import (
	"context"
	"errors"

	pb "github.com/RTradeLtd/TxPB/v3/go"
)

// Connected is a convenience function for checking if the given peer is connected to our node
func (c *Client) Connected(ctx context.Context, peer string) (bool, error) {
	resp, err := c.ConnMgmt(ctx, &pb.ConnMgmtRequest{
		RequestType: pb.CONNMGMTREQTYPE_CM_STATUS,
		PeerIDs:     []string{peer},
	})
	if err != nil {
		return false, err
	}
	return resp.GetConnected()[peer], nil
}

// ConnectToPeer is a convenience function to connect to a single peer
func (c *Client) ConnectToPeer(ctx context.Context, multiaddr string) error {
	_, err := c.ConnMgmt(ctx, &pb.ConnMgmtRequest{
		RequestType: pb.CONNMGMTREQTYPE_CM_CONNECT,
		MultiAddrs:  []string{multiaddr},
	})
	return err
}

// DisconnectFromPeer is a convenience function to disconnect from a single peer
func (c *Client) DisconnectFromPeer(ctx context.Context, id string) error {
	resp, err := c.ConnMgmt(ctx, &pb.ConnMgmtRequest{
		RequestType: pb.CONNMGMTREQTYPE_CM_DISCONNECT,
		PeerIDs:     []string{id},
	})
	if err != nil {
		return err
	}
	if !resp.GetStatus()[id].GetDisconnected() {
		return errors.New(resp.GetStatus()[id].GetReason())
	}
	return nil
}

// GetPeerCount is a convenience function for returning the number of peers we are connected to
func (c *Client) GetPeerCount(ctx context.Context) (int, error) {
	resp, err := c.ConnMgmt(ctx, &pb.ConnMgmtRequest{
		RequestType: pb.CONNMGMTREQTYPE_CM_GET_PEERS,
	})
	if err != nil {
		return 0, err
	}
	return len(resp.GetPeerIDs()), nil
}
