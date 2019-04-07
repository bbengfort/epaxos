package epaxos_test

import (
	"encoding/json"
	"io/ioutil"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/bbengfort/epaxos"
)

var _ = Describe("Config", func() {

	It("should validate a correct configuration", func() {
		conf := &Config{
			Name:      "foo",
			Seed:      42,
			Timeout:   "300ms",
			Aggregate: true,
			Thrifty:   true,
			LogLevel:  2,
			Uptime:    "15m",
			Metrics:   "metrics.json",
		}
		Ω(conf.Validate()).Should(Succeed())
	})

	It("should be valid with loaded defaults", func() {
		conf := new(Config)

		confPath, err := conf.GetPath()
		Ω(confPath).Should(BeZero())
		Ω(err).Should(HaveOccurred())

		Ω(conf.Load()).Should(Succeed())

		// Validate configuration defaults
		Ω(conf.Timeout).Should(Equal("500ms"))
		Ω(conf.Aggregate).Should(BeFalse())
		Ω(conf.Thrifty).Should(BeFalse())
		Ω(conf.LogLevel).Should(Equal(3))

		// Validate non configurations
		Ω(conf.Name).Should(BeZero())
		Ω(conf.Seed).Should(BeZero())
		Ω(conf.Peers).Should(BeZero())
		Ω(conf.Uptime).Should(BeZero())
		Ω(conf.Metrics).Should(BeZero())
	})

	It("should be able to parse durations", func() {
		conf := &Config{Timeout: "10s", Uptime: "10s"}

		duration, err := conf.GetTimeout()
		Ω(err).ShouldNot(HaveOccurred())
		Ω(duration).Should(Equal(10 * time.Second))

		duration, err = conf.GetUptime()
		Ω(err).ShouldNot(HaveOccurred())
		Ω(duration).Should(Equal(10 * time.Second))
	})

	Describe("quorum and peers configuration", func() {

		var config *Config

		BeforeEach(func() {
			data, err := ioutil.ReadFile("testdata/config.json")
			Ω(err).ShouldNot(HaveOccurred())

			err = json.Unmarshal(data, &config)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(config.Validate()).Should(Succeed())
		})

		It("should be able to get it's local peer config", func() {
			Ω(config.GetName()).Should(Equal("bravo"))

			peer, err := config.GetPeer()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(peer.PID).Should(Equal(uint32(2)))
		})

		Describe("size 3 quorum", func() {

			BeforeEach(func() {
				config.Peers = config.Peers[:3]
				Ω(config.Peers).Should(HaveLen(3))
			})

			It("should get 2 remotes to connect to", func() {
				remotes, err := config.GetRemotes()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(remotes).Should(HaveLen(2))

				for _, remote := range remotes {
					Ω(remote.PID).ShouldNot(Equal(uint32(2)))
				}

			})

			It("should return nil if not thrifty", func() {
				config.Thrifty = false
				Ω(config.GetThrifty()).Should(BeNil())
			})

			It("should return 1 thrifty remote", func() {
				config.Thrifty = true
				thrifty := config.GetThrifty()
				Ω(thrifty).Should(HaveLen(1))
				Ω(thrifty).Should(ContainElement(uint32(3)))
			})

			It("should return 1 thrifty remote circularly", func() {
				config.Thrifty = true
				config.Name = "charlie"
				thrifty := config.GetThrifty()
				Ω(thrifty).Should(HaveLen(1))
				Ω(thrifty).Should(ContainElement(uint32(1)))
			})

			It("should return a quorum of 2 replicas", func() {
				Ω(config.GetQuorum()).Should(Equal(uint32(2)))
			})

		})

		Describe("size 5 quorum", func() {

			BeforeEach(func() {
				config.Peers = config.Peers[:5]
				Ω(config.Peers).Should(HaveLen(5))
			})

			It("should get 4 remotes to connect to", func() {
				remotes, err := config.GetRemotes()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(remotes).Should(HaveLen(4))

				for _, remote := range remotes {
					Ω(remote.PID).ShouldNot(Equal(uint32(2)))
				}

			})

			It("should return nil if not thrifty", func() {
				config.Thrifty = false
				Ω(config.GetThrifty()).Should(BeNil())
			})

			It("should return 2 thrifty remotes", func() {
				config.Thrifty = true
				thrifty := config.GetThrifty()
				Ω(thrifty).Should(HaveLen(2))
				Ω(thrifty).Should(ContainElement(uint32(3)))
				Ω(thrifty).Should(ContainElement(uint32(4)))
			})

			It("should return 2 thrifty remotes circularly", func() {
				config.Thrifty = true
				config.Name = "delta"
				thrifty := config.GetThrifty()
				Ω(thrifty).Should(HaveLen(2))
				Ω(thrifty).Should(ContainElement(uint32(5)))
				Ω(thrifty).Should(ContainElement(uint32(1)))
			})

			It("should return a quorum of 3 replicas", func() {
				Ω(config.GetQuorum()).Should(Equal(uint32(3)))
			})

		})

		Describe("size 7 quorum", func() {

			BeforeEach(func() {
				config.Peers = config.Peers[:7]
				Ω(config.Peers).Should(HaveLen(7))
			})

			It("should get 6 remotes to connect to", func() {
				remotes, err := config.GetRemotes()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(remotes).Should(HaveLen(6))

				for _, remote := range remotes {
					Ω(remote.PID).ShouldNot(Equal(uint32(2)))
				}

			})

			It("should return nil if not thrifty", func() {
				config.Thrifty = false
				Ω(config.GetThrifty()).Should(BeNil())
			})

			It("should return 3 thrifty remotes", func() {
				config.Thrifty = true
				thrifty := config.GetThrifty()
				Ω(thrifty).Should(HaveLen(3))
				Ω(thrifty).Should(ContainElement(uint32(3)))
				Ω(thrifty).Should(ContainElement(uint32(4)))
				Ω(thrifty).Should(ContainElement(uint32(5)))
			})

			It("should return 3 thrifty remotes circularly", func() {
				config.Thrifty = true
				config.Name = "echo"
				thrifty := config.GetThrifty()
				Ω(thrifty).Should(HaveLen(3))
				Ω(thrifty).Should(ContainElement(uint32(6)))
				Ω(thrifty).Should(ContainElement(uint32(7)))
				Ω(thrifty).Should(ContainElement(uint32(1)))
			})

			It("should return a quorum of 4 replicas", func() {
				Ω(config.GetQuorum()).Should(Equal(uint32(4)))
			})

		})

	})

})
