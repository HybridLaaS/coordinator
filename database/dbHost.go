package database

import (
	"encoding/json"
	"errors"
)

var ErrHostExists = errors.New("host already exists")

const HOSTS_STATEMENT = `CREATE TABLE IF NOT EXISTS hosts (
	name TEXT PRIMARY KEY NOT NULL,
	health INTEGER NOT NULL DEFAULT 3,
	cpu_count INTEGER NOT NULL,
	cpu_speed_mhz INTEGER NOT NULL,
	cpu_cores INTEGER NOT NULL,
	memory_total_mib INTEGER NOT NULL,
	memory_speed_mhz INTEGER NOT NULL,
	virtual_storage_size_mib INTEGER NOT NULL,
	networking_provider TEXT NOT NULL,
	networking_speed_mbps INTEGER NOT NULL,
	ipmi_address TEXT NOT NULL,
	ipmi_username TEXT NOT NULL,
	ipmi_password TEXT NOT NULL,
	ipmi_redfish_version INTEGER NOT NULL
);`

const INSERT_HOST_STATEMENT = `INSERT INTO hosts (name, health, cpu_count, cpu_speed_mhz, cpu_cores, memory_total_mib, memory_speed_mhz, virtual_storage_size_mib, networking_provider, networking_speed_mbps, ipmi_address, ipmi_username, ipmi_password, ipmi_redfish_version) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
const SELECT_HOST_STATEMENT = `SELECT name, health, cpu_count, cpu_speed_mhz, cpu_cores, memory_total_mib, memory_speed_mhz, virtual_storage_size_mib, networking_provider, networking_speed_mbps, ipmi_address, ipmi_username, ipmi_password, ipmi_redfish_version FROM hosts WHERE name = ?;`
const DELETE_HOST_STATEMENT = `DELETE FROM hosts WHERE name = ?;`
const UPDATE_HOST_HEALTH_STATEMENT = `UPDATE hosts SET health = ? WHERE name = ?;`
const UPDATE_HOST_SPECS_STATEMENT = `UPDATE hosts SET cpu_count = ?, cpu_speed_mhz = ?, cpu_cores = ?, memory_total_mib = ?, memory_speed_mhz = ?, virtual_storage_size_mib = ? WHERE name = ?;`
const UPDATE_HOST_NETWORKING_STATEMENT = `UPDATE hosts SET networking_provider = ?, networking_speed_mbps = ? WHERE name = ?;`
const UPDATE_HOST_IPMI_STATEMENT = `UPDATE hosts SET ipmi_address = ?, ipmi_username = ?, ipmi_password = ?, ipmi_redfish_version = ? WHERE name = ?;`

const (
	HostHealthGood = iota
	HostHealthDegraded
	HostHealthBad
	HostHealthUnknown
)

const (
	HostRedfishVersion_Dell_iDRAC_7 = iota
	HostRedfishVersion_Dell_iDRAC_8
	HostRedfishVersion_Dell_iDRAC_9
)

type DBHost struct {
	Name     string `json:"name"`
	Health   int    `json:"health"`
	Hardware struct {
		CPU struct {
			Count    int `json:"count"`
			SpeedMHz int `json:"speed_mhz"`
			Cores    int `json:"cores"`
		} `json:"cpu"`
		Memory struct {
			SizeMiB  int `json:"total_mib"`
			SpeedMHz int `json:"speed_mhz"`
		} `json:"memory"`
		VirtualStorageSizeMiB int `json:"virtual_storage_size_mib"`
	} `json:"hardware"`
	Networking struct {
		Provider  string `json:"provider"`
		SpeedMbps int    `json:"speed_mbps"`
	} `json:"networking"`
	IPMI struct {
		Address        string `json:"-"`
		Username       string `json:"-"`
		Password       string `json:"-"`
		RedfishVersion int    `json:"-"`
	} `json:"-"`
}

func (h *DBHost) JSON() []byte {
	json, _ := json.Marshal(h)
	return json
}

func HostExists(name string) bool {
	rows, err := QueuedQuery(SELECT_HOST_STATEMENT, name)

	if err != nil {
		return false
	}

	defer rows.Close()
	return rows.Next()
}

func CreateHost(name string, health int, cpuCount, cpuSpeedMHz, cpuCores, memoryTotalMiB, memorySpeedMHz, virtualStorageSizeMiB int, networkingProvider string, networkingSpeedMbps int, ipmiAddress, ipmiUsername, ipmiPassword string, ipmiRedfishVersion int) (*DBHost, error) {
	if HostExists(name) {
		return nil, ErrHostExists
	}

	if err := QueuedExec(INSERT_HOST_STATEMENT, name, health, cpuCount, cpuSpeedMHz, cpuCores, memoryTotalMiB, memorySpeedMHz, virtualStorageSizeMiB, networkingProvider, networkingSpeedMbps, ipmiAddress, ipmiUsername, ipmiPassword, ipmiRedfishVersion); err != nil {
		return nil, err
	}

	return GetHost(name)
}

func GetHost(name string) (*DBHost, error) {
	rows, err := QueuedQuery(SELECT_HOST_STATEMENT, name)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var host DBHost
	err = rows.Scan(&host.Name, &host.Health, &host.Hardware.CPU.Count, &host.Hardware.CPU.SpeedMHz, &host.Hardware.CPU.Cores, &host.Hardware.Memory.SizeMiB, &host.Hardware.Memory.SpeedMHz, &host.Hardware.VirtualStorageSizeMiB, &host.Networking.Provider, &host.Networking.SpeedMbps, &host.IPMI.Address, &host.IPMI.Username, &host.IPMI.Password, &host.IPMI.RedfishVersion)

	if err != nil {
		return nil, err
	}

	return &host, nil
}

func DeleteHost(name string) error {
	return QueuedExec(DELETE_HOST_STATEMENT, name)
}

func UpdateHostHealth(name string, health int) error {
	return QueuedExec(UPDATE_HOST_HEALTH_STATEMENT, health, name)
}

func UpdateHostSpecs(name string, cpuCount, cpuSpeedMHz, cpuCores, memoryTotalMiB, memorySpeedMHz, virtualStorageSizeMiB int) error {
	return QueuedExec(UPDATE_HOST_SPECS_STATEMENT, cpuCount, cpuSpeedMHz, cpuCores, memoryTotalMiB, memorySpeedMHz, virtualStorageSizeMiB, name)
}

func UpdateHostNetworking(name, provider string, speedMbps int) error {
	return QueuedExec(UPDATE_HOST_NETWORKING_STATEMENT, provider, speedMbps, name)
}

func UpdateHostIPMI(name, address, username, password string, redfishVersion int) error {
	return QueuedExec(UPDATE_HOST_IPMI_STATEMENT, address, username, password, redfishVersion, name)
}

func (h *DBHost) T() {}
