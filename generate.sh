#!/bin/bash

protoc --go_out=. --go-grpc_out=. greet/greetpb/greet.proto