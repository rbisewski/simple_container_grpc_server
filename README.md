# Simple container gRPC server

The repo contains code for a simple gRPC server that manages containers, which
is written in golang.

The most simple and the quickest way to play around with containerd is to use
Rancher's k3c utility, which offers all of the raw speed of containerd with
the featureset of docker, so that is the container engine this app uses.

Consider reading the *Requirements* at a minimum to get the hang of how
to setup this job scheduler locally, although the author makes no promises
as to whether this will run on your specific Linux system.

# Requirements

The program itself was designed around a standard Debian Linux environment,
with the following requirements:

* k3c
* golang 1.10+
* libcontainer
* grpc libraries
* protobuf library includes

In the event that this program does not appear to work on a particular
non-mainstream distro, feel free to contact me if you need assistance
and I will make note of it in future versions of this readme.

That said, an easy Makefile target has been provided to prepare a system:

```bash
make prep
```

The above command will install k3c and go get the dependencies you need,
and then also install protoc dependencies. Consider examining the Makefile
if you are unable to get your system prepared, or contact the author for
additional assistence.

# Installation

Having prepared your system environment as per the above section, you can
build both the server binary and the client binary, along with the protobufs
like so:

```bash
make
```

# Basic Usage Instructions

This section assumes you have prepared your system and then built the server
and client as per the prior sections above.

Start a background instance of k3c like so:

```bash
make run_daemon
```

Once the k3c daemon has started, then start up an instance of the gRPC server
like below. Note that the server runs on 0.0.0.0:30001, so make sure you have
no other service using that port.

```bash
make run_server
```

Afterwards, you can run the client... 

To print a list of containers:

```bash
./bin/client -list
```

To provision a new postgres container:

```bash
./bin/client -provision
```

To destroy the postgres container:

```bash
./bin/client -destroy
```

# Author

This was created by Robert Bisewski at Ibis Cybernetics. For more
information, please contact us at:

* Website: www.ibiscybernetics.com

* Email: contact@ibiscybernetics.com
