/*
 * gRPC server
 */

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"

	"../pkg/proto/containerd"
	"google.golang.org/grpc"
)

const (
	username     = "postgres"
	externalPort = "34352"
	password     = "tdLnXEXC4pCx3RPsBWYENHW04fNeOYMO"
)

type server struct {
	containerd.UnimplementedContainerdServiceServer
}

func main() {
	log.Println("Server running ...")

	lis, err := net.Listen("tcp", ":30001")
	if err != nil {
		log.Fatalln(err)
	}

	srv := grpc.NewServer()
	containerd.RegisterContainerdServiceServer(srv, &server{})

	log.Fatalln(srv.Serve(lis))
}

func (s *server) Containerd(ctx context.Context, request *containerd.ContainerdRequest) (*containerd.ContainerdResponse, error) {

	log.Println(fmt.Sprintf("Request: %s", request.GetCommand()))

	output := ""

	switch cmd := request.GetCommand(); cmd {

	//
	// PROVISION
	//
	case "provision":

		// check if server can communicate a list of container and also check if a postgres container already exists
		stdout, err := exec.Command("/usr/local/bin/k3c", "ps", "-q").Output()
		if err != nil {
			output = err.Error()
			break
		}
		listOfContainerIds := strings.TrimSpace(string(stdout))
		if len(listOfContainerIds) != 0 {
			output = "Error: container already exists."
			break
		}

		// attempt to make a new container, and then get the ID
		stdout, err = exec.Command("/usr/local/bin/k3c",
			"run",
			"-d",
			"--name",
			"example_postgres_container",
			"-p", externalPort+":5432",
			"--env", "POSTGRES_PASSWORD="+password,
			"docker.io/library/postgres:latest").Output()

		if err != nil {
			output = "Error: Container of name 'example_postgres_container' already exists; doing nothing...\n" + err.Error()
			break
		}

		id := strings.TrimSpace(string(stdout))

		// attempt to obtain the fully qualified domain name of the host, else assume localhost
		hostname := ""
		stdout, err = exec.Command("/usr/bin/hostname", "--fqdn").Output()
		if err != nil {
			hostname = "127.0.0.1"
		} else {
			hostname = strings.TrimSpace(string(stdout))
		}

		output += "\n\nProvisioned container details:\n\n"

		output += "ID:       " + id
		output += "\nType:     Postgres server\n"

		output += "\nHostname: " + hostname
		output += "\nPort:     " + externalPort
		output += "\nUser:     " + username
		output += "\nPassword: " + password
		output += "\n"

		output += "\nSample connection command using the psql client:\n"
		output += "\npsql --user " + username + " --host " + hostname + " --port " + externalPort
		output += "\n"
	//
	// LIST
	//
	case "list":
		stdout, err := exec.Command("/usr/local/bin/k3c", "ps", "-a").Output()
		if err != nil {
			output = err.Error()
		} else {
			output = fmt.Sprintf("\n%s", stdout)
		}
	//
	// DESTROY
	//
	case "destroy":
		stdout, err := exec.Command("/usr/local/bin/k3c", "rm", "-f", "example_postgres_container").Output()
		if err != nil {
			output = "Error: Unable to destroy container as it does not exist.\n" + err.Error()
		} else {
			output = fmt.Sprintf("Destroyed the following container:\n%s", stdout)
		}
	default:
		// do nothing
	}

	return &containerd.ContainerdResponse{Confirmation: fmt.Sprintf("%s", output)}, nil
}
