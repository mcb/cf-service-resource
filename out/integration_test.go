package out_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"github.com/mcb/cf-service-resource"
	"github.com/mcb/cf-service-resource/out"
)

var _ = Describe("Out", func() {
	var (
		tmpDir  string
		cmd     *exec.Cmd
		request out.Request
	)

	BeforeEach(func() {
		var err error

		tmpDir, err = ioutil.TempDir("", "cf_resource_out")
		Ω(err).ShouldNot(HaveOccurred())

		err = os.Mkdir(filepath.Join(tmpDir, "project"), 0755)
		Ω(err).ShouldNot(HaveOccurred())

		err = ioutil.WriteFile(filepath.Join(tmpDir, "project", "manifest.yml"), []byte{}, 0555)
		Ω(err).ShouldNot(HaveOccurred())

		request = out.Request{
			Source: resource.Source{
				API:           "https://api.run.pivotal.io",
				Username:      "awesome@example.com",
				Password:      "hunter2",
				Organization:  "org",
				Space:         "space",
				SkipCertCheck: true,
			},
			Params: out.Params{
				Service:        "p-mysql",
				Plan:           "512mb",
				InstanceName:   "mysql-test",
				ManifestPath:   "project/manifest.yml",
				CurrentAppName: "bar-app",
			},
		}
	})

	JustBeforeEach(func() {
		assetsPath, err := filepath.Abs("assets")
		Ω(err).ShouldNot(HaveOccurred())

		stdin := &bytes.Buffer{}

		err = json.NewEncoder(stdin).Encode(request)
		Ω(err).ShouldNot(HaveOccurred())

		cmd = exec.Command(binPath, tmpDir)
		cmd.Stdin = stdin
		cmd.Dir = tmpDir

		newEnv := []string{}
		for _, envVar := range os.Environ() {
			if strings.HasPrefix(envVar, "PATH=") {
				newEnv = append(newEnv, fmt.Sprintf("PATH=%s:%s", assetsPath, os.Getenv("PATH")))
			} else {
				newEnv = append(newEnv, envVar)
			}
		}

		cmd.Env = newEnv
	})

	AfterEach(func() {
		err := os.RemoveAll(tmpDir)
		Ω(err).ShouldNot(HaveOccurred())
	})

	Context("when my manifest do not contain a glob", func() {
		It("returns an error", func() {
			session, err := gexec.Start(
				cmd,
				GinkgoWriter,
				GinkgoWriter,
			)
			Ω(err).ShouldNot(HaveOccurred())

			Eventually(session).Should(gexec.Exit(1))

			errMsg := fmt.Sprintf("error invalid manifest file: applications required")
			Ω(session.Err).Should(gbytes.Say(errMsg))
		})
	})

	Context("when my manifest contains a glob", func() {
		var tmpFileManifest *os.File

		BeforeEach(func() {
			var err error

			tmpFileManifest, err = ioutil.TempFile(tmpDir, "manifest-some-glob.yml_")
			Ω(err).ShouldNot(HaveOccurred())
			_, err = tmpFileManifest.WriteString("applications:\n - name: foo-app\n")
			Ω(err).ShouldNot(HaveOccurred())

			request.Params.ManifestPath = "manifest-*.yml_*"
		})

		Context("when one file matches", func() {
			It("create/bind service and restage an application to Cloud Foundry", func() {
				session, err := gexec.Start(
					cmd,
					GinkgoWriter,
					GinkgoWriter,
				)
				Ω(err).ShouldNot(HaveOccurred())

				Eventually(session).Should(gexec.Exit(0))

				var response out.Response
				err = json.Unmarshal(session.Out.Contents(), &response)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(response.Version.Timestamp).Should(BeTemporally("~", time.Now(), time.Second))

				// shim outputs arguments
				Ω(session.Err).Should(gbytes.Say("cf api https://api.run.pivotal.io --skip-ssl-validation"))
				Ω(session.Err).Should(gbytes.Say("cf auth awesome@example.com hunter2"))
				Ω(session.Err).Should(gbytes.Say("cf target -o org -s space"))
				Ω(session.Err).Should(gbytes.Say("cf create-service p-mysql 512mb mysql-test"))
				Ω(session.Err).Should(gbytes.Say("cf bind-service foo-app mysql-test"))
				Ω(session.Err).Should(gbytes.Say("cf restage foo-app"))

				// color should be always
				Ω(session.Err).Should(gbytes.Say("CF_COLOR=true"))
			})
		})

		Context("when one file matches but name is missing", func() {
			BeforeEach(func() {
				var err error
				err = os.Remove(tmpFileManifest.Name())
				Ω(err).ShouldNot(HaveOccurred())
				tmpFileManifest, err = ioutil.TempFile(tmpDir, "manifest-some-glob.yml_")
				Ω(err).ShouldNot(HaveOccurred())
				_, err = tmpFileManifest.WriteString("applications:\n - env: foo=bar\n")
				Ω(err).ShouldNot(HaveOccurred())
			})
			It("returns an error", func() {
				session, err := gexec.Start(
					cmd,
					GinkgoWriter,
					GinkgoWriter,
				)
				Ω(err).ShouldNot(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))

				errMsg := fmt.Sprintf("error invalid manifest file: name required")
				Ω(session.Err).Should(gbytes.Say(errMsg))
			})
		})

		Context("when no files match the manifest path", func() {
			BeforeEach(func() {
				request.Params.ManifestPath = "nope-*"
			})

			It("returns an error", func() {
				session, err := gexec.Start(
					cmd,
					GinkgoWriter,
					GinkgoWriter,
				)
				Ω(err).ShouldNot(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))
				errMsg := fmt.Sprintf("error invalid manifest path: found 0 files instead of 1 at path: %s", filepath.Join(tmpDir, `nope-\*`))
				Ω(session.Err).Should(gbytes.Say(errMsg))
			})
		})

		Context("when more than one file matches the manifest path", func() {
			BeforeEach(func() {
				_, err := ioutil.TempFile(tmpDir, "manifest-some-glob.yml_")
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("returns an error", func() {
				session, err := gexec.Start(
					cmd,
					GinkgoWriter,
					GinkgoWriter,
				)
				Ω(err).ShouldNot(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))
				errMsg := fmt.Sprintf("error invalid manifest path: found 2 files instead of 1 at path: %s", filepath.Join(tmpDir, `manifest-\*.yml_\*`))
				Ω(session.Err).Should(gbytes.Say(errMsg))
			})
		})
	})

	Context("when manifest is empty but current_app_name is filled", func() {
		BeforeEach(func() {
			request.Params.ManifestPath =""
		})

		It("create/bind service and restage an application to Cloud foundry", func() {
			session, err := gexec.Start(
				cmd,
				GinkgoWriter,
				GinkgoWriter,
			)
			Ω(err).ShouldNot(HaveOccurred())

			Eventually(session).Should(gexec.Exit(0))

			var response out.Response
			err = json.Unmarshal(session.Out.Contents(), &response)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(response.Version.Timestamp).Should(BeTemporally("~", time.Now(), time.Second))

			// shim outputs arguments
			Ω(session.Err).Should(gbytes.Say("cf api https://api.run.pivotal.io --skip-ssl-validation"))
			Ω(session.Err).Should(gbytes.Say("cf auth awesome@example.com hunter2"))
			Ω(session.Err).Should(gbytes.Say("cf target -o org -s space"))
			Ω(session.Err).Should(gbytes.Say("cf create-service p-mysql 512mb mysql-test"))
			Ω(session.Err).Should(gbytes.Say("cf bind-service bar-app mysql-test"))
			Ω(session.Err).Should(gbytes.Say("cf restage bar-app"))

			// color shoud be always
			Ω(session.Err).Should(gbytes.Say("CF_COLOR=true"))
		})
	})

	Context("when required option is empty", func() {
		Context("service is empty", func() {
			BeforeEach(func() {
				request.Params.Service = ""
			})

			It("return an error", func() {
				session, err := gexec.Start(
					cmd,
					GinkgoWriter,
					GinkgoWriter,
				)
				Ω(err).ShouldNot(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))

				errMsg := fmt.Sprintf("error parameter required: service")
				Ω(session.Err).Should(gbytes.Say(errMsg))
			})
		})

		Context("plan is empty", func() {
			BeforeEach(func() {
				request.Params.Plan = ""
			})

			It("return an error", func() {
				session, err := gexec.Start(
					cmd,
					GinkgoWriter,
					GinkgoWriter,
				)
				Ω(err).ShouldNot(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))

				errMsg := fmt.Sprintf("error parameter required: plan")
				Ω(session.Err).Should(gbytes.Say(errMsg))
			})
		})

                Context("instance_name is empty", func() {
                        BeforeEach(func() {
                                request.Params.InstanceName = ""
                        })

                        It("return an error", func() {
                                session, err := gexec.Start(
                                        cmd,
                                        GinkgoWriter,
                                        GinkgoWriter,
                                )
                                Ω(err).ShouldNot(HaveOccurred())

                                Eventually(session).Should(gexec.Exit(1))

                                errMsg := fmt.Sprintf("error parameter required: instance_name")
                                Ω(session.Err).Should(gbytes.Say(errMsg))
                        })
                })

                Context("manifest and current_app_name is empty", func() {
                        BeforeEach(func() {
                                request.Params.ManifestPath = ""
				request.Params.CurrentAppName = ""
                        })

                        It("return an error", func() {
                                session, err := gexec.Start(
                                        cmd,
                                        GinkgoWriter,
                                        GinkgoWriter,
                                )
                                Ω(err).ShouldNot(HaveOccurred())

                                Eventually(session).Should(gexec.Exit(1))

                                errMsg := fmt.Sprintf("error parameter required: either manifest/current_app_name")
                                Ω(session.Err).Should(gbytes.Say(errMsg))
                        })
                })
	})
})
