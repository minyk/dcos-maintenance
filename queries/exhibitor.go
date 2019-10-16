package queries

import (
	"encoding/json"
	"github.com/minyk/dcos-maintenance/client"
	"path"
)

type Exhibitor struct {
	PrefixZnode   func() string
	PrefixCluster func() string
}

func NewExhibitor() *Exhibitor {
	return &Exhibitor{
		PrefixZnode:   func() string { return "/exhibitor/exhibitor/v1/explorer/znode" },
		PrefixCluster: func() string { return "/exhibitor/v1/cluster" },
	}
}

func (q *Exhibitor) DeleteZNode(znode string) error {
	joinedPath := path.Join(q.PrefixZnode(), znode)
	responseContent, err := client.HTTPServiceDelete(joinedPath)
	if err != nil {
		return err
	}

	result := ExhibitorResult{}
	err = json.Unmarshal(responseContent, &result)
	if err != nil {
		return err
	}

	if result.Succeeded {
		client.PrintMessage("Delete succeeded: %s", znode)
	} else {
		client.PrintMessage("Delete failed: %s", result.Message)
	}
	return nil
}

type ExhibitorResult struct {
	Succeeded bool   `json:"succeeded"`
	Message   string `json:"message"`
}
