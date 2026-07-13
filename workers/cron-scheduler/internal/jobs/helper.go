package data_transfer

import (
	"errors"
	"net"
	"strconv"

	"github.com/HenryNg101/cron-scheduler/internal/model"
)

func mapRow(headers, row []string) map[string]string {
	m := make(map[string]string)
	for i := range headers {
		m[headers[i]] = row[i]
	}
	return m
}

func parseServer(record map[string]string) (*model.Server, error) {
	name := record["name"]
	if name == "" {
		return nil, errors.New("name is required")
	}

	status, err := strconv.ParseBool(record["status"])
	if err != nil {
		return nil, errors.New("invalid status")
	}

	ip := record["ipv4_address"]
	if net.ParseIP(ip) == nil {
		return nil, errors.New("invalid IP")
	}

	port, err := strconv.Atoi(record["port"])
	if err != nil || port < 0 || port > 65535 {
		return nil, errors.New("invalid port")
	}

	protocol := record["protocol"]
	if protocol == "" {
		protocol = "tcp"
	}

	return &model.Server{
		Name:        name,
		Status:      status,
		IPv4Address: ip,
		Port:        uint(port),
		Protocol:    protocol,
	}, nil
}
