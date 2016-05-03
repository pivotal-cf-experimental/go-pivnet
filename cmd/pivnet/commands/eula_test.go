package commands_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/go-pivnet"
	"github.com/pivotal-cf-experimental/go-pivnet/cmd/pivnet/commands"

	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("eula commands", func() {
	var (
		server *ghttp.Server
		host   string

		eulas []pivnet.EULA

		outBuffer bytes.Buffer
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		host = server.URL()
		commands.Pivnet.Host = host
		outBuffer = bytes.Buffer{}
		commands.OutWriter = &outBuffer

		eulas = []pivnet.EULA{
			{
				ID:   1234,
				Name: "some eula",
				Slug: "some-eula",
			},
			{
				ID:   2345,
				Name: "another eula",
				Slug: "another-eula",
			},
		}
	})

	AfterEach(func() {
		server.Close()
	})

	It("lists all EULAs", func() {
		eulasResponse := pivnet.EULAsResponse{
			EULAs: eulas,
		}

		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", fmt.Sprintf("%s/eulas", apiPrefix)),
				ghttp.RespondWithJSONEncoded(http.StatusOK, eulasResponse),
			),
		)

		eulasCommand := commands.EULAsCommand{}

		err := eulasCommand.Execute(nil)
		Expect(err).NotTo(HaveOccurred())

		var returnedEULAs []pivnet.EULA

		err = json.Unmarshal(outBuffer.Bytes(), &returnedEULAs)
		Expect(err).NotTo(HaveOccurred())

		Expect(returnedEULAs).To(Equal(eulas))
	})

	It("shows specific EULA", func() {
		eulaResponse := eulas[0]

		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", fmt.Sprintf("%s/eulas/%s", apiPrefix, eulas[0].Slug)),
				ghttp.RespondWithJSONEncoded(http.StatusOK, eulaResponse),
			),
		)

		eulaCommand := commands.EULACommand{}
		eulaCommand.EULASlug = eulas[0].Slug

		err := eulaCommand.Execute(nil)
		Expect(err).NotTo(HaveOccurred())

		var returnedEULA pivnet.EULA

		err = json.Unmarshal(outBuffer.Bytes(), &returnedEULA)
		Expect(err).NotTo(HaveOccurred())

		Expect(returnedEULA).To(Equal(eulas[0]))
	})

	It("accepts EULA", func() {
		releases := []pivnet.Release{
			{
				ID:          1234,
				Version:     "version 0.2.3",
				Description: "Some release with some description.",
			},
			{
				ID:          2345,
				Version:     "version 0.3.4",
				Description: "Another release with another description.",
			},
		}

		releasesResponse := pivnet.ReleasesResponse{
			Releases: releases,
		}

		productSlug := "some-product-slug"

		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest(
					"GET",
					fmt.Sprintf("%s/products/%s/releases", apiPrefix, productSlug),
				),
				ghttp.RespondWithJSONEncoded(http.StatusOK, releasesResponse),
			),
		)

		eulaAcceptanceResponse := pivnet.EULAAcceptanceResponse{
			AcceptedAt: "now",
		}

		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest(
					"POST",
					fmt.Sprintf(
						"%s/products/%s/releases/%d/eula_acceptance",
						apiPrefix,
						productSlug,
						releases[0].ID,
					),
				),
				ghttp.RespondWithJSONEncoded(http.StatusOK, eulaAcceptanceResponse),
			),
		)

		acceptEULACommand := commands.AcceptEULACommand{}
		acceptEULACommand.ProductSlug = productSlug
		acceptEULACommand.ReleaseVersion = releases[0].Version

		err := acceptEULACommand.Execute(nil)
		Expect(err).NotTo(HaveOccurred())
	})
})
