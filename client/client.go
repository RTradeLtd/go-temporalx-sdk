package client

import (
	"crypto/tls"

	pb "github.com/RTradeLtd/TxPB/v3/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Client is a golang thin-client to
// interact with a TemporalX server
type Client struct {
	conn *grpc.ClientConn
	pb.FileAPIClient
	pb.NodeAPIClient
	pb.StatusAPIClient
	pb.PubSubAPIClient
	pb.AdminAPIClient
	pb.NameSysAPIClient
}

// Opts is used to configure the temporalx grpc client
type Opts struct {
	ListenAddress string
	Insecure      bool
}

// NewClient dials a connection to the grpc server and registered api endpoints
func NewClient(opts Opts, dialOpts ...grpc.DialOption) (*Client, error) {
	if !opts.Insecure {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
	} else {
		dialOpts = append(dialOpts, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(opts.ListenAddress, dialOpts...)
	if err != nil {
		return nil, err
	}
	client := &Client{conn: conn}
	client.FileAPIClient = pb.NewFileAPIClient(client.conn)
	client.NodeAPIClient = pb.NewNodeAPIClient(client.conn)
	client.StatusAPIClient = pb.NewStatusAPIClient(client.conn)
	client.PubSubAPIClient = pb.NewPubSubAPIClient(client.conn)
	client.AdminAPIClient = pb.NewAdminAPIClient(client.conn)
	client.NameSysAPIClient = pb.NewNameSysAPIClient(client.conn)
	return client, nil
}

// Close is used to terminate our grpc client connection
func (c *Client) Close() error {
	return c.conn.Close()
}
