#!/bin/bash
export EPAXOS="../../cmd/epaxos/main.go"
export SERVE="go run $EPAXOS -c config.json serve"
export COMMIT="go run $EPAXOS -c config.json commit -k foo -v $(ts)"
export BENCH="go run $EPAXOS -c config.json bench -r 5000 -b"
