/*
Package epaxos implements the ePaxos consensus algorithm.
*/
package epaxos

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/bbengfort/x/noplog"
	"google.golang.org/grpc/grpclog"
)

//===========================================================================
// Package Initialization
//===========================================================================

// PackageVersion of the current ePaxos implementation
const PackageVersion = "0.1"

// Initialize the package and random numbers, etc.
func init() {
	// Set the random seed to something different each time.
	rand.Seed(time.Now().UnixNano())

	// Initialize our debug logging with our prefix
	SetLogger(log.New(os.Stdout, "[epaxos] ", log.Lmicroseconds))
	cautionCounter = new(counter)
	cautionCounter.init()

	// Stop the grpc verbose logging
	grpclog.SetLogger(noplog.New())
}
