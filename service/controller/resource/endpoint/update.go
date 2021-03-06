package endpoint

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/v4/pkg/resource/crud"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) ApplyUpdateChange(ctx context.Context, obj, updateState interface{}) error {
	endpointToUpdate, err := toK8sEndpoint(updateState)
	if err != nil {
		return microerror.Mask(err)
	}

	if !isEmptyEndpoint(endpointToUpdate) {
		r.logger.Debugf(ctx, "updating endpoint '%s'", endpointToUpdate.GetName())

		_, err = r.k8sClient.CoreV1().Endpoints(endpointToUpdate.Namespace).Update(ctx, endpointToUpdate, v1.UpdateOptions{})
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.Debugf(ctx, "updated endpoint '%s'", endpointToUpdate.GetName())
	} else {
		r.logger.Debugf(ctx, "not updating endpoint '%s'", endpointToUpdate.GetName())
	}

	return nil
}

func (r *Resource) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*crud.Patch, error) {
	createState, err := r.newCreateChange(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	updateState, err := r.newUpdateChange(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	patch := crud.NewPatch()

	patch.SetCreateChange(createState)
	patch.SetUpdateChange(updateState)

	return patch, nil
}

func (r *Resource) newUpdateChange(ctx context.Context, obj, currentState, desiredState interface{}) (*corev1.Endpoints, error) {
	currentEndpoint, err := toEndpoint(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	desiredEndpoint, err := toEndpoint(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var updateChange *corev1.Endpoints
	{
		ips := ipsForUpdateChange(currentEndpoint.IPs, desiredEndpoint.IPs)

		e := &Endpoint{
			ResourceVersion:  currentEndpoint.ResourceVersion,
			ServiceName:      currentEndpoint.ServiceName,
			ServiceNamespace: currentEndpoint.ServiceNamespace,
		}

		if !containsStrings(currentEndpoint.IPs, desiredEndpoint.IPs) {
			e.Addresses = ipsToAddresses(ips)
			e.IPs = ips
			e.Ports = currentEndpoint.Ports
		}

		updateChange = r.newK8sEndpoint(e)
	}

	return updateChange, nil
}

// containsStrings returns true when all items in b are present in a. We use
// this to determine if all of the desired IPs are already present in the
// current IPs, which are the IPs from the k8s endpoint. In case all desired IPs
// are already in the endpoint, we do not need to update the endpoint against
// the Kubernetes API.
func containsStrings(a []string, b []string) bool {
	for _, s := range b {
		if !containsString(a, s) {
			return false
		}
	}

	return true
}

func containsString(list []string, s string) bool {
	for _, item := range list {
		if item == s {
			return true
		}
	}

	return false
}

func ipsForUpdateChange(currentIPs []string, desiredIPs []string) []string {
	var ips []string

	ips = append(ips, currentIPs...)

	for _, ip := range desiredIPs {
		if !containsIP(ips, ip) {
			ips = append(ips, ip)
		}
	}

	if len(currentIPs) > 0 {
		return ips
	}

	return nil
}
