package tools

import (
	"encoding/json"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/service"
)

const (
	FluentedElasticsearchImageName = "fluentd_elasticsearch/fluentd"
	FluentedElasticsearchTag       = "v2.8.0"
	ElasticSearchImageName         = "elasticsearch/elasticsearch"
	ElasticSearchTag               = "7.6.2"
)

type EFK struct {
	Cluster *Cluster
	Tool    *model.ClusterTool
}

func (c EFK) setDefaultValue() {
	systemService := service.NewSystemSettingService()
	locahostName := systemService.GetLocalHostName()
	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(c.Tool.Vars), &values)
	values["fluentd-elasticsearch.image.repository"] = fmt.Sprintf("%s:%d/%s", locahostName, constant.LocalDockerRepositoryPort, FluentedElasticsearchImageName)
	values["fluentd-elasticsearch.imageTag"] = FluentedElasticsearchTag
	values["elasticsearch.image"] = fmt.Sprintf("%s:%d/%s", locahostName, constant.LocalDockerRepositoryPort, ElasticSearchImageName)
	values["elasticsearch.imageTag"] = ElasticSearchTag
	values["elasticsearch.replicas"] = 1
	str, _ := json.Marshal(&values)
	c.Tool.Vars = string(str)
}

func NewEFK(cluster *Cluster, tool *model.ClusterTool) (*EFK, error) {
	p := &EFK{
		Tool:    tool,
		Cluster: cluster,
	}
	return p, nil
}

func (c EFK) Install() error {
	c.setDefaultValue()
	if err := installChart(c.Cluster.HelmClient, c.Tool, constant.EFKChartName); err != nil {
		return err
	}
	if err := createRoute(constant.DefaultEFKIngressName, constant.DefaultEFKIngress, constant.DefaultEFKServiceName, 8080, c.Cluster.KubeClient); err != nil {
		return err
	}
	if err := waitForRunning(constant.DefaultEFKDeploymentName, 1, c.Cluster.KubeClient); err != nil {
		return err
	}
	return nil
}

func (c EFK) Uninstall() error {
	return uninstall(c.Tool, constant.DefaultEFKIngressName, c.Cluster.HelmClient, c.Cluster.KubeClient)
}
