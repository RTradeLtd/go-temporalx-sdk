package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"time"

	pb "github.com/RTradeLtd/TxPB/v3/go"
	"github.com/schollz/progressbar/v2"
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

// UploadFile is a convenience function for uploading a single file
func (c *Client) UploadFile(
	ctx context.Context,
	file io.Reader,
	fileSize int64,
	opts *pb.UploadOptions,
	printProgress bool,
	grpcOpts ...grpc.CallOption,
) (*pb.PutResponse, error) {
	stream, err := c.FileAPIClient.UploadFile(ctx, grpcOpts...)
	if err != nil {
		return nil, err
	}

	// declare file options
	if err := stream.Send(&pb.UploadRequest{Options: opts}); err != nil {
		return nil, err
	}
	// upload file - chunked at 5mb each
	buf := make([]byte, 4194294)
	var pt *progressTracker
	if printProgress {
		pt = newPT(fileSize)
	}
	for {
		n, err := file.Read(buf)
		if err != nil && err == io.EOF {
			// only break if we haven't read any bytes, otherwise exit
			if n == 0 {
				break
			}
		} else if err != nil && err != io.EOF {
			return nil, err
		}
		if err := stream.Send(&pb.UploadRequest{Blob: &pb.Blob{Content: buf[:n]}}); err != nil {
			return nil, err
		}
		if printProgress {
			pt.Update(n)
		}
	}
	if printProgress {
		fmt.Println("")
	}
	// done
	return stream.CloseAndRecv()
}

// DownloadFile is a convenience function for downloading a single file
func (c *Client) DownloadFile(
	ctx context.Context,
	download *pb.DownloadRequest,
	printProgress bool,
	grpcOpts ...grpc.CallOption,
) (*bytes.Buffer, error) {
	stream, err := c.FileAPIClient.DownloadFile(ctx, download, grpcOpts...)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	var pt *progressTracker
	if printProgress {
		pt = newPT(int64(buf.Len()))
	}
	lastPrint := time.Time{}
	for {
		b, err := stream.Recv()
		if err != nil && err == io.EOF {
			// only exit at end of stream if we have nothing left to process
			if b == nil {
				break
			}
		} else if err != nil && err != io.EOF {
			return nil, err
		}
		if download.RangeStart != 0 && b.GetBlob().GetRangeStart() == 0 {
			// we don't need backwards support here since we should update all the servers to support range download.
			return nil, errors.New("range download is not supported on the connected server")
		}
		count := len(b.Blob.Content)
		if _, err := buf.Write(b.GetBlob().GetContent()); err != nil {
			return nil, err
		}
		if printProgress {
			now := time.Now()
			if now.After(lastPrint.Add(time.Second / 100)) {
				lastPrint = now
				if pt.bar.GetMax() < buf.Len() {
					pt.bar.ChangeMax(buf.Len())
					fmt.Println("changed max")
				}
				fmt.Println(count)
				pt.Update(count)
			}
		}
	}
	if err := stream.CloseSend(); err != nil {
		return nil, err
	}
	return buf, nil
}

// Close is used to terminate our grpc client connection
func (c *Client) Close() error {
	return c.conn.Close()
}

type progressTracker struct {
	bar *progressbar.ProgressBar
}

func newPT(maxBytes int64) *progressTracker {
	return &progressTracker{
		bar: progressbar.NewOptions64(
			maxBytes,
			progressbar.OptionSetRenderBlankState(true),
		),
	}
}

func (p *progressTracker) Update(sent int) {
	p.bar.Add(sent)
}
