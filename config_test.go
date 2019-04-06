package epaxos_test

import (
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

})
