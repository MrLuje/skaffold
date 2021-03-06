/*
Copyright 2020 The Skaffold Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package deploy

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/runner/runcontext"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/latest"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/util"
	"github.com/GoogleContainerTools/skaffold/testutil"
)

func TestKpt_Deploy(t *testing.T) {
	tests := []struct {
		description string
		expected    []string
		shouldErr   bool
	}{
		{
			description: "nil",
		},
	}
	for _, test := range tests {
		testutil.Run(t, test.description, func(t *testutil.T) {
			k := NewKptDeployer(&runcontext.RunContext{}, nil)
			res, err := k.Deploy(context.Background(), nil, nil)
			t.CheckErrorAndDeepEqual(test.shouldErr, err, test.expected, res)
		})
	}
}

func TestKpt_Dependencies(t *testing.T) {
	tests := []struct {
		description string
		expected    []string
		shouldErr   bool
	}{
		{
			description: "nil",
		},
	}
	for _, test := range tests {
		testutil.Run(t, test.description, func(t *testutil.T) {
			k := NewKptDeployer(&runcontext.RunContext{}, nil)
			res, err := k.Dependencies()
			t.CheckErrorAndDeepEqual(test.shouldErr, err, test.expected, res)
		})
	}
}

func TestKpt_Cleanup(t *testing.T) {
	tests := []struct {
		description string
		shouldErr   bool
	}{
		{
			description: "nil",
		},
	}
	for _, test := range tests {
		testutil.Run(t, test.description, func(t *testutil.T) {
			k := NewKptDeployer(&runcontext.RunContext{}, nil)
			err := k.Cleanup(context.Background(), nil)
			t.CheckError(test.shouldErr, err)
		})
	}
}

func TestKpt_Render(t *testing.T) {
	tests := []struct {
		description string
		shouldErr   bool
	}{
		{},
	}
	for _, test := range tests {
		testutil.Run(t, test.description, func(t *testutil.T) {
			k := NewKptDeployer(&runcontext.RunContext{}, nil)
			err := k.Render(context.Background(), nil, nil, false, "")
			t.CheckError(test.shouldErr, err)
		})
	}
}

func TestKpt_GetApplyDir(t *testing.T) {
	tests := []struct {
		description string
		applyDir    string
		expected    string
		commands    util.Command
		shouldErr   bool
	}{
		{
			description: "specified an invalid applyDir",
			applyDir:    "invalid_path",
			shouldErr:   true,
		},
		{
			description: "specified a valid applyDir",
			applyDir:    "valid_path",
			expected:    "valid_path",
		},
		{
			description: "unspecified applyDir",
			expected:    ".kpt-hydrated",
			commands:    testutil.CmdRun("kpt live init .kpt-hydrated"),
		},
		{
			description: "existing template resource in .kpt-hydrated",
			expected:    ".kpt-hydrated",
		},
	}
	for _, test := range tests {
		testutil.Run(t, test.description, func(t *testutil.T) {
			t.Override(&util.DefaultExecCommand, test.commands)
			tmpDir := t.NewTempDir().Chdir()

			if test.applyDir == test.expected {
				os.Mkdir(test.applyDir, 0755)
			}

			if test.description == "existing template resource in .kpt-hydrated" {
				tmpDir.Touch(".kpt-hydrated/inventory-template.yaml")
			}

			k := NewKptDeployer(&runcontext.RunContext{
				WorkingDir: ".",
				Cfg: latest.Pipeline{
					Deploy: latest.DeployConfig{
						DeployType: latest.DeployType{
							KptDeploy: &latest.KptDeploy{
								ApplyDir: test.applyDir,
							},
						},
					},
				},
			}, nil)

			applyDir, err := k.getApplyDir(context.Background())

			t.CheckErrorAndDeepEqual(test.shouldErr, err, test.expected, applyDir)
		})
	}
}

func TestKpt_KptCommandArgs(t *testing.T) {
	tests := []struct {
		description string
		dir         string
		commands    []string
		flags       []string
		globalFlags []string
		expected    []string
	}{
		{
			description: "empty",
		},
		{
			description: "all inputs have len >0",
			dir:         "test",
			commands:    []string{"live", "apply"},
			flags:       []string{"--fn-path", "kpt-func.yaml"},
			globalFlags: []string{"-h"},
			expected:    strings.Split("live apply test --fn-path kpt-func.yaml -h", " "),
		},
		{
			description: "empty dir",
			commands:    []string{"live", "apply"},
			flags:       []string{"--fn-path", "kpt-func.yaml"},
			globalFlags: []string{"-h"},
			expected:    strings.Split("live apply --fn-path kpt-func.yaml -h", " "),
		},
		{
			description: "empty commands",
			dir:         "test",
			flags:       []string{"--fn-path", "kpt-func.yaml"},
			globalFlags: []string{"-h"},
			expected:    strings.Split("test --fn-path kpt-func.yaml -h", " "),
		},
		{
			description: "empty flags",
			dir:         "test",
			commands:    []string{"live", "apply"},
			globalFlags: []string{"-h"},
			expected:    strings.Split("live apply test -h", " "),
		},
		{
			description: "empty globalFlags",
			dir:         "test",
			commands:    []string{"live", "apply"},
			flags:       []string{"--fn-path", "kpt-func.yaml"},
			expected:    strings.Split("live apply test --fn-path kpt-func.yaml", " "),
		},
	}
	for _, test := range tests {
		testutil.Run(t, test.description, func(t *testutil.T) {
			res := kptCommandArgs(test.dir, test.commands, test.flags, test.globalFlags)
			t.CheckDeepEqual(test.expected, res)
		})
	}
}
