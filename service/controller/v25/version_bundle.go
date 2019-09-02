package v25

import (
	"github.com/giantswarm/versionbundle"
)

func VersionBundle() versionbundle.Bundle {
	return versionbundle.Bundle{
		Changelogs: []versionbundle.Changelog{
			{
				Component:   "ignition",
				Description: "Add name label for default, giantswarm and kube-system namespaces.",
				Kind:        versionbundle.KindAdded,
			},
			{
				Component:   "ignition",
				Description: "Use v1 stable for giantswarm-critical priority class.",
				Kind:        versionbundle.KindFixed,
			},
			{
				Component:   "ignition",
				Description: "Introduce explicit resource reservation for OS resources and container runtime.",
				Kind:        versionbundle.KindAdded,
			},
			{
				Component:   "kvm-operator",
				Description: "Fix NTP server configuration on tenant nodes.",
				Kind:        versionbundle.KindFixed,
			},
			{
				Component:   "kubernetes",
				Description: "Update kubernetes to 1.14.6 (CVE-2019-9512, CVE-2019-9514) https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG-1.14.md#v1146",
				Kind:        versionbundle.KindChanged,
			},
		},
		Components: []versionbundle.Component{
			{
				Name:    "calico",
				Version: "3.7.2",
			},
			{
				Name:    "containerlinux",
				Version: "2135.4.0",
			},
			{
				Name:    "docker",
				Version: "18.06.1",
			},
			{
				Name:    "etcd",
				Version: "3.3.13",
			},
			{
				Name:    "kubernetes",
				Version: "1.14.6",
			},
		},
		Name:    "kvm-operator",
		Version: "3.8.0",
	}
}