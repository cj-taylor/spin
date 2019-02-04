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
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spinnaker/spin/cmd/gateclient"
	"github.com/spinnaker/spin/util"
)

type SaveOptions struct {
	*projectOptions
	projectFile    string
	projectName    string
	ownerEmail     string
	cloudProviders *[]string
}

var (
	saveProjectShort = "Save the provided application"
	saveProjectLong  = "Save the specified application"
)

func NewSaveCmd(projOptions projectOptions) *cobra.Command {
	options := SaveOptions{
		projectOptions: &projOptions,
	}
	cmd := &cobra.Command{
		Use:   "save",
		Short: saveProjectShort,
		Long:  saveProjectLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			return saveProject(cmd, options)
		},
	}
	cmd.PersistentFlags().StringVarP(&options.projectName, "project-name", "", "", "name of the application")
	cmd.PersistentFlags().StringVarP(&options.ownerEmail, "owner-email", "", "", "email of the application owner")

	return cmd
}

func saveProject(cmd *cobra.Command, options SaveOptions) error {
	gateClient, err := gateclient.NewGateClient(cmd.InheritedFlags())
	if err != nil {
		return err
	}

	initialProject, err := util.ParseJsonFromFileOrStdin(options.projectFile, true)
	if err != nil {
		return fmt.Errorf("Could not parse supplied project: %v.\n", err)
	}

	var app map[string]interface{}
	if initialProject != nil && len(initialProject) > 0 {
		app = initialProject
		if options.projectName != "" {
			util.UI.Warn("Overriding project name with explicit flag values.\n")
			app["name"] = options.projectName
		}
		if options.ownerEmail != "" {
			util.UI.Warn("Overriding project owner email with explicit flag values.\n")
			app["email"] = options.ownerEmail
		}
	} else {
		if options.projectName == "" || options.ownerEmail == "" {
			return errors.New("Required project parameter missing, exiting...")
		}
		app = map[string]interface{}{
			"name":           options.projectName,
			"email":          options.ownerEmail,
		}
	}

	createAppTask := map[string]interface{}{
		"job":         []interface{}{map[string]interface{}{"type": "createApplication", "application": app}},
		"application": app["name"],
		"description": fmt.Sprintf("Create Application: %s", app["name"]),
	}

	ref, _, err := gateClient.TaskControllerApi.TaskUsingPOST1(gateClient.Context, createAppTask)
	if err != nil {
		return err
	}

	toks := strings.Split(ref["ref"].(string), "/")
	id := toks[len(toks)-1]

	task, resp, err := gateClient.TaskControllerApi.GetTaskUsingGET1(gateClient.Context, id)
	attempts := 0
	for (task == nil || !taskCompleted(task)) && attempts < 5 {
		toks := strings.Split(ref["ref"].(string), "/")
		id := toks[len(toks)-1]

		task, resp, err = gateClient.TaskControllerApi.GetTaskUsingGET1(gateClient.Context, id)
		attempts += 1
		time.Sleep(time.Duration(attempts*attempts) * time.Second)
	}

	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("Encountered an error saving application, status code: %d\n", resp.StatusCode)
	}
	if !taskSucceeded(task) {
		return fmt.Errorf("Encountered an error saving application, task output was: %v\n", task)
	}

	util.UI.Info(util.Colorize().Color(fmt.Sprintf("[reset][bold][green]Application save succeeded")))
	return nil
}

// TODO(jacobkiefer): Consider generalizing if we need these functions elsewhere.
func taskCompleted(task map[string]interface{}) bool {
	taskStatus, exists := task["status"]
	if !exists {
		return false
	}

	COMPLETED := [...]string{"SUCCEEDED", "STOPPED", "SKIPPED", "TERMINAL", "FAILED_CONTINUE"}
	for _, status := range COMPLETED {
		if taskStatus == status {
			return true
		}
	}
	return false
}

func taskSucceeded(task map[string]interface{}) bool {
	taskStatus, exists := task["status"]
	if !exists {
		return false
	}

	SUCCESSFUL := [...]string{"SUCCEEDED", "STOPPED", "SKIPPED"}
	for _, status := range SUCCESSFUL {
		if taskStatus == status {
			return true
		}
	}
	return false
}
