package ovhclient

import (
	"testing"
)

func TestDNSRecord_Structure(t *testing.T) {
	record := DNSRecord{
		ID:        123456,
		Zone:      "ganima.xyz",
		SubDomain: "www",
		Type:      "A",
		Target:    "1.2.3.4",
		TTL:       3600,
	}

	if record.ID != 123456 {
		t.Errorf("expected ID 123456, got %d", record.ID)
	}
	if record.Type != "A" {
		t.Errorf("expected Type A, got %s", record.Type)
	}
	if record.Target != "1.2.3.4" {
		t.Errorf("expected Target 1.2.3.4, got %s", record.Target)
	}
}

func TestRecordCreation_Structure(t *testing.T) {
	record := RecordCreation{
		SubDomain: "test",
		Type:      "CNAME",
		Target:    "example.com",
		TTL:       0,
	}

	if record.SubDomain != "test" {
		t.Errorf("expected SubDomain test, got %s", record.SubDomain)
	}
	if record.Type != "CNAME" {
		t.Errorf("expected Type CNAME, got %s", record.Type)
	}
	if record.Target != "example.com" {
		t.Errorf("expected Target example.com, got %s", record.Target)
	}
}

func TestRecordUpdate_Structure(t *testing.T) {
	record := RecordUpdate{
		SubDomain: "new",
		Target:    "new.example.com",
		TTL:       7200,
	}

	if record.SubDomain != "new" {
		t.Errorf("expected SubDomain new, got %s", record.SubDomain)
	}
	if record.Target != "new.example.com" {
		t.Errorf("expected Target new.example.com, got %s", record.Target)
	}
	if record.TTL != 7200 {
		t.Errorf("expected TTL 7200, got %d", record.TTL)
	}
}
