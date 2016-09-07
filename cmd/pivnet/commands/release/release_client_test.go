package release_test

import (
	"bytes"
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	pivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/go-pivnet/cmd/pivnet/commands/release"
	"github.com/pivotal-cf/go-pivnet/cmd/pivnet/commands/release/releasefakes"
	"github.com/pivotal-cf/go-pivnet/cmd/pivnet/errorhandler/errorhandlerfakes"
	"github.com/pivotal-cf/go-pivnet/cmd/pivnet/printer"
)

var _ = Describe("release commands", func() {
	var (
		fakePivnetClient *releasefakes.FakePivnetClient

		fakeErrorHandler *errorhandlerfakes.FakeErrorHandler

		outBuffer bytes.Buffer

		releases []pivnet.Release

		client *release.ReleaseClient
	)

	BeforeEach(func() {
		fakePivnetClient = &releasefakes.FakePivnetClient{}

		outBuffer = bytes.Buffer{}

		fakeErrorHandler = &errorhandlerfakes.FakeErrorHandler{}

		releases = []pivnet.Release{
			{
				ID: 1234,
			},
			{
				ID: 2345,
			},
		}

		client = release.NewReleaseClient(
			fakePivnetClient,
			fakeErrorHandler,
			printer.PrintAsJSON,
			&outBuffer,
			printer.NewPrinter(&outBuffer),
		)
	})

	Describe("List", func() {
		var (
			productSlug    string
			releaseVersion string
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = ""

			fakePivnetClient.ReleasesForProductSlugReturns(releases, nil)
		})

		It("lists all Releases", func() {
			err := client.List(productSlug)
			Expect(err).NotTo(HaveOccurred())

			var returnedReleases []pivnet.Release
			err = json.Unmarshal(outBuffer.Bytes(), &returnedReleases)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedReleases).To(Equal(releases))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("releases error")
				fakePivnetClient.ReleasesForProductSlugReturns(nil, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.List(productSlug)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("Get", func() {
		var (
			productSlug    string
			releaseVersion string
			releaseID      int
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = ""
			releaseID = releases[0].ID

			fakePivnetClient.ReleaseForProductVersionReturns(releases[0], nil)
			fakePivnetClient.ReleaseETagReturns("some-etag", nil)
		})

		It("gets Release", func() {
			err := client.Get(productSlug, releaseVersion)
			Expect(err).NotTo(HaveOccurred())

			var returnedRelease pivnet.Release
			err = json.Unmarshal(outBuffer.Bytes(), &returnedRelease)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnedRelease).To(Equal(releases[0]))
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("release error")
				fakePivnetClient.ReleaseForProductVersionReturns(pivnet.Release{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Get(productSlug, releaseVersion)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when there is an error getting etag", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("release error")
				fakePivnetClient.ReleaseETagReturns("", expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Get(productSlug, releaseVersion)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})

	Describe("Delete", func() {
		var (
			productSlug    string
			releaseVersion string
		)

		BeforeEach(func() {
			productSlug = "some-product-slug"
			releaseVersion = releases[0].Version

			fakePivnetClient.ReleaseForProductVersionReturns(releases[0], nil)
			fakePivnetClient.DeleteReleaseReturns(nil)
		})

		It("deletes Release", func() {
			err := client.Delete(productSlug, releaseVersion)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("release error")
				fakePivnetClient.DeleteReleaseReturns(expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Delete(productSlug, releaseVersion)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})

		Context("when there is an error getting release", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = errors.New("release error")
				fakePivnetClient.ReleaseForProductVersionReturns(pivnet.Release{}, expectedErr)
			})

			It("invokes the error handler", func() {
				err := client.Delete(productSlug, releaseVersion)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
				Expect(fakeErrorHandler.HandleErrorArgsForCall(0)).To(Equal(expectedErr))
			})
		})
	})
})