package epaxos_test

import (
	"encoding/json"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/bbengfort/epaxos"
	"github.com/bbengfort/epaxos/pb"
)

var _ = Describe("Log", func() {

	var logs *Logs
	var config *Config

	BeforeSuite(func() {
		data, err := ioutil.ReadFile("testdata/config.json")
		Ω(err).ShouldNot(HaveOccurred())

		err = json.Unmarshal(data, &config)
		Ω(err).ShouldNot(HaveOccurred())

		Ω(config.Validate()).Should(Succeed())

		// Set the quroum size to 5
		config.Peers = config.Peers[:5]
	})

	Describe("from empty logs", func() {

		BeforeEach(func() {
			logs = NewLog(config)
		})

		It("should be able to create an instance for a replica", func() {
			inst, err := logs.Create(2, []*pb.Operation{{Type: pb.AccessType_WRITE, Key: "foo", Value: []byte("bar")}})
			Ω(err).ShouldNot(HaveOccurred())

			slot, err := logs.LastApplied(2)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(slot).Should(Equal(inst.Slot))
		})

		It("should not create an instance for a replica not in the quorum", func() {
			inst, err := logs.Create(32, []*pb.Operation{{Type: pb.AccessType_WRITE, Key: "foo", Value: []byte("bar")}})
			Ω(inst).Should(BeNil())
			Ω(err).Should(MatchError("no log for replica with PID 32"))
		})

		It("should be able to insert an instance for a replica", func() {
			inst := &pb.Instance{
				Replica: 4,
				Slot:    0,
				Seq:     1,
				Deps:    make(map[uint32]uint64),
				Status:  pb.Status_INITIAL,
				Acks:    1,
				Ops:     []*pb.Operation{{Type: pb.AccessType_WRITE, Key: "foo", Value: []byte("bar")}},
			}

			Ω(logs.Insert(inst)).Should(Succeed())
			slot, err := logs.LastApplied(4)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(slot).Should(Equal(inst.Slot))
		})

		It("should not insert an instance for a replica not in the quorum", func() {
			inst := &pb.Instance{
				Replica: 48,
				Slot:    0,
				Seq:     1,
				Deps:    make(map[uint32]uint64),
				Status:  pb.Status_INITIAL,
				Acks:    1,
				Ops:     []*pb.Operation{{Type: pb.AccessType_WRITE, Key: "foo", Value: []byte("bar")}},
			}

			Ω(logs.Insert(inst)).Should(MatchError("no log for replica with PID 48"))
		})

		It("should not insert an instance ahead of the current slot", func() {
			inst := &pb.Instance{
				Replica: 4,
				Slot:    10,
				Seq:     1,
				Deps:    make(map[uint32]uint64),
				Status:  pb.Status_INITIAL,
				Acks:    1,
				Ops:     []*pb.Operation{{Type: pb.AccessType_WRITE, Key: "foo", Value: []byte("bar")}},
			}

			Ω(logs.Insert(inst)).Should(MatchError("cannot insert into slot 10 when next slot is 0"))
		})

		It("should not insert an instance twice", func() {
			inst := &pb.Instance{
				Replica: 4,
				Slot:    0,
				Seq:     1,
				Deps:    make(map[uint32]uint64),
				Status:  pb.Status_INITIAL,
				Acks:    1,
				Ops:     []*pb.Operation{{Type: pb.AccessType_WRITE, Key: "foo", Value: []byte("bar")}},
			}

			Ω(logs.Insert(inst)).Should(Succeed())
			Ω(logs.Insert(inst)).Should(MatchError("there is already an instance in slot 0"))
		})

	})

})
