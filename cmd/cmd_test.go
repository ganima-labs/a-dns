package cmd

import (
	"testing"
)

func TestNewRootCmd_Help(t *testing.T) {
	cmd := NewRootCmd()

	if cmd.Use != "a-dns" {
		t.Errorf("expected Use a-dns, got %s", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("expected Short to be set")
	}
}

func TestNewRootCmd_ContainsCommands(t *testing.T) {
	cmd := NewRootCmd()

	commands := cmd.Commands()
	if len(commands) == 0 {
		t.Error("expected commands to be registered")
	}

	commandNames := make(map[string]bool)
	for _, c := range commands {
		commandNames[c.Name()] = true
	}

	expectedCommands := []string{
		"list-zones",
		"list-records",
		"add-record",
		"delete-record",
		"update-record",
		"setup",
	}

	for _, expected := range expectedCommands {
		if !commandNames[expected] {
			t.Errorf("expected command %s to be registered", expected)
		}
	}
}

func TestNewListZonesCmd(t *testing.T) {
	cmd := NewListZonesCmd()

	if cmd.Use != "list-zones" {
		t.Errorf("expected Use list-zones, got %s", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("expected Short to be set")
	}
}

func TestNewListRecordsCmd(t *testing.T) {
	cmd := NewListRecordsCmd()

	if cmd.Use != "list-records [zone]" {
		t.Errorf("expected Use list-records [zone], got %s", cmd.Use)
	}

	flags := cmd.Flags()
	if flags.Lookup("zone") == nil {
		t.Error("expected zone flag to be defined")
	}
	if flags.Lookup("type") == nil {
		t.Error("expected type flag to be defined")
	}
	if flags.Lookup("subdomain") == nil {
		t.Error("expected subdomain flag to be defined")
	}
}

func TestNewAddRecordCmd(t *testing.T) {
	cmd := NewAddRecordCmd()

	if cmd.Use != "add-record [zone] [type] [target]" {
		t.Errorf("unexpected Use: %s", cmd.Use)
	}

	flags := cmd.Flags()
	if flags.Lookup("zone") == nil {
		t.Error("expected zone flag to be defined")
	}
	if flags.Lookup("subdomain") == nil {
		t.Error("expected subdomain flag to be defined")
	}
	if flags.Lookup("ttl") == nil {
		t.Error("expected ttl flag to be defined")
	}
}

func TestNewDeleteRecordCmd(t *testing.T) {
	cmd := NewDeleteRecordCmd()

	if cmd.Use != "delete-record [zone] [record-id]" {
		t.Errorf("unexpected Use: %s", cmd.Use)
	}
}

func TestNewUpdateRecordCmd(t *testing.T) {
	cmd := NewUpdateRecordCmd()

	if cmd.Use != "update-record [zone] [record-id] [type] [target]" {
		t.Errorf("unexpected Use: %s", cmd.Use)
	}
}

func TestNewSetupCmd(t *testing.T) {
	cmd := NewSetupCmd()

	if cmd.Use != "setup" {
		t.Errorf("expected Use setup, got %s", cmd.Use)
	}

	flags := cmd.Flags()
	if flags.Lookup("endpoint") == nil {
		t.Error("expected endpoint flag to be defined")
	}
	if flags.Lookup("method") == nil {
		t.Error("expected method flag to be defined")
	}
	if flags.Lookup("app-key") == nil {
		t.Error("expected app-key flag to be defined")
	}
	if flags.Lookup("oauth2-client-id") == nil {
		t.Error("expected oauth2-client-id flag to be defined")
	}
	if flags.Lookup("default-zone") == nil {
		t.Error("expected default-zone flag to be defined")
	}
}
