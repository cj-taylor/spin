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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	NAME  = "project"
	EMAIL = "appowner@spinnaker-test.net"
)

func TestProjectSave_basic(t *testing.T) {
	ts := testGateProjectSaveSuccess()
	defer ts.Close()
	currentCmd := NewSaveCmd(projectOptions{})
	rootCmd := getRootCmdForTest()
	projCmd := NewProjectCmd(os.Stdout)
	projCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(projCmd)

	args := []string{
		"project", "save",
		"--gate-endpoint=" + ts.URL,
		"--project-name", NAME,
		"--owner-email", EMAIL,
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Command failed with: %s", err)
	}
}

func TestProjectSave_fail(t *testing.T) {
	ts := GateServerFail()
	defer ts.Close()
	args := []string{
		"project", "save",
		"--project-name", NAME,
		"--owner-email", EMAIL,
		"--cloud-providers", "gce,kubernetes",
		"--gate-endpoint=" + ts.URL,
	}
	currentCmd := NewSaveCmd(projectOptions{})
	rootCmd := getRootCmdForTest()
	projCmd := NewProjectCmd(os.Stdout)
	projCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(projCmd)

	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	if err == nil {
		t.Fatalf("Command failed with: %s", err)
	}
}

func TestProjectSave_flags(t *testing.T) {
	ts := testGateProjectSaveSuccess()
	defer ts.Close()

	args := []string{
		"project", "save",
		"--gate-endpoint=" + ts.URL,
	}
	currentCmd := NewSaveCmd(projectOptions{})
	rootCmd := getRootCmdForTest()
	projCmd := NewProjectCmd(os.Stdout)
	projCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(projCmd)

	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	if err == nil {
		t.Fatalf("Command failed with: %s", err)
	}
}

func TestProjectSave_missingname(t *testing.T) {
	ts := testGateProjectSaveSuccess()
	defer ts.Close()

	args := []string{
		"project", "save",
		"--owner-email", EMAIL,
		"--gate-endpoint=" + ts.URL,
	}
	currentCmd := NewSaveCmd(projectOptions{})
	rootCmd := getRootCmdForTest()
	projCmd := NewProjectCmd(os.Stdout)
	projCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(projCmd)

	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	if err == nil {
		t.Fatalf("Command failed with: %s", err)
	}
}

func TestProjectSave_missingemail(t *testing.T) {
	ts := testGateProjectSaveSuccess()
	defer ts.Close()

	args := []string{
		"project", "save",
		"--project-name", NAME,
		"--cloud-providers", "gce,kubernetes",
		"--gate-endpoint", ts.URL,
	}
	currentCmd := NewSaveCmd(projectOptions{})
	rootCmd := getRootCmdForTest()
	projCmd := NewProjectCmd(os.Stdout)
	projCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(projCmd)

	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	if err == nil {
		t.Fatalf("Command failed with: %s", err)
	}
}

func TestProjectSave_filebasic(t *testing.T) {
	ts := testGateProjectSaveSuccess()
	defer ts.Close()

	tempFile := tempAppFile(testProjJsonStr)
	if tempFile == nil {
		t.Fatal("Could not create temp project file.")
	}
	defer os.Remove(tempFile.Name())

	args := []string{
		"project", "save",
		"--file", tempFile.Name(),
		"--gate-endpoint", ts.URL,
	}

	currentCmd := NewSaveCmd(projectOptions{})
	rootCmd := getRootCmdForTest()
	projCmd := NewProjectCmd(os.Stdout)
	projCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(projCmd)

	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Command failed with: %s", err)
	}
}

func TestProjectSave_stdinbasic(t *testing.T) {
	ts := testGateProjectSaveSuccess()
	defer ts.Close()

	tempFile := tempAppFile(testProjJsonStr)
	if tempFile == nil {
		t.Fatal("Could not create temp app file.")
	}
	defer os.Remove(tempFile.Name())

	// Prepare Stdin for test reading.
	tempFile.Seek(0, 0)
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin = tempFile

	args := []string{
		"project", "save",
		"--gate-endpoint", ts.URL,
	}

	currentCmd := NewSaveCmd(projectOptions{})
	rootCmd := getRootCmdForTest()
	projCmd := NewProjectCmd(os.Stdout)
	projCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(projCmd)

	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Command failed with: %s", err)
	}
}
// testGatePipelineExecuteSuccess spins up a local http server that we will configure the GateClient
// to direct requests to. Responds with successful responses to pipeline execute API calls.
func testGateProjectSaveSuccess() *httptest.Server {
	mux := http.NewServeMux()
	mux.Handle("/tasks", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]string{
			"ref": "/tasks/id",
		}
		b, _ := json.Marshal(&payload)
		fmt.Fprintln(w, string(b))
	}))
	mux.Handle("/tasks/id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]string{
			"status": "SUCCEEDED",
		}
		b, _ := json.Marshal(&payload)
		fmt.Fprintln(w, string(b))
	}))
	return httptest.NewServer(mux)
}

func tempAppFile(appContent string) *os.File {
	tempFile, _ := ioutil.TempFile("" /* /tmp dir. */, "app-spec")
	bytes, err := tempFile.Write([]byte(appContent))
	if err != nil || bytes == 0 {
		fmt.Println("Could not write temp file.")
		return nil
	}
	return tempFile
}

const testProjJsonStr = `
{
   "email" : "someone@example.com",
   "name" : "sampleproject"
}

`
