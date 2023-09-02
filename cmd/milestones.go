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
	getMilestonesCmd = &cobra.Command{
		Use:   "get",
		Short: "Get milestones",
		Long:  `Get milestones.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			milestones, _, err := client.Issues.ListMilestones(
				ctx,
				config.Username,
				config.Repositories[0].Name,
				&github.MilestoneListOptions{
					State:     "all",
					Sort:      "due_on",
					Direction: "asc",
				})
			if err != nil {
				return err
			}

			fmt.Println("##################################################")
			fmt.Println("\t\t Milestones")
			fmt.Println("##################################################")
			fmt.Printf("Repo: %s\n", config.Repositories[0].Name)
			fmt.Println("--------------------------------------------------")
			for _, milestone := range milestones {
				var (
					dueOn    string
					closedAt string
				)

				if milestone.DueOn != nil {
					dueOn = milestone.DueOn.Format("02 Jan 06")
				}
				if milestone.ClosedAt != nil {
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
				milestoneJSON, err := json.Marshal(m)
				if err != nil {
					return err
				}
				fmt.Println(string(milestoneJSON))
			}
			return nil
		},
	}
)

func init() {
	Root.AddCommand(milestonesCmd)
	milestonesCmd.AddCommand(getMilestonesCmd)
}
