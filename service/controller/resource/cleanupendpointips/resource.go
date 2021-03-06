package cleanupendpointips

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/tenantcluster/v4/pkg/tenantcluster"
	"k8s.io/client-go/kubernetes"
)

const (
	Name = "cleanupendpointips"
)

type Config struct {
	K8sClient     kubernetes.Interface
	Logger        micrologger.Logger
	TenantCluster tenantcluster.Interface
}

type Resource struct {
	k8sClient     kubernetes.Interface
	logger        micrologger.Logger
	tenantCluster tenantcluster.Interface
}

func New(config Config) (*Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.TenantCluster == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.TenantCluster must not be empty", config)
	}

	r := &Resource{
		k8sClient:     config.K8sClient,
		logger:        config.Logger,
		tenantCluster: config.TenantCluster,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}
