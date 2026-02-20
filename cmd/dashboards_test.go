// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.

package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/datadog-labs/pup/pkg/client"
	"github.com/datadog-labs/pup/pkg/config"
)

func TestDashboardsCmd(t *testing.T) {
	if dashboardsCmd == nil {
		t.Fatal("dashboardsCmd is nil")
	}

	if dashboardsCmd.Use != "dashboards" {
		t.Errorf("Use = %s, want dashboards", dashboardsCmd.Use)
	}

	if dashboardsCmd.Short == "" {
		t.Error("Short description is empty")
	}
}

func TestDashboardsCmd_Subcommands(t *testing.T) {
	expectedCommands := []string{"list", "get", "delete", "create", "update"}

	commands := dashboardsCmd.Commands()

	commandMap := make(map[string]bool)
	for _, cmd := range commands {
		commandMap[cmd.Name()] = true
	}

	for _, expected := range expectedCommands {
		if !commandMap[expected] {
			t.Errorf("Missing subcommand: %s", expected)
		}
	}
}

// Helper function to setup dashboards test client
func setupDashboardsTestClient(t *testing.T) func() {
	t.Helper()

	origClient := ddClient
	origCfg := cfg
	origFactory := clientFactory

	cfg = &config.Config{
		Site:        "datadoghq.com",
		APIKey:      "test-api-key-12345678",
		AppKey:      "test-app-key-12345678",
		AutoApprove: false,
	}

	clientFactory = func(c *config.Config) (*client.Client, error) {
		return nil, fmt.Errorf("mock client: no real API connection in tests")
	}

	ddClient = nil

	return func() {
		ddClient = origClient
		cfg = origCfg
		clientFactory = origFactory
	}
}

func TestRunDashboardsList(t *testing.T) {
	cleanup := setupDashboardsTestClient(t)
	defer cleanup()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "fails on client creation",
			wantErr: true, // Mock client error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			outputWriter = &buf
			defer func() { outputWriter = os.Stdout }()

			err := runDashboardsList(dashboardsListCmd, []string{})

			if (err != nil) != tt.wantErr {
				t.Errorf("runDashboardsList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRunDashboardsGet(t *testing.T) {
	cleanup := setupDashboardsTestClient(t)
	defer cleanup()

	tests := []struct {
		name        string
		dashboardID string
		wantErr     bool
	}{
		{
			name:        "with valid dashboard ID",
			dashboardID: "abc-123-xyz",
			wantErr:     true, // Mock client error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			outputWriter = &buf
			defer func() { outputWriter = os.Stdout }()

			err := runDashboardsGet(dashboardsGetCmd, []string{tt.dashboardID})

			if (err != nil) != tt.wantErr {
				t.Errorf("runDashboardsGet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRunDashboardsDelete_AutoApprove(t *testing.T) {
	cleanup := setupDashboardsTestClient(t)
	defer cleanup()

	// Set auto-approve
	cfg.AutoApprove = true

	tests := []struct {
		name        string
		dashboardID string
		wantErr     bool
	}{
		{
			name:        "with auto-approve",
			dashboardID: "abc-123-xyz",
			wantErr:     true, // Mock client error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			outputWriter = &buf
			defer func() { outputWriter = os.Stdout }()

			err := runDashboardsDelete(dashboardsDeleteCmd, []string{tt.dashboardID})

			if (err != nil) != tt.wantErr {
				t.Errorf("runDashboardsDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRunDashboardsDelete_WithConfirmation(t *testing.T) {
	cleanup := setupDashboardsTestClient(t)
	defer cleanup()

	// Disable auto-approve
	cfg.AutoApprove = false

	tests := []struct {
		name        string
		dashboardID string
		input       string
		wantErr     bool
	}{
		{
			name:        "fails on client creation (mock)",
			dashboardID: "abc-123-xyz",
			input:       "n\n",
			wantErr:     true, // getClient() called before confirmation
		},
		{
			name:        "fails on client creation with yes (mock)",
			dashboardID: "abc-123-xyz",
			input:       "yes\n",
			wantErr:     true, // getClient() called before confirmation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			outputWriter = &buf
			defer func() { outputWriter = os.Stdout }()

			// Simulate input
			inputReader = strings.NewReader(tt.input)
			defer func() { inputReader = os.Stdin }()

			err := runDashboardsDelete(dashboardsDeleteCmd, []string{tt.dashboardID})

			if (err != nil) != tt.wantErr {
				t.Errorf("runDashboardsDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDashboardsCreateCmd(t *testing.T) {
	if dashboardsCreateCmd == nil {
		t.Fatal("dashboardsCreateCmd is nil")
	}

	if dashboardsCreateCmd.Use != "create" {
		t.Errorf("Use = %s, want create", dashboardsCreateCmd.Use)
	}

	if dashboardsCreateCmd.Short == "" {
		t.Error("Short description is empty")
	}

	if dashboardsCreateCmd.RunE == nil {
		t.Error("RunE is nil")
	}

	flags := dashboardsCreateCmd.Flags()
	if flags.Lookup("body") == nil {
		t.Error("Missing --body flag")
	}
}

func TestDashboardsCreateCmd_BodyRequired(t *testing.T) {
	if dashboardsCreateCmd.Flags().Lookup("body") == nil {
		t.Fatal("--body flag not found")
	}

	if err := dashboardsCreateCmd.ValidateRequiredFlags(); err == nil {
		t.Error("expected --body to be required")
	}
}

func TestDashboardsUpdateCmd(t *testing.T) {
	if dashboardsUpdateCmd == nil {
		t.Fatal("dashboardsUpdateCmd is nil")
	}

	if dashboardsUpdateCmd.Use != "update [dashboard-id]" {
		t.Errorf("Use = %s, want 'update [dashboard-id]'", dashboardsUpdateCmd.Use)
	}

	if dashboardsUpdateCmd.Short == "" {
		t.Error("Short description is empty")
	}

	if dashboardsUpdateCmd.RunE == nil {
		t.Error("RunE is nil")
	}

	if dashboardsUpdateCmd.Args == nil {
		t.Error("Args validator is nil")
	}

	flags := dashboardsUpdateCmd.Flags()
	if flags.Lookup("body") == nil {
		t.Error("Missing --body flag")
	}
}

func TestDashboardsUpdateCmd_BodyRequired(t *testing.T) {
	if dashboardsUpdateCmd.Flags().Lookup("body") == nil {
		t.Fatal("--body flag not found")
	}

	if err := dashboardsUpdateCmd.ValidateRequiredFlags(); err == nil {
		t.Error("expected --body to be required")
	}
}

func TestRunDashboardsCreate(t *testing.T) {
	cleanup := setupDashboardsTestClient(t)
	defer cleanup()

	tmpFile := t.TempDir() + "/dashboard.json"
	content := []byte(`{"title":"Test Dashboard","layout_type":"ordered","widgets":[]}`)
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		body    string
		wantErr bool
	}{
		{
			name:    "fails on client creation",
			body:    "@" + tmpFile,
			wantErr: true, // Mock client error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			outputWriter = &buf
			defer func() { outputWriter = os.Stdout }()

			if err := dashboardsCreateCmd.Flags().Set("body", tt.body); err != nil {
				t.Fatal(err)
			}

			err := runDashboardsCreate(dashboardsCreateCmd, []string{})

			if (err != nil) != tt.wantErr {
				t.Errorf("runDashboardsCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRunDashboardsUpdate(t *testing.T) {
	cleanup := setupDashboardsTestClient(t)
	defer cleanup()

	tmpFile := t.TempDir() + "/dashboard.json"
	content := []byte(`{"title":"Updated Dashboard","layout_type":"ordered","widgets":[]}`)
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		dashboardID string
		body        string
		wantErr     bool
	}{
		{
			name:        "fails on client creation",
			dashboardID: "abc-def-123",
			body:        "@" + tmpFile,
			wantErr:     true, // Mock client error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			outputWriter = &buf
			defer func() { outputWriter = os.Stdout }()

			if err := dashboardsUpdateCmd.Flags().Set("body", tt.body); err != nil {
				t.Fatal(err)
			}

			err := runDashboardsUpdate(dashboardsUpdateCmd, []string{tt.dashboardID})

			if (err != nil) != tt.wantErr {
				t.Errorf("runDashboardsUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDashboardsCmd_ParentChild(t *testing.T) {
	commands := dashboardsCmd.Commands()

	for _, cmd := range commands {
		if cmd.Parent() != dashboardsCmd {
			t.Errorf("Command %s parent is not dashboardsCmd", cmd.Use)
		}
	}
}
