MicroService Foundation
=======================

> **ALPHA:** Lile is currently considered "Alpha" in that things may change. Currently I am gathering feedback and will finalise Lile shortly to avoid breaking changes going forward.

This package provides a foundation to build [gRPC](google.golang.org/grpc) based microservices.

It takes care of:

- Server initialization
- Logging
- Metrics
- Middleware

# Dependencies

You should vendor the dependencies needed in your project. Depending on the features used the
requirements will vary.

- [Go >= v1.8](https://golang.org/dl)
- [dep](https://github.com/golang/dep)
- [gRPC](google.golang.org/grpc)
- [Go Protocol Buffers](github.com/golang/protobuf)
- [cmux](github.com/cockroachdb/cmux)
- [metrics](github.com/codahale/metrics)
- [getconf](github.com/jllopis/getconf)
- [libkv](github.com/abronan/libkv)
- [consul api](github.com/hashicorp/consul/api)
- [etcd v3](github.com/coreos/etcd/clientv3)
- [Trace](golang.org/x/net/trace)

Also, the protobuf compiler will be needed if you want to rebuild the go files:

- [Compilador Protocol Buffers (protoc) v3.5.0](https://github.com/google/protobuf/releases)

# Guide

- [Installation](#installation)
- [Service Definition](#service-definition)
- [RPC Methods](#rpc-methods)

## Installation

## Service Definition

## RPC Methods
