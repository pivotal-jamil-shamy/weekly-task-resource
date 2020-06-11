package main_test

import (
	"encoding/json"
	"os/exec"

	"github.com/pivotal-cf-experimental/cron-resource/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Check", func() {
	var checkCmd *exec.Cmd

	BeforeEach(func() {
		checkCmd = exec.Command(checkPath)
	})

	var request models.CheckRequest
	var session *gexec.Session

	BeforeEach(func() {
		request = models.CheckRequest{
			Source: models.Source{
				Location:   "America/Toronto",
				HourToFire: 17,
				DayToFire:  "Sunday",
			},
		}
	})

	JustBeforeEach(func() {
		stdin, err := checkCmd.StdinPipe()
		Expect(err).ShouldNot(HaveOccurred())

		session, err = gexec.Start(checkCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())

		err = json.NewEncoder(stdin).Encode(request)
		Expect(err).ShouldNot(HaveOccurred())
	})

	Describe("resource config validation", func() {
		Describe("Day To Fire", func() {

			Context("when given day is invalid", func() {
				BeforeEach(func() {
					request.Source.DayToFire = "snoopyday"
				})

				It("exits with status code 1 and a proper message", func() {
					Eventually(session.Err).Should(gbytes.Say("\"day_to_fire\" should be one of the following: \"Sunday,Monday,Tuesday,Wednesday,Thursday,Friday,Saturday\""))
					Eventually(session).Should(gexec.Exit(1))
				})
			})

			Context("when given day is a valid weekday", func() {
				BeforeEach(func() {
					request.Source.DayToFire = "Saturday"
				})

				It("exits with status code 0 and no errors", func() {
					Eventually(session.Err.Closed).Should(BeTrue())
					Eventually(session.Err.Contents).Should(BeEmpty())
					Eventually(session).Should(gexec.Exit(0))
				})

				Context("when the given day is mixed case", func() {
					BeforeEach(func() {
						request.Source.DayToFire = "fRiDaY"
					})

					It("exits with status code 0 and no errors", func() {
						Eventually(session.Err.Closed).Should(BeTrue())
						Eventually(session.Err.Contents).Should(BeEmpty())
						Eventually(session).Should(gexec.Exit(0))
					})
				})

				Context("when the given day is not trimmed", func() {
					BeforeEach(func() {
						request.Source.DayToFire = "    Tuesday    "
					})

					It("exits with status code 0 and no errors", func() {
						Eventually(session.Err.Closed).Should(BeTrue())
						Eventually(session.Err.Contents).Should(BeEmpty())
						Eventually(session).Should(gexec.Exit(0))
					})
				})
			})

		})

		Describe("Hour To Fire", func() {
			Context("when given hour is less than zero", func() {
				BeforeEach(func() {
					request.Source.HourToFire = -2
				})

				It("exits with status code 1 and a proper message", func() {
					Eventually(session.Err).Should(gbytes.Say("\"hour_to_fire\" should be in the 0-23 range"))
					Eventually(session).Should(gexec.Exit(1))
				})
			})

			Context("when given hour is greater than 23", func() {
				BeforeEach(func() {
					request.Source.HourToFire = 25
				})

				It("exits with status code 1 and a proper message", func() {
					Eventually(session.Err).Should(gbytes.Say("\"hour_to_fire\" should be in the 0-23 range"))
					Eventually(session).Should(gexec.Exit(1))
				})
			})

			Context("when given hour is within range", func() {
				BeforeEach(func() {
					request.Source.HourToFire = 23
				})

				It("exits with status code 0 and no error", func() {
					Eventually(session.Err.Closed).Should(BeTrue())
					Eventually(session.Err.Contents).Should(BeEmpty())
					Eventually(session).Should(gexec.Exit(0))
				})
			})
		})

		Describe("Timezone", func() {
			Context("when given timezone is invalid", func() {
				BeforeEach(func() {
					request.Source.Location = "America/China"
				})

				It("exits with status code 1 and a proper message", func() {
					Eventually(session.Err).Should(gbytes.Say("unknown time zone"))
					Eventually(session).Should(gexec.Exit(1))
				})
			})

			Context("when given timezone is valid", func() {
				BeforeEach(func() {
					request.Source.Location = "Asia/Beirut"
				})

				It("exits with status code 0 and no error", func() {
					Eventually(session.Err.Closed).Should(BeTrue())
					Eventually(session.Err.Contents).Should(BeEmpty())
					Eventually(session).Should(gexec.Exit(0))
				})
			})
		})
	})
})
