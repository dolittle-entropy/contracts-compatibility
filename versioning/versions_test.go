package versioning_test

import (
	"dolittle.io/contracts-compatibility/versioning"
	"github.com/coreos/go-semver/semver"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewReleaseVersions", func() {
	var toParse []string
	var parsed semver.Versions
	JustBeforeEach(func() {
		parsed = versioning.NewReleaseVersions(toParse)
	})

	When("there are no versions to parse", func() {
		BeforeEach(func() {
			toParse = []string{}
		})

		It("should return an empty slice", func() {
			Expect(parsed).To(BeEmpty())
		})
	})

	When("there are three valid versions to parse", func() {
		BeforeEach(func() {
			toParse = []string{
				"1.2.3",
				"4.5.6",
				"7.8.9",
			}
		})

		It("should return the three versions", func() {
			Expect(parsed).To(ConsistOf(
				semver.New("1.2.3"),
				semver.New("4.5.6"),
				semver.New("7.8.9"),
			))
		})
	})

	When("there are two releases and one pre-release", func() {
		BeforeEach(func() {
			toParse = []string{
				"1.2.3",
				"4.5.6-something",
				"7.8.9",
			}
		})

		It("should return the two released versions", func() {
			Expect(parsed).To(ConsistOf(
				semver.New("1.2.3"),
				semver.New("7.8.9"),
			))
		})
	})

	When("there is one invalid and one valid version", func() {
		BeforeEach(func() {
			toParse = []string{
				"some random stuff",
				"7.8.9",
			}
		})

		It("should just the valid version", func() {
			Expect(parsed).To(ConsistOf(
				semver.New("7.8.9"),
			))
		})
	})
})
