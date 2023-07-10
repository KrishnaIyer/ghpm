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

	"github.com/spf13/cobra"
	"krishnaiyer.dev/golang/datasink/pkg/device"
	conf "krishnaiyer.dev/golang/dry/pkg/config"
	logger "krishnaiyer.dev/golang/dry/pkg/logger"
)

const (
	defaultBufferSize = 64
)

// Config contains the configuration.
type Config struct {
}

var (
	config  = &Config{}
	manager *conf.Manager
	baseCtx = context.Background()
	devices = make(map[string]device.Device)

	// Root is the root of the commands.
	Root = &cobra.Command{
		Use:           "datasink",
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         "datasink is tool that acts as acts as a server with multiple protocols (ex: mqtt, websocket) for incoming traffic and writes to a time series database",
		Long:          `datasink is tool that acts as acts as a server with multiple protocols (ex: mqtt, websocket) for incoming traffic and writes to a time series database. More documentation at https://krishnaiyer.dev/golang/datasink`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := manager.ReadFromFile(cmd.Flags())
			if err != nil {
				panic(err)
			}
			err = manager.Unmarshal(&config)
			if err != nil {
				panic(err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(baseCtx)
			defer cancel()

			l, err := logger.New(ctx, false)
			if err != nil {
				panic(err)
			}
			ctx = logger.NewContextWithLogger(ctx, l)
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
	manager = conf.New("config")
	manager.InitFlags(*config)
	// This line is needed to persist the config file to subcommands.
	manager.AddConfigFlag(manager.Flags())
	Root.PersistentFlags().AddFlagSet(manager.Flags())
	Root.AddCommand(VersionCommand(Root))
}
