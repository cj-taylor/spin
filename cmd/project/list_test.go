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
)

func TestProjectList_basic(t *testing.T) {
	ts := testProjectListSuccess()
	defer ts.Close()

	currentCmd := NewListCmd(projectOptions{})
	rootCmd := getRootCmdForTest()
	projCmd := NewProjectCmd(os.Stdout)
	projCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(projCmd)

	args := []string{"project", "list", "--gate-endpoint=" + ts.URL}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Command failed with: %s", err)
	}
}

func TestProjectList_malformed(t *testing.T) {
	ts := testGateProjectListMalformed()
	defer ts.Close()

	currentCmd := NewListCmd(projectOptions{})
	rootCmd := getRootCmdForTest()
	projCmd := NewProjectCmd(os.Stdout)
	projCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(projCmd)

	args := []string{"project", "list", "--gate-endpoint=" + ts.URL}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	if err == nil {
		t.Fatalf("Command failed with: %s", err)
	}
}

func TestProjectList_fail(t *testing.T) {
	ts := GateServerFail()
	defer ts.Close()

	currentCmd := NewListCmd(projectOptions{})
	rootCmd := getRootCmdForTest()
	projCmd := NewProjectCmd(os.Stdout)
	projCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(projCmd)

	args := []string{"project", "list", "--gate-endpoint=" + ts.URL}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Command failed with: %s", err)
	}
}

// testProjectListSuccess spins up a local http server that we will configure the GateClient
// to direct requests to. Responds with a 200 and a well-formed project list.
func testProjectListSuccess() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, strings.TrimSpace(projectListJson))
	}))
}

// testGateProjectListMalformed returns a malformed list of projects.
func testGateProjectListMalformed() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, strings.TrimSpace(malformedProjectListJson))
	}))
}

const malformedProjectListJson = `
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
]
`

const projectListJson = `
[
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
]
`
