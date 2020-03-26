package client

import (
	"context"
	"fmt"
	"io"
	"time"

	pb "github.com/RTradeLtd/TxPB/v3/go"
	au "github.com/logrusorgru/aurora"
)

var (
	// PrintMsgExecFunc is used to print pubsub messages as they are received
	PrintMsgExecFunc = func(msg *pb.PubSubResponse) {
		for _, msg := range msg.GetMessage() {
			fmt.Printf(
				"%s\n",
				au.Bold(au.Green(string(msg.GetData()))),
			)
		}
	}
)

// PSSubscribe is a helper method to subscribe to a given topic
// It takes in a msgExecFunc argument that allows you to craft your own functions
// that do different executions whenever a message is received.
// If msgExecFunc is empty then we use a default function that prints received messages
func (c *Client) PSSubscribe(ctx context.Context, topic string, discover bool, msgExecFunc func(msg *pb.PubSubResponse)) error {
	if msgExecFunc == nil {
		msgExecFunc = PrintMsgExecFunc
	}
	stream, err := c.PubSub(ctx)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	errChan := make(chan error, 1)
	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				errChan <- err
				return
			}
			msgExecFunc(msg)
		}
	}()
	if err := stream.Send(&pb.PubSubRequest{
		RequestType: pb.PSREQTYPE_PS_SUBSCRIBE,
		Topics:      []string{topic},
	}); err != nil {
		return err
	}
	for {
		select {
		case err := <-errChan:
			if err == io.EOF {
				return nil
			}
		case <-stream.Context().Done():
			return nil
		case <-ctx.Done():
			return stream.CloseSend()
		}
	}
}

// PSPublish is a helper function used to publish pubusb messages
func (c *Client) PSPublish(ctx context.Context, topic string, data []byte) error {
	stream, err := c.PubSub(ctx)
	if err != nil {
		return err
	}
	if err := stream.Send(&pb.PubSubRequest{
		RequestType: pb.PSREQTYPE_PS_PUBLISH,
		Topics:      []string{topic},
		Data:        data,
	}); err != nil {
		return err
	}
	// TODO(bonedaddy): redesign pubsub to send back
	// a confirmation after a publish to work around this issue
	time.Sleep(time.Second * 1)
	return stream.CloseSend()
}

// PSListPeers is a helper function used to list pubsub peers
func (c *Client) PSListPeers(ctx context.Context, topics ...string) ([]*pb.PubSubPeer, error) {
	stream, err := c.PubSub(ctx)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(&pb.PubSubRequest{
		RequestType: pb.PSREQTYPE_PS_LIST_PEERS,
		Topics:      topics,
	}); err != nil {
		return nil, err
	}
	msg, err := stream.Recv()
	if err != nil {
		return nil, err
	}
	return msg.GetPeers(), nil
}

// PSGetTopics is a helper function used to fetch subscribed pubsub topics
func (c *Client) PSGetTopics(ctx context.Context) ([]string, error) {
	stream, err := c.PubSub(ctx)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(&pb.PubSubRequest{
		RequestType: pb.PSREQTYPE_PS_GET_TOPICS,
	}); err != nil {
		return nil, err
	}
	msg, err := stream.Recv()
	if err != nil {
		return nil, err
	}
	return msg.GetTopics(), nil
}
