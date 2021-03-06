// Copyright (c) 2018, Google, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package application

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/spinnaker/spin/util"

	"github.com/spf13/cobra"
	"github.com/spinnaker/spin/cmd/gateclient"
)

type GetOptions struct {
	*applicationOptions
	output string
}

var (
	getApplicationShort   = "Get the specified application"
	getApplicationLong    = "Get the specified application"
	getApplicationExample = "usage: spin application get [options] application-name"
)

func NewGetCmd(appOptions applicationOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"get"},
		Short:   getApplicationShort,
		Long:    getApplicationLong,
		Example: getApplicationExample,
		RunE:    getApplication,
	}

	return cmd
}

func getApplication(cmd *cobra.Command, args []string) error {
	gateClient, err := gateclient.NewGateClient(cmd.InheritedFlags())
	if err != nil {
		return err
	}
	if len(args) == 0 || args[0] == "" {
		return errors.New("application name required")
	}
	applicationName := args[0]
	app, resp, err := gateClient.ApplicationControllerApi.GetApplicationUsingGET(gateClient.Context, applicationName, map[string]interface{}{})
	if resp != nil {
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Application '%s' not found\n", applicationName)
		} else if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Encountered an error getting application, status code: %d\n", resp.StatusCode)
		}
	}

	if err != nil {
		return err
	}

	util.UI.JsonOutput(app, util.UI.OutputFormat)
	return nil
}
