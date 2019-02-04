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
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/spinnaker/spin/util"

	"github.com/spf13/cobra"
)

const (
	PROJECT = "proj"
)

func getRootCmdForTest() *cobra.Command {
	rootCmd := &cobra.Command{}
	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.spin/config)")
	rootCmd.PersistentFlags().String("gate-endpoint", "", "Gate (API server) endpoint. Default http://localhost:8084")
	rootCmd.PersistentFlags().Bool("insecure", false, "Ignore Certificate Errors")
	rootCmd.PersistentFlags().Bool("quiet", false, "Squelch non-essential output")
	rootCmd.PersistentFlags().Bool("no-color", false, "Disable color")
	rootCmd.PersistentFlags().String("output", "", "Configure output formatting")
	util.InitUI(false, false, "")
	return rootCmd
}

func TestProjectGet_basic(t *testing.T) {
	ts := testGateProjectGetSuccess()
	defer ts.Close()
	currentCmd := NewGetCmd(projectOptions{})
	rootCmd := getRootCmdForTest()
	projCmd := NewProjectCmd(os.Stdout)
	projCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(projCmd)

	args := []string{"project", "get", PROJECT, "--gate-endpoint=" + ts.URL}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Command failed with: %s", err)
	}
}

func TestProjectGet_flags(t *testing.T) {
	ts := testGateProjectGetSuccess()
	defer ts.Close()
	currentCmd := NewGetCmd(projectOptions{})
	rootCmd := getRootCmdForTest()
	projCmd := NewProjectCmd(os.Stdout)
	projCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(projCmd)
	args := []string{"roject", "get", "--gate-endpoint", ts.URL} // Missing positional arg.
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	if err == nil { // Success is actually failure here, flags are malformed.
		t.Fatalf("Command failed with: %s", err)
	}
}

func TestProjectGet_malformed(t *testing.T) {
	ts := testGateProjectGetMalformed()
	defer ts.Close()

	currentCmd := NewGetCmd(projectOptions{})
	rootCmd := getRootCmdForTest()
	projCmd := NewProjectCmd(os.Stdout)
	projCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(projCmd)

	args := []string{"project", "get", PROJECT, "--gate-endpoint=" + ts.URL}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	if err == nil { // Success is actually failure here, return payload is malformed.
		t.Fatalf("Command failed with: %d", err)
	}
}

func TestProjectGet_fail(t *testing.T) {
	ts := GateServerFail()
	defer ts.Close()

	currentCmd := NewGetCmd(projectOptions{})
	rootCmd := getRootCmdForTest()
	projCmd := NewProjectCmd(os.Stdout)
	projCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(projCmd)

	args := []string{"project", "get", PROJECT, "--gate-endpoint=" + ts.URL}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	if err == nil { // Success is actually failure here, return payload is malformed.
		t.Fatalf("Command failed with: %d", err)
	}
}

// testGateProjectGetSuccess spins up a local http server that we will configure the GateClient
// to direct requests to. Responds with a 200 and a well-formed pipeline list.
func testGateProjectGetSuccess() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, strings.TrimSpace(projectJson))
	}))
}

// testGateProjectGetMalformed returns a malformed list of pipeline configs.
func testGateProjectGetMalformed() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, strings.TrimSpace(malformedProjectGetJson))
	}))
}

// GateServerFail spins up a local http server that we will configure the GateClient
// to direct requests to. Responds with a 500 InternalServerError.
func GateServerFail() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(jacobkiefer): Mock more robust errors once implemented upstream.
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
}


const malformedProjectGetJson = `
  "config": {
    "applications": null,
    "clusters": null,
    "pipelineConfigs": null
  },
  "createTs": 1549251430546,
  "email": "proj",
  "id": "9a347a80-d824-4c37-a937-feb73d0769c8",
  "lastModifiedBy": "anonymous",
  "name": "proj",
  "updateTs": 1549251430000
}
`

const projectJson = `
{
  "config": {
    "applications": null,
    "clusters": null,
    "pipelineConfigs": null
  },
  "createTs": 1549251430546,
  "email": "proj",
  "id": "9a347a80-d824-4c37-a937-feb73d0769c8",
  "lastModifiedBy": "anonymous",
  "name": "proj",
  "updateTs": 1549251430000
}
`
