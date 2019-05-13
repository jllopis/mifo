package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"bitbucket.org/acbapis/acbapis/lib/go/common"
	"bitbucket.org/acbapis/acbapis/lib/go/status"
	"google.golang.org/grpc"
)

var (
	serverPort = flag.String("port", "58000", "remote gRPC server port")

	remoteConn *grpc.ClientConn
)

func main() {
	flag.Parse()
	dial()
	defer remoteConn.Close()

	client := status.NewStatusServiceClient(remoteConn)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	serverTime, err := client.GetServerTime(ctx, &common.EmptyMessage{})
	if err != nil {
		log.Printf("error getting server time: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Time server: %s\n", time.Unix(serverTime.GetValue()/1e9, 0))

	serverVersion, err := client.GetVersion(ctx, &common.EmptyMessage{})
	if err != nil {
		log.Printf("error getting server version: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Server version: %s\n", serverVersion)

}

// dialCompeticion opens a connection with Competicion service
func dial() error {
	var err error
	host := "localhost:" + *serverPort
	remoteConn, err = grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Error dialing %s: %s", host, err)
		return err
	}

	return nil
}
