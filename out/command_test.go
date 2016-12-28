package out_test

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mcb/cf-service-resource"
	"github.com/mcb/cf-service-resource/out"
	"github.com/mcb/cf-service-resource/out/fakes"
)

var _ = Describe("Out Command", func() {
	var (
		cloudFoundry *fakes.FakePAAS
		request      out.Request
		command      *out.Command
	)

	BeforeEach(func() {
		cloudFoundry = &fakes.FakePAAS{}
		command = out.NewCommand(cloudFoundry)

		request = out.Request{
			Source: resource.Source{
				API:          "https://api.run.pivotal.io",
				Username:     "awesome@example.com",
				Password:     "hunter2",
				Organization: "secret",
				Space:        "volcano-base",
			},
			Params: out.Params{
				Service: "p-mysql",
				Plan: "512mb",
				InstanceName: "mysql-test",
				CurrentAppName: "foobar",
			},
		}
	})

	Describe("running the command", func() {
		It("pushes an application into cloud foundry", func() {
			response, err := command.Run(request)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(response.Version.Timestamp).Should(BeTemporally("~", time.Now(), time.Second))
			Ω(response.Metadata[0]).Should(Equal(
				resource.MetadataPair{
					Name:  "organization",
					Value: "secret",
				},
			))
			Ω(response.Metadata[1]).Should(Equal(
				resource.MetadataPair{
					Name:  "space",
					Value: "volcano-base",
				},
			))

			By("logging in")
			Ω(cloudFoundry.LoginCallCount()).Should(Equal(1))

			api, username, password, insecure := cloudFoundry.LoginArgsForCall(0)
			Ω(api).Should(Equal("https://api.run.pivotal.io"))
			Ω(username).Should(Equal("awesome@example.com"))
			Ω(password).Should(Equal("hunter2"))
			Ω(insecure).Should(Equal(false))

			By("targetting the organization and space")
			Ω(cloudFoundry.TargetCallCount()).Should(Equal(1))

			org, space := cloudFoundry.TargetArgsForCall(0)
			Ω(org).Should(Equal("secret"))
			Ω(space).Should(Equal("volcano-base"))

			By("create service instance")
			Ω(cloudFoundry.CreateServiceCallCount()).Should(Equal(1))

			service, plan, instanceName := cloudFoundry.CreateServiceArgsForCall(0)

			Ω(service).Should(Equal("p-mysql"))
			Ω(plan).Should(Equal("512mb"))
			Ω(instanceName).Should(Equal("mysql-test"))

			By("bind service instance to app")
			Ω(cloudFoundry.BindServiceCallCount()).Should(Equal(1))

			currentAppName, instanceName := cloudFoundry.BindServiceArgsForCall(0)

			Ω(currentAppName).Should(Equal("foobar"))
			Ω(instanceName).Should(Equal("mysql-test"))

			By("restage app")
			Ω(cloudFoundry.RestageAppCallCount()).Should(Equal(1))

			currentAppName = cloudFoundry.RestageAppArgsForCall(0)

			Ω(currentAppName).Should(Equal("foobar"))

		})

		Describe("handling any errors", func() {
			var expectedError error

			BeforeEach(func() {
				expectedError = errors.New("it all went wrong")
			})

			It("from logging in", func() {
				cloudFoundry.LoginReturns(expectedError)

				_, err := command.Run(request)
				Ω(err).Should(MatchError(expectedError))
			})

			It("from targetting an org and space", func() {
				cloudFoundry.TargetReturns(expectedError)

				_, err := command.Run(request)
				Ω(err).Should(MatchError(expectedError))
			})

			It("from create service instance", func() {
				cloudFoundry.CreateServiceReturns(expectedError)

				_, err := command.Run(request)
				Ω(err).Should(MatchError(expectedError))
			})

			It("from binding service to app", func() {
				cloudFoundry.BindServiceReturns(expectedError)

				_, err := command.Run(request)
				Ω(err).Should(MatchError(expectedError))
			})

			It("from restage app", func() {
				cloudFoundry.RestageAppReturns(expectedError)

				_, err := command.Run(request)
				Ω(err).Should(MatchError(expectedError))
			})
		})

		It("lets people skip the certificate check", func() {
			request = out.Request{
				Source: resource.Source{
					API:           "https://api.run.pivotal.io",
					Username:      "awesome@example.com",
					Password:      "hunter2",
					Organization:  "secret",
					Space:         "volcano-base",
					SkipCertCheck: true,
				},
				Params: out.Params{
					Service:        "foo",
					Plan:           "bar",
					InstanceName:   "baz",
					CurrentAppName: "fox",
				},
			}

			_, err := command.Run(request)
			Ω(err).ShouldNot(HaveOccurred())

			By("logging in")
			Ω(cloudFoundry.LoginCallCount()).Should(Equal(1))

			_, _, _, insecure := cloudFoundry.LoginArgsForCall(0)
			Ω(insecure).Should(Equal(true))
		})
	})
})
