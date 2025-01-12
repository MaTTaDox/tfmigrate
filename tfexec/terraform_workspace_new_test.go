package tfexec

import (
	"context"
	"testing"
)

func TestTerraformCLIWorkspaceNew(t *testing.T) {
	cases := []struct {
		desc         string
		mockCommands []*mockCommand
		workspace    string
		dir          string
		opts         []string
		ok           bool
	}{
		{
			desc: "no workspace, no dir, no opts",
			mockCommands: []*mockCommand{
				{
					args:     []string{"terraform", "workspace", "new"},
					exitCode: 1,
				},
			},
			ok: false,
		},
		{
			desc: "with workspace",
			mockCommands: []*mockCommand{
				{
					args:     []string{"terraform", "workspace", "new", "foo"},
					exitCode: 0,
				},
			},
			workspace: "foo",
			ok:        true,
		},
		{
			desc: "with workspace and dir",
			mockCommands: []*mockCommand{
				{
					args:     []string{"terraform", "workspace", "new", "foo", "bar"},
					exitCode: 0,
				},
			},
			workspace: "foo",
			dir:       "bar",
			ok:        true,
		},
		{
			desc: "with workspace and dir and opts",
			mockCommands: []*mockCommand{
				{
					args:     []string{"terraform", "workspace", "new", "-lock=true", "-lock-timeout=0s", "foo", "bar"},
					exitCode: 0,
				},
			},
			workspace: "foo",
			dir:       "bar",
			opts:      []string{"-lock=true", "-lock-timeout=0s"},
			ok:        true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			e := NewMockExecutor(tc.mockCommands)
			terraformCLI := NewTerraformCLI(e)
			err := terraformCLI.WorkspaceNew(context.Background(), tc.workspace, tc.dir, tc.opts...)
			if tc.ok && err != nil {
				t.Fatalf("unexpected err: %s", err)
			}
			if !tc.ok && err == nil {
				t.Fatal("expected to return an error, but no error")
			}
		})
	}
}

func TestAccTerraformCLIWorkspaceNew(t *testing.T) {
	SkipUnlessAcceptanceTestEnabled(t)

	source := `resource "null_resource" "foo" {}`
	e := SetupTestAcc(t, source)
	terraformCLI := NewTerraformCLI(e)

	err := terraformCLI.Init(context.Background(), "", "-input=false", "-no-color")
	if err != nil {
		t.Fatalf("failed to run terraform init: %s", err)
	}

	err = terraformCLI.WorkspaceNew(context.Background(), "myworkspace", "")
	if err != nil {
		t.Fatalf("failed to create a new workspace: %s", err)
	}

	got, err := terraformCLI.WorkspaceShow(context.Background())
	if err != nil {
		t.Fatalf("failed to run terraform workspace show: %s", err)
	}

	if got == "" {
		t.Error("The current workspace doesn't match the workspace that was just created")
	}
}
