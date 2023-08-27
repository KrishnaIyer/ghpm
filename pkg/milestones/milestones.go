// Copyright © 2023 Krishna Iyer Easwaran
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package milestones

import "krishnaiyer.dev/golang/ghpm/pkg/client"

type Config struct {
	Username string `name:"username" description:"The GitHub user or organization name"`
	Repo     string `name:"repo" description:"The GitHub repository name"`
}

// Manager is a GitHub manager.
type Manager struct {
	client *client.Client
}

// New creates a new manager.
func New(client *client.Client) *Manager {
	return &Manager{
		client: client,
	}
}

// ListMilestones lists all milestones.
func ListMilestones() {
}

// ListPastMilestones lists all past milestones.
func ListPastMilestones() {
}
