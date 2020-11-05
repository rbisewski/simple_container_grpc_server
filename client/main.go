/*
 * client
 */
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"../pkg/proto/containerd"
	"google.golang.org/grpc"
)

var (
	provision = false
	list      = false
	destroy   = false
)

func init() {
	flag.BoolVar(&provision, "provision", false,
		"Create a new postgres instance.")

	flag.BoolVar(&list, "list", false,
		"List containers.")

	flag.BoolVar(&destroy, "destroy", false,
		"Destroy a container.")
}

func main() {

	flag.Parse()

	if !provision && !list && !destroy {
		flag.Usage()
		os.Exit(0)
	}

	log.Println("Attempting to connect to server ...")

	conn, err := grpc.Dial(":30001", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := containerd.NewContainerdServiceClient(conn)

	cmd := "default"
	if provision {
		cmd = "provision"
	} else if list {
		cmd = "list"
	} else if destroy {
		cmd = "destroy"
	}

	request := &containerd.ContainerdRequest{Command: cmd}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.Containerd(ctx, request)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Response:", response.GetConfirmation())
	os.Exit(0)
}
