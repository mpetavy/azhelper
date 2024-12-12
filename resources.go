package main

import (
	"context"
	"encoding/json"
	"github.com/mpetavy/common"
	"os/exec"
	"time"
)

type Resource struct {
	ChangedTime      time.Time   `json:"changedTime"`
	CreatedTime      time.Time   `json:"createdTime"`
	ExtendedLocation interface{} `json:"extendedLocation"`
	Id               string      `json:"id"`
	Identity         struct {
		PrincipalId            string      `json:"principalId"`
		TenantId               string      `json:"tenantId"`
		Type                   string      `json:"type"`
		UserAssignedIdentities interface{} `json:"userAssignedIdentities"`
	} `json:"identity"`
	Kind              interface{} `json:"kind"`
	Location          string      `json:"location"`
	ManagedBy         interface{} `json:"managedBy"`
	Name              string      `json:"name"`
	Plan              interface{} `json:"plan"`
	Properties        interface{} `json:"properties"`
	ProvisioningState string      `json:"provisioningState"`
	ResourceGroup     string      `json:"resourceGroup"`
	Sku               struct {
		Capacity interface{} `json:"capacity"`
		Family   interface{} `json:"family"`
		Model    interface{} `json:"model"`
		Name     string      `json:"name"`
		Size     interface{} `json:"size"`
		Tier     string      `json:"tier"`
	} `json:"sku"`
	Tags struct {
		ContactEmailAddress string `json:"ContactEmailAddress"`
	} `json:"tags"`
	Type string `json:"type"`
}

func ReadResources() ([]Resource, error) {
	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cmd := exec.Command("az", "resource", "list")

	ba, err := common.ExecuteCmd(cmd)
	if common.Error(err) {
		return nil, err
	}

	list := make([]Resource, 0)
	err = json.Unmarshal(ba, &list)
	if common.Error(err) {
		return nil, err
	}

	return list, nil
}
