package ovhclient

import (
	"fmt"
	"strconv"

	"github.com/ovh/go-ovh/ovh"
)

type DNSRecord struct {
	ID        int64  `json:"id"`
	Zone      string `json:"zone"`
	SubDomain string `json:"subDomain"`
	Type      string `json:"type"`
	Target    string `json:"target"`
	TTL       int    `json:"ttl"`
}

type RecordCreation struct {
	SubDomain string `json:"subDomain,omitempty"`
	Type      string `json:"type"`
	Target    string `json:"target"`
	TTL       int    `json:"ttl,omitempty"`
}

type RecordUpdate struct {
	SubDomain string `json:"subDomain,omitempty"`
	Target    string `json:"target,omitempty"`
	TTL       int    `json:"ttl,omitempty"`
}

func ListZones(client *ovh.Client) ([]string, error) {
	var zones []string
	err := client.Get("/domain/zone", &zones)
	return zones, err
}

func ListRecords(client *ovh.Client, zone, recordType, subDomain string) ([]DNSRecord, error) {
	endpoint := fmt.Sprintf("/domain/zone/%s/record", zone)

	type Filter struct {
		FieldType string `json:"fieldType,omitempty"`
		SubDomain string `json:"subDomain,omitempty"`
	}

	var filters []string
	if recordType != "" {
		filters = append(filters, "fieldType="+recordType)
	}
	if subDomain != "" {
		filters = append(filters, "subDomain="+subDomain)
	}

	if len(filters) > 0 {
		endpoint += "?" + filters[0]
	}

	var records []DNSRecord
	err := client.Get(endpoint, &records)
	return records, err
}

func AddRecord(client *ovh.Client, zone, subDomain, recordType, target string, ttl int) (*DNSRecord, error) {
	endpoint := fmt.Sprintf("/domain/zone/%s/record", zone)

	record := RecordCreation{
		SubDomain: subDomain,
		Type:      recordType,
		Target:    target,
	}
	if ttl > 0 {
		record.TTL = ttl
	}

	var result DNSRecord
	err := client.Post(endpoint, record, &result)
	return &result, err
}

func UpdateRecord(client *ovh.Client, zone, recordID, subDomain, recordType, target string, ttl int) (*DNSRecord, error) {
	endpoint := fmt.Sprintf("/domain/zone/%s/record/%s", zone, recordID)

	record := RecordUpdate{}
	if subDomain != "" {
		record.SubDomain = subDomain
	}
	if target != "" {
		record.Target = target
	}
	if ttl > 0 {
		record.TTL = ttl
	}

	var result DNSRecord
	err := client.Put(endpoint, record, &result)
	return &result, err
}

func DeleteRecord(client *ovh.Client, zone, recordID string) error {
	id, err := strconv.ParseInt(recordID, 10, 64)
	if err != nil {
		return fmt.Errorf("ID invalide: %v", err)
	}

	endpoint := fmt.Sprintf("/domain/zone/%s/record/%d", zone, id)
	return client.Delete(endpoint, nil)
}

func RefreshZone(client *ovh.Client, zone string) error {
	endpoint := fmt.Sprintf("/domain/zone/%s/refresh", zone)
	return client.Post(endpoint, nil, nil)
}
