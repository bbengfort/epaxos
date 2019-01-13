package epaxos_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEpaxos(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Epaxos Suite")
}
