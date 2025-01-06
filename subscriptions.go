package main

import (
	"context"
	"encoding/json"
	"github.com/mpetavy/common"
	"os/exec"
	"time"
)

type Subscription struct {
	CloudName           string        `json:"cloudName"`
	HomeTenantId        string        `json:"homeTenantId"`
	Id                  string        `json:"id"`
	IsDefault           bool          `json:"isDefault"`
	ManagedByTenants    []interface{} `json:"managedByTenants"`
	Name                string        `json:"name"`
	State               string        `json:"state"`
	TenantDefaultDomain string        `json:"tenantDefaultDomain"`
	TenantDisplayName   string        `json:"tenantDisplayName"`
	TenantId            string        `json:"tenantId"`
	User                struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"user"`
}

func ReadSubscriptions() ([]Subscription, error) {
	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cmd := exec.Command("az", "account", "list")

	ba, err := common.RunCmd(cmd)
	if common.Error(err) {
		return nil, err
	}

	list := make([]Subscription, 0)
	err = json.Unmarshal(ba, &list)
	if common.Error(err) {
		return nil, err
	}

	return list, nil
}
