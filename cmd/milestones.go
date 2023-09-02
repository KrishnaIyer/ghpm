// Copyright © 2023 Krishna Iyer Easwaran
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/go-github/v54/github"
	"github.com/spf13/cobra"
)

// Milestone is a GitHub repository milestone.
// This custom type exists since `github.Milestone` contains a lot of information that we don't need.
type Milestone struct {
	Title        string `json:"title,omitempty"`
	Description  string `json:"description,omitempty"`
	OpenIssues   int    `json:"open_issues,omitempty"`
	ClosedIssues int    `json:"closed_issues,omitempty"`
	ClosedAt     string `json:"closed_at,omitempty"`
	DueOn        string `json:"due_on,omitempty"`
}

var (
	milestonesCmd = &cobra.Command{
		Use:   "milestones",
		Short: "Manage milestones",
		Long:  `Manage milestones.`,
	}
	createMilestonesCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a milestone in all repositories",
		Long:  `Create a milestone in all the repositories in the configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			title, _ := cmd.Flags().GetString("title")
			description, _ := cmd.Flags().GetString("description")
			dueOn, _ := cmd.Flags().GetString("due-on")
			due, err := time.Parse("2006-01-02", dueOn)
			// Add 24 hours to the due date since GitHub API expects the due date to be at the end of the day.
			due = due.Add(24 * time.Hour)
			if err != nil {
				return err
			}
			for _, repo := range config.Repositories {
				userName := config.Username
				if repo.Username != "" {
					userName = repo.Username
				}
				milestone, _, err := client.Issues.CreateMilestone(
					ctx,
					userName,
					repo.Name,
					&github.Milestone{
						Title:       &title,
						Description: &description,
						DueOn:       &github.Timestamp{Time: due.UTC()},
					})
				if err != nil {
					return err
				}
				fmt.Printf("Created milestone %s in repository %s\n", *milestone.Title, repo.Name)
			}
			return nil
		},
	}
	getMilestonesCmd = &cobra.Command{
		Use:   "get",
		Short: "Get milestones",
		Long:  `Get milestones.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			overdue, _ := cmd.Flags().GetBool("overdue")
			closed, _ := cmd.Flags().GetBool("closed")

			byRepo := map[string][][]byte{}
			for _, repo := range config.Repositories {
				userName := config.Username
				if repo.Username != "" {
					userName = repo.Username
				}
				milestones, _, err := client.Issues.ListMilestones(
					ctx,
					userName,
					repo.Name,
					&github.MilestoneListOptions{
						State:     "all",
						Sort:      "due_on",
						Direction: "asc",
					})
				if err != nil {
					return err
				}
				for _, milestone := range milestones {
					var (
						dueOn    string
						closedAt string
					)
					if milestone.DueOn != nil {
						dueOn = milestone.DueOn.Format("02 Jan 06")
					}
					if overdue && time.Now().Before(milestone.DueOn.Time) {
						continue
					}
					if closed {
						if milestone.ClosedAt == nil {
							continue
						}
						closedAt = milestone.ClosedAt.Format("02 Jan 06 ")
					}

					m := Milestone{
						Title:        *milestone.Title,
						Description:  *milestone.Description,
						OpenIssues:   *milestone.OpenIssues,
						ClosedIssues: *milestone.ClosedIssues,
						ClosedAt:     closedAt,
						DueOn:        dueOn,
					}
					jsonValue, err := json.Marshal(m)
					if err != nil {
						return err
					}
					byRepo[repo.Name] = append(byRepo[repo.Name], jsonValue)
				}
			}
			fmt.Println("##################################################")
			fmt.Println("\t\t Milestones")
			fmt.Println("##################################################")
			for repo, milestones := range byRepo {
				fmt.Printf("Repository: %s\n", repo)
				for i, milestone := range milestones {
					fmt.Printf("%d. %s\n", i+1, milestone)
				}
				fmt.Println("--------------------------------------------------")
			}
			return nil
		},
	}
	updateMilestonesCmd = &cobra.Command{
		Use:   "update [title]",
		Short: "Update a milestone in all repositories",
		Long: `Update a milestone in all the repositories in the configuration.
Milestones not found are skipped.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("no title provided")
			}
			title := args[0]
			description, _ := cmd.Flags().GetString("description")
			dueOn, _ := cmd.Flags().GetString("due-on")
			due, err := time.Parse("2006-01-02", dueOn)
			// Add 24 hours to the due date since GitHub API expects the due date to be at the end of the day.
			due = due.Add(24 * time.Hour)
			fmt.Println(due)
			if err != nil {
				return err
			}
			for _, repo := range config.Repositories {
				userName := config.Username
				if repo.Username != "" {
					userName = repo.Username
				}
				milestones, _, err := client.Issues.ListMilestones(
					ctx,
					userName,
					repo.Name,
					&github.MilestoneListOptions{
						State:     "all",
						Sort:      "due_on",
						Direction: "asc",
					})
				if err != nil {
					return err
				}
				for _, milestone := range milestones {
					if *milestone.Title != title {
						continue
					}
					milestone, _, err := client.Issues.EditMilestone(
						ctx,
						userName,
						repo.Name,
						*milestone.Number,
						&github.Milestone{
							Description: &description,
							DueOn:       &github.Timestamp{Time: due.UTC()},
						})
					if err != nil {
						return err
					}
					fmt.Printf("Updated milestone %s in repository %s\n", *milestone.Title, repo.Name)
				}
			}
			return nil
		},
	}
)

func init() {
	Root.AddCommand(milestonesCmd)
	milestonesCmd.AddCommand(getMilestonesCmd)
	getMilestonesCmd.Flags().Bool("overdue", false, "Show only overdue milestones")
	getMilestonesCmd.Flags().Bool("closed", false, "Show only closed milestones")
	milestonesCmd.AddCommand(createMilestonesCmd)
	createMilestonesCmd.Flags().String("title", "", "Title of the milestone")
	createMilestonesCmd.Flags().String("description", "", "Description of the milestone")
	createMilestonesCmd.Flags().String("due-on", "", "Due date of the milestone (YYYY-MM-DD)")
	milestonesCmd.AddCommand(updateMilestonesCmd)
	updateMilestonesCmd.Flags().String("description", "", "Description of the milestone")
	updateMilestonesCmd.Flags().String("due-on", "", "Due date of the milestone (YYYY-MM-DD)")
}
