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
	"context"
	"log"

	"github.com/google/go-github/v54/github"
	"github.com/spf13/cobra"
	conf "krishnaiyer.dev/golang/dry/pkg/config"
	logger "krishnaiyer.dev/golang/dry/pkg/logger"
)

// Config contains the configuration.
type Config struct {
	Token        string `name:"token" description:"The GitHub token"`
	Username     string `name:"username" description:"The GitHub user or organization name"`
	Repositories []struct {
		Name     string `name:"name" description:"The GitHub repository name"`
		Username string `name:"username" description:"The GitHub user or organization name. Overrides global username"`
	} `name:"repositories" description:"The GitHub repositories"`
}

const (
	baseURL = "https://api.github.com"
)

var (
	config  = &Config{}
	manager *conf.Manager
	client  *github.Client
	ctx     context.Context

	// Root is the root of the commands.
	Root = &cobra.Command{
		Use:           "ghpm",
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         "ghpm is a tool to manage github repositories",
		Long:          `ghpm is a tool to manage github repositories.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := manager.ReadFromFile(cmd.Flags())
			if err != nil {
				panic(err)
			}
			err = manager.Unmarshal(&config)
			if err != nil {
				panic(err)
			}
			ctx = context.Background()
			l, err := logger.New(ctx, false)
			if err != nil {
				panic(err)
			}
			ctx = logger.NewContextWithLogger(ctx, l)
			client = github.NewTokenClient(ctx, config.Token)
			return nil
		},
	}
)

// Execute ...
func Execute() {
	if err := Root.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}

func init() {
	manager = conf.New("ghpm")
	manager.InitFlags(*config)
	// This line is needed to persist the config file to subcommands.
	manager.AddConfigFlag(manager.Flags())
	Root.PersistentFlags().AddFlagSet(manager.Flags())
	Root.AddCommand(VersionCommand(Root))
}
