package controller

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"k8s.io/client-go/pkg/api/errors"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

var (
	clusterResourceCreation = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cluster_resource_creation_milliseconds",
		Help: "Time taken to create cluster resource, in milliseconds",
	})
)

func init() {
	prometheus.MustRegister(clusterResourceCreation)
}

// createClusterResource creates the 'cluster' ThirdPartyResource,
// if it does not exist already.
func (c *controller) createClusterResource() error {
	clusterTPR := v1beta1.ThirdPartyResource{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "ThirdPartyResource",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name: ClusterThirdPartyResourceName,
		},
		Description: "A specification of a Kubernetes cluster",
		Versions: []v1beta1.APIVersion{
			v1beta1.APIVersion{
				Name: "v1",
			},
		},
	}

	tprClient := c.clientset.Extensions().ThirdPartyResources()

	start := time.Now()

	log.Println("Creating cluster resource")
	var err error
	if _, err = tprClient.Create(&clusterTPR); err != nil && !errors.IsAlreadyExists(err) {
		return err
	}
	if errors.IsAlreadyExists(err) {
		log.Println("Cluster resource already exists")
	} else {
		log.Println("Cluster resource created")
	}

	clusterResourceCreation.Set(float64(time.Since(start) / time.Millisecond))

	return nil
}
