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

package project

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spinnaker/spin/cmd/gateclient"
	"github.com/spinnaker/spin/util"
)

type ListOptions struct {
	*projectOptions
}

var (
	listProjectShort   = "List the all projects"
	listProjectLong    = "List the all projects"
	listProjectExample = "usage: spin project list [options]"
)

func NewListCmd(projOptions projectOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   listProjectShort,
		Long:    listProjectLong,
		Example: listProjectExample,
		RunE:    listProject,
	}
	return cmd
}

func listProject(cmd *cobra.Command, args []string) error {
	gateClient, err := gateclient.NewGateClient(cmd.InheritedFlags())
	if err != nil {
		return err
	}
	projectList, resp, err := gateClient.ProjectControllerApi.AllUsingGET3(gateClient.Context)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error getting projects, status code: %d\n", resp.StatusCode)
	}

	util.UI.JsonOutput(projectList, util.UI.OutputFormat)
	return nil
}
