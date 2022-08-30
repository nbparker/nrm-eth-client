# nrm-eth-client

Send protos to server for storage

## Specification

Currently reads a folder of json updates and sends to the server.

## Customisation

Can be customised to send varied types of protobuf update to the server.

How to:

- Replace `GenericUpdate` protobuf in the proto file with resource of choice (do not rename)
- Generate new protos and replace existing in [proto folder](proto/nrm) for both client AND server

## Usage

- Build (see [Dockerfile](Dockerfile) for example)
- Run with parameters:

  - Read JSON from folder:

    `go run cmd/client/main.go --addr="nrm-eth-server:50051" --updates_folder="updates"`
