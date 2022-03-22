package main_test

import (
	main "dolittle.io/contracts-compatibility"
	"github.com/coreos/go-semver/semver"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("VersionsCompatibilityIsOverridden", func() {
	var runtimeVersion *semver.Version
	var sdk string
	var sdkVersion *semver.Version
	BeforeEach(func() {
		runtimeVersion = &semver.Version{}
		sdk = ""
		sdkVersion = &semver.Version{}
	})

	var isOverridden bool
	JustBeforeEach(func() {
		isOverridden = main.VersionsCompatibilityIsOverridden(runtimeVersion, sdk, sdkVersion)
	})

	When("checking versions before and after ReverseCall ping-pong change", func() {
		Context("for DotNet SDK", func() {
			BeforeEach(func() {
				sdk = "DotNET"
			})

			Context("and the SDK is too old", func() {
				BeforeEach(func() {
					runtimeVersion.Major = 6
					sdkVersion.Major = 8
				})

				It("should be overridden", func() {
					Expect(isOverridden).To(BeTrue())
				})
			})

			Context("and the Runtime is too old", func() {
				BeforeEach(func() {
					runtimeVersion.Major = 5
					sdkVersion.Major = 9
				})

				It("should be overridden", func() {
					Expect(isOverridden).To(BeTrue())
				})
			})

			Context("and both are old", func() {
				BeforeEach(func() {
					runtimeVersion.Major = 4
					sdkVersion.Major = 6
				})

				It("should not be overridden", func() {
					Expect(isOverridden).To(BeFalse())
				})
			})

			Context("and both are new", func() {
				BeforeEach(func() {
					runtimeVersion.Major = 12
					sdkVersion.Major = 15
				})

				It("should not be overridden", func() {
					Expect(isOverridden).To(BeFalse())
				})
			})
		})

		Context("for JavaScript SDK", func() {
			BeforeEach(func() {
				sdk = "JavaScript"
			})

			Context("and the SDK is too old", func() {
				BeforeEach(func() {
					runtimeVersion.Major = 7
					sdkVersion.Major = 13
				})

				It("should be overridden", func() {
					Expect(isOverridden).To(BeTrue())
				})
			})

			Context("and the Runtime is too old", func() {
				BeforeEach(func() {
					runtimeVersion.Major = 5
					sdkVersion.Major = 15
				})

				It("should be overridden", func() {
					Expect(isOverridden).To(BeTrue())
				})
			})

			Context("and both are old", func() {
				BeforeEach(func() {
					runtimeVersion.Major = 4
					sdkVersion.Major = 14
				})

				It("should not be overridden", func() {
					Expect(isOverridden).To(BeFalse())
				})
			})

			Context("and both are new", func() {
				BeforeEach(func() {
					runtimeVersion.Major = 6
					sdkVersion.Major = 15
				})

				It("should not be overridden", func() {
					Expect(isOverridden).To(BeFalse())
				})
			})
		})
	})
})
