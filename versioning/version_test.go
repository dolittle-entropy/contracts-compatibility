package versioning_test

import (
	"dolittle.io/contracts-compatibility/versioning"
	"github.com/coreos/go-semver/semver"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewReleaseVersion", func() {
	var toParse string
	var parsed *semver.Version
	var err error
	JustBeforeEach(func() {
		parsed, err = versioning.NewReleaseVersion(toParse)
	})

	When("the version is not a valid SemVer", func() {
		BeforeEach(func() {
			toParse = "some string"
		})

		It("should fail", func() {
			Expect(parsed).To(BeNil())
			Expect(err).ToNot(BeNil())
		})
	})

	When("the version is a valid pre-release SemVer", func() {
		BeforeEach(func() {
			toParse = "1.2.3-something"
		})

		It("should fail", func() {
			Expect(parsed).To(BeNil())
			Expect(err).ToNot(BeNil())
		})
	})

	When("the version is a valid release SemVer", func() {
		BeforeEach(func() {
			toParse = "4.5.6"
		})

		It("should not fail", func() {
			Expect(parsed.Major).To(Equal(int64(4)))
			Expect(parsed.Minor).To(Equal(int64(5)))
			Expect(parsed.Patch).To(Equal(int64(6)))
			Expect(err).ToNot(BeNil())
		})
	})
})
