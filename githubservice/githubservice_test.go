package githubservice

/* Admittedly, this test file is not going to be of any use to anyone else. */

import (
	"fmt"
	"strconv"
	"testing"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func Test(t *testing.T) {
	g := Goblin(t)

	// special hook for gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })
	s := New("PERSONAL ACCESS TOKEN HERE")

	g.Describe("Github Service", func() {

		g.It("Should find product backlog items in marvin", func() {
			issues, err := s.Backlog("RobotsAndPencils", "marvin")

			Expect(issues).ToNot(BeNil())
			Expect(err).To(BeNil())
		})

		g.It("Should not find sprint backlog items in marvin", func() {
			issues, err := s.Sprint("RobotsAndPencils", "marvin")

			Expect(issues).To(BeNil())
			Expect(err).To(BeNil())
		})

		g.It("Should find in progress items in marvin", func() {
			issues, err := s.InProgress("RobotsAndPencils", "marvin")

			Expect(issues).ToNot(BeNil())
			Expect(err).To(BeNil())
		})

		g.It("Should find ready for QA items in marvin", func() {
			issues, err := s.ReadyForQA("RobotsAndPencils", "marvin")

			Expect(issues).To(BeNil())
			Expect(err).To(BeNil())
		})

		g.It("Should find ready for review items in pencilcase", func() {
			issues, err := s.ReadyForReview("RobotsAndPencils", "pencilcase")

			Expect(issues).ToNot(BeNil())
			Expect(err).To(BeNil())
		})

		g.It("Should not find passed QA items in marvin", func() {
			issues, err := s.QAPass("RobotsAndPencils", "marvin")

			Expect(issues).To(BeNil())
			Expect(err).To(BeNil())
		})

		g.It("Should load all of the issues for gambit", func() {
			issues, _ := s.loadIssuesForRepo("RobotsAndPencils", "gambit", "")

			Expect(len(issues)).To(BeEquivalentTo(36))
		})

		g.It("Should find 4 sprint tasks in gambit", func() {
			issues, err := s.Sprint("RobotsAndPencils", "gambit")

			Expect(issues).ToNot(BeNil())
			Expect(len(issues)).To(BeEquivalentTo(3))
			Expect(err).To(BeNil())
		})

		g.It("Should find items assigned to nealsanche in gambitandroid", func() {
			issues, err := s.AssignedTo("RobotsAndPencils", "gambitandroid", "nealsanche")

			Expect(issues).ToNot(BeNil())
			Expect(err).To(BeNil())
		})

		g.It("Should load issues for assignee by org", func() {
			issues, err := s.loadIssuesForAssignee("RobotsAndPencils", "nealsanche")

			Expect(issues).ToNot(BeNil())
			Expect(err).To(BeNil())

			fmt.Println("Issue count: " + strconv.Itoa(len(issues)))
		})

	})
}
