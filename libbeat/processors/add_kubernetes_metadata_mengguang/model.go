package add_kubernetes_metadata_mengguang

import (
	"strings"

	"github.com/elastic/beats/libbeat/common"
)

type Event int

const (
	EventAdd Event = iota
	EventUpdate
	EventDelete
)

type Node struct {
	Name string `json:"name"`
}

type Pod struct {
	Name string `json:"name"`
}

type Container struct {
	Image string `json:"image"`
	Name  string `json:"name"`
}

type Kubernetes struct {
	Pod         Pod               `json:"pod"`
	Namespace   string            `json:"namespace"`
	Node        Node              `json:"node"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	Container   Container         `json:"container,omitempty"`
	Containers  []Container       `json:"containers,omitempty"`
}

type Metadata struct {
	event Event

	ContainerId string `json:"containerId"`
	PodName     string `json:"podName"`

	Kubernetes Kubernetes `json:"kubernetes"`
}

func (m Metadata) ToMapStr() common.MapStr {
	labelMap := common.MapStr{}
	for k, v := range m.Kubernetes.Labels {
		labelMap[k] = v
	}

	annotationsMap := common.MapStr{}
	for k, v := range m.Kubernetes.Annotations {
		annotationsMap[k] = v
	}
	meta := common.MapStr{
		"pod": common.MapStr{
			"name": m.Kubernetes.Pod.Name,
		},
		"node": common.MapStr{
			"name": m.Kubernetes.Node.Name,
		},
		"namespace": m.Kubernetes.Namespace,
	}

	if len(labelMap) != 0 {
		meta["labels"] = labelMap
	}

	if len(annotationsMap) != 0 {
		meta["annotations"] = annotationsMap
	}

	for _, c := range m.Kubernetes.Containers {
		if !strings.HasPrefix(c.Name, "filebeat") {
			meta["container"] = common.MapStr{
				"name":  c.Name,
				"image": c.Image,
			}
			break
		}
	}

	return meta
}
