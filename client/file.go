package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	pb "github.com/RTradeLtd/TxPB/v3/go"
	"google.golang.org/grpc"
)

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
