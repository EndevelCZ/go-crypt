package grpc

import (
	"context"
	"fmt"
	"io"

	"github.com/EndevelCZ/go-crypt/gcspb"
	"google.golang.org/grpc"
)

type GcsClientGRPC struct {
	conn      *grpc.ClientConn
	client    gcspb.GuploadServiceClient
	chunkSize int
}

func NewGcsClientGRPC(addr string) (*GcsClientGRPC, error) {
	var (
		grpcOpts = []grpc.DialOption{}
		// grpcCreds credentials.TransportCredentials
	)
	if addr == "" {
		return nil, fmt.Errorf("address must be specified")
	}
	grpcOpts = append(grpcOpts, grpc.WithInsecure())
	conn, err := grpc.Dial(addr, grpcOpts...)
	if err != nil {
		return nil, err
	}
	client := gcspb.NewGuploadServiceClient(conn)
	return &GcsClientGRPC{
		chunkSize: 1 << 10,
		conn:      conn,
		client:    client,
	}, nil
}

func (c *GcsClientGRPC) Upload(ctx context.Context, r io.Reader) error {

	stream, err := c.client.Upload(ctx)
	if err != nil {
		return err
	}
	defer stream.CloseSend()
	buf := make([]byte, c.chunkSize)
	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		err = stream.Send(&gcspb.Chunk{
			Content: buf[:n],
		})
		if err != nil {
			return err
		}
	}
	status, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	if status.Code != gcspb.UploadStatusCode_Ok {
		return fmt.Errorf("upload failed - msg: %s", status.Message)
	}
	return nil
}
