package grpc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"strconv"

	"cloud.google.com/go/storage"
	"github.com/EndevelCZ/go-crypt/gcspb"
	"github.com/sirupsen/logrus"
	"github.com/xtgo/uuid"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

// https://github.com/grpc/grpc-go/tree/master/examples/route_guide
// https://grpc.io/docs/tutorials/basic/go.html#example-code-and-setup
// https://github.com/cirocosta/gupload
type GcsServerGRPC struct {
	server     *grpc.Server
	port       int
	gcsClient  *storage.Client
	bucketName string
	ctx        context.Context
}

func NewGcsServerGRPC(port int, bucketName, gcsServiceAccountPath string) (*GcsServerGRPC, error) {
	if port == 0 {
		return nil, fmt.Errorf("Port must be specified")
	}
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(gcsServiceAccountPath))
	if err != nil {
		return nil, err
	}
	// gcsClient := gcp.NewClient(stiface.AdaptClient(client), bucketName, ctx)
	return &GcsServerGRPC{
		server:     nil,
		port:       port,
		bucketName: bucketName,
		ctx:        ctx,
		gcsClient:  client,
	}, nil
}
func (g *GcsServerGRPC) Listen() (err error) {
	var (
		listener net.Listener
		grpcOpts = []grpc.ServerOption{}
		// grpcCreds credentials.TransportCredentials
	)

	listener, err = net.Listen("tcp", ":"+strconv.Itoa(g.port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d", g.port)

	}

	// if s.certificate != "" && s.key != "" {
	// 	grpcCreds, err = credentials.NewServerTLSFromFile(
	// 		s.certificate, s.key)
	// 	if err != nil {
	// 		err = errors.Wrapf(err,
	// 			"failed to create tls grpc server using cert %s and key %s",
	// 			s.certificate, s.key)
	// 		return
	// 	}

	// 	grpcOpts = append(grpcOpts, grpc.Creds(grpcCreds))
	// }

	g.server = grpc.NewServer(grpcOpts...)
	gcspb.RegisterGuploadServiceServer(g.server, g)

	err = g.server.Serve(listener)
	if err != nil {
		err = fmt.Errorf("errored listening for grpc connections")
		return
	}

	return
}

// gcspb.GuploadService_UploadServer is interface
//type GuploadService_UploadServer interface {
//	SendAndClose(*UploadStatus) error
//	Recv() (*Chunk, error)
//	grpc.ServerStream
//}
func (g *GcsServerGRPC) Upload(stream gcspb.GuploadService_UploadServer) error {
	chunk := &gcspb.Chunk{}
	objectName := uuid.NewRandom()
	// if err != nil {
	// 	return err
	// }
	// r := bytes.NewReader(chunk.GetContent())
	// g.gcsClient.UploadObject(r ,objectName)
	wc := g.gcsClient.Bucket(g.bucketName).Object(objectName.String()).NewWriter(g.ctx)
	var err error
	for {
		chunk, err = stream.Recv()
		logrus.Info("chunk:", chunk)
		if err != nil {
			if err == io.EOF {
				logrus.Info("upload received")
				// fmt.Println(chunk)
				if err := wc.Close(); err != nil {
					return err
				}
				break
			}
			return err
		}
		r := bytes.NewReader(chunk.GetContent())
		if _, err := io.Copy(wc, r); err != nil {
			return err
		}
	}

	err = stream.SendAndClose(&gcspb.UploadStatus{
		Message: "Upload received with success",
		Code:    gcspb.UploadStatusCode_Ok,
	})
	if err != nil {
		return err
	}
	return nil
}
