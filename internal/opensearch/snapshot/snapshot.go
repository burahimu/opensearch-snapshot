package snapshot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dim-ops/opensearch-snapshot/internal/config"
	"github.com/opensearch-project/opensearch-go"
)

var snapshotRepository = "snapshot-" + time.Now().Local().Format("2006-01-02")

type Snapshot interface {
	CreateRepository(client *opensearch.Client, cfg *config.Config) error
	CreateSnapshot(client *opensearch.Client) error
}

func CreateRepository(client *opensearch.Client, cfg *config.Config) error {

	fmt.Println()

	payload := map[string]interface{}{
		"type": "s3",
		"settings": map[string]string{
			"bucket":    cfg.Opensearch.Bucket,
			"base_path": "snapshot-opensearch",
			"region":    cfg.Opensearch.Region,
			"role_arn":  cfg.Opensearch.RoleARN,
		},
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// this method doesn't return an error despite an HTTP 4XX or 5XX, so I catch the status
	rsp, err := client.API.Snapshot.CreateRepository(snapshotRepository, bytes.NewReader(jsonPayload))
	if err != nil || rsp.StatusCode != 200 {
		return fmt.Errorf("HTTP - %d", rsp.StatusCode)
	}

	return nil
}

func CreateSnapshot(client *opensearch.Client) error {

	// this method doesn't return an error despite an HTTP 4XX or 5XX, so I catch the status
	rsp, err := client.API.Snapshot.Create(snapshotRepository, "snapshot")
	if err != nil || rsp.StatusCode != 200 {
		return fmt.Errorf("HTTP - %d", rsp.StatusCode)
	}

	return nil
}
