// Package ipmimon provides interface to IPMI monitoing
package ipmimon

import (
	"context"
	"fmt"
	"os/exec"
	"syscall"

	"github.com/gocarina/gocsv"
)

// Item represents one line from the report
//
//	ipmimonitoring --comma-separated-output | head -n1
//	ID,Name,Type,State,Reading,Units,Event
type Item struct {
	ID      int    `csv:"ID" json:"id"`
	Name    string `csv:"Name" json:"name"`
	Type    string `csv:"Type" json:"type"`
	State   string `csv:"State" json:"state"`
	Reading string `csv:"Reading" json:"reading"`
	Units   string `csv:"Units" json:"units"`
	Event   string `csv:"Event" json:"event"`
}

// Report represent parsed report
type Report = []Item

// ParseCSV parses report csv into Report
func ParseCSV(data []byte) (Report, error) {
	items := make([]Item, 0)

	err := gocsv.UnmarshalBytes(data, &items)
	if err != nil {
		return nil, fmt.Errorf("csv parse error: %w", err)
	}

	return Report(items), nil
}

// GetReport run ipmimonitoring and parse report
func GetReport(ctx context.Context) (Report, error) {
	cmd := exec.CommandContext(ctx, "ipmimonitoring", "--comma-separated-output")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGKILL,
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ipmimonitoring command failed: %w", err)
	}

	return ParseCSV(out)
}
