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
		daysOfActivity := 1

		g.It("Should find active repos for RobotsAndPencils", func() {
			repos, err := s.loadActiveReposForOrganization("RobotsAndPencils", daysOfActivity)

			Expect(repos).ToNot(BeNil())
			Expect(err).To(BeNil())
		})

		g.It("Should find open pull requests for active repos in RobotsAndPencils", func() {
			pullRequests, err := s.loadOpenPRsForOrganization("RobotsAndPencils", daysOfActivity)

			Expect(pullRequests).ToNot(BeNil())
			Expect(err).To(BeNil())
		})

		g.It("Should find product backlog items in pencilcase", func() {
			issues, err := s.Backlog("RobotsAndPencils", "pencilcase")

			Expect(issues).ToNot(BeNil())
			Expect(err).To(BeNil())
		})

		g.It("Should not find sprint backlog items in marvin", func() {
			issues, err := s.Sprint("RobotsAndPencils", "marvin")

			Expect(issues).To(BeNil())
			Expect(err).To(BeNil())
		})

		g.It("Should find in progress items in pencilcase", func() {
			issues, err := s.InProgress("RobotsAndPencils", "pencilcase")

			Expect(issues).ToNot(BeNil())
			Expect(err).To(BeNil())
		})

		g.It("Should find ready for QA items in pencilcase", func() {
			issues, err := s.ReadyForQA("RobotsAndPencils", "pencilcase")

			Expect(issues).ToNot(BeNil())
			Expect(err).To(BeNil())
		})

		g.It("Should not find ready for review items in marvin", func() {
			issues, err := s.ReadyForReview("RobotsAndPencils", "marvin")

			Expect(issues).To(BeNil())
			Expect(err).To(BeNil())
		})

		g.It("Should not find passed QA items in marvin", func() {
			issues, err := s.QAPass("RobotsAndPencils", "marvin")

			Expect(issues).To(BeNil())
			Expect(err).To(BeNil())
		})

		g.It("Should load all of the issues for gambit", func() {
			issues, _ := s.loadIssuesForRepo("RobotsAndPencils", "gambit", "")

			Expect(len(issues)).To(BeEquivalentTo(59))
		})

		g.It("Should find 4 sprint tasks in gambit", func() {
			issues, err := s.Sprint("RobotsAndPencils", "gambit")

			Expect(issues).ToNot(BeNil())
			Expect(len(issues)).To(BeEquivalentTo(5))
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

		g.It("Should find commits for a repo", func() {
			commits, err := s.loadCommitsForRepo("RobotsAndPencils", "marvin", "")

			Expect(commits).ToNot(BeNil())
			Expect(err).To(BeNil())
		})

		g.It("Should find PR commits for a repo", func() {
			commits, err := s.loadCommitsFromAllRepoPRs("RobotsAndPencils", "marvin")

			Expect(commits).ToNot(BeNil())
			Expect(err).To(BeNil())
		})

		g.It("Should find commits to master", func() {
			commits, err := s.CommitsToMaster("RobotsAndPencils", "marvin")

			Expect(commits).ToNot(BeNil())
			Expect(err).To(BeNil())
		})
	})
}
