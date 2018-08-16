package add_kubernetes_metadata_mengguang

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/processors"
)

type kubernetesMetadata struct {
	host       string
	podName    string
	httpClient *http.Client
}

func init() {
	processors.RegisterPlugin("add_kubernetes_metadata_mengguang", newAddKubernetesMetadata)
}

func newAddKubernetesMetadata(c *common.Config) (processors.Processor, error) {
	config := struct {
		Host string `config:"host"`
	}{
		Host: "127.0.0.1",
	}

	err := c.Unpack(&config)
	if err != nil {
		return nil, errors.Wrap(err, "fail to unpack the add_kubernetes_metadata_mengguang configuration")
	}

	podName, exists := os.LookupEnv("POD_NAME")
	if !exists {
		logp.Warn("add_kubernetes_metadata_mengguang", "Missing env POD_NAME")
	}

	return &kubernetesMetadata{
		host:       config.Host,
		podName:    podName,
		httpClient: &http.Client{},
	}, nil
}

func (km *kubernetesMetadata) Run(event *beat.Event) (*beat.Event, error) {
	req, _ := http.NewRequest(http.MethodGet, "http://"+km.host+"/api/metadata?podName="+km.podName, nil)
	resp, err := km.httpClient.Do(req)
	if err != nil {
		return event, errors.Wrap(err, "fail to get kubernetes metadata from mengguang")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return event, errors.Wrap(err, "fail to get kubernetes metadata from mengguang when parse body")
	}

	var metadata Metadata
	err = json.Unmarshal(body, &metadata)
	if err != nil {
		return event, errors.Wrap(err, "fail to unmarshal response body")
	}
	m := metadata.ToMapStr()
	meta := common.MapStr{}
	metaIface, ok := event.Fields["kubernetes"]
	if !ok {
		event.Fields["kubernetes"] = common.MapStr{}
	} else {
		meta = metaIface.(common.MapStr)
	}

	meta.Update(m)
	event.Fields["kubernetes"] = meta

	return event, nil
}

func (km *kubernetesMetadata) String() string {
	return "add_kubernetes_metadata_mengguang=[host=" + km.host + ", pod_name=" + km.podName + "]"
}
