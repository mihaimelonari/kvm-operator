package endpoint

import (
	"context"
	"reflect"
	"strconv"
	"testing"

	g8sfake "github.com/giantswarm/apiextensions/v3/pkg/clientset/versioned/fake"
	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

func Test_Resource_Endpoint_containsStrings(t *testing.T) {
	testCases := []struct {
		name            string
		lista           []string
		listb           []string
		containsStrings bool
	}{
		{
			name: "case 0",
			lista: []string{
				"a",
				"b",
				"c",
				"d",
			},
			listb: []string{
				"a",
			},
			containsStrings: true,
		},
		{
			name: "case 1",
			lista: []string{
				"a",
				"b",
				"c",
				"d",
			},
			listb: []string{
				"b",
				"c",
			},
			containsStrings: true,
		},
		{
			name: "case 2",
			lista: []string{
				"a",
				"b",
				"c",
				"d",
			},
			listb: []string{
				"a",
				"b",
				"c",
				"d",
			},
			containsStrings: true,
		},
		{
			name: "case 3",
			lista: []string{
				"a",
				"b",
				"c",
				"d",
			},
			listb: []string{
				"x",
			},
			containsStrings: false,
		},
		{
			name: "case 4",
			lista: []string{
				"a",
				"b",
				"c",
				"d",
			},
			listb: []string{
				"x",
				"y",
			},
			containsStrings: false,
		},
		{
			name: "case 5",
			lista: []string{
				"a",
				"b",
				"c",
				"d",
			},
			listb: []string{
				"a",
				"x",
				"y",
				"z",
			},
			containsStrings: false,
		},
		{
			name: "case 6",
			lista: []string{
				"a",
				"b",
				"c",
				"d",
			},
			listb: []string{
				"a",
				"x",
				"c",
				"d",
			},
			containsStrings: false,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			containsStrings := containsStrings(tc.lista, tc.listb)

			if containsStrings != tc.containsStrings {
				t.Fatalf("\n\n%s\n", cmp.Diff(containsStrings, tc.containsStrings))
			}
		})
	}
}

func Test_Resource_Endpoint_newUpdateChange(t *testing.T) {
	testCases := []struct {
		currentState *Endpoint
		desiredState *Endpoint
		updateChange *corev1.Endpoints
	}{
		{
			currentState: &Endpoint{
				IPs: []string{
					"1.1.1.1",
				},
				Ports: []corev1.EndpointPort{
					{
						Port: 1234,
					},
				},
				ServiceName:      "TestService",
				ServiceNamespace: "TestNamespace",
			},
			desiredState: &Endpoint{
				IPs: []string{
					"1.1.1.1",
				},
				ServiceName:      "TestService",
				ServiceNamespace: "TestNamespace",
			},
			updateChange: &corev1.Endpoints{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "TestService",
					Namespace: "TestNamespace",
				},
			},
		},
		{
			currentState: &Endpoint{
				IPs: []string{
					"1.1.1.1",
					"1.2.3.4",
				},
				Ports: []corev1.EndpointPort{
					{
						Port: 1234,
					},
				},
				ServiceName:      "TestService",
				ServiceNamespace: "TestNamespace",
			},
			desiredState: &Endpoint{
				IPs: []string{
					"1.1.1.1",
				},
				ServiceName:      "TestService",
				ServiceNamespace: "TestNamespace",
			},
			updateChange: &corev1.Endpoints{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "TestService",
					Namespace: "TestNamespace",
				},
			},
		},
		{
			currentState: &Endpoint{
				IPs: []string{
					"5.5.5.5",
					"1.2.3.4",
				},
				Ports: []corev1.EndpointPort{
					{
						Port: 1234,
					},
				},
				ServiceName:      "TestService",
				ServiceNamespace: "TestNamespace",
			},
			desiredState: &Endpoint{
				IPs: []string{
					"1.1.1.1",
				},
				ServiceName:      "TestService",
				ServiceNamespace: "TestNamespace",
			},
			updateChange: &corev1.Endpoints{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "TestService",
					Namespace: "TestNamespace",
				},
				Subsets: []corev1.EndpointSubset{
					{
						Ports: []corev1.EndpointPort{
							{
								Port: 1234,
							},
						},
						Addresses: []corev1.EndpointAddress{
							{
								IP: "5.5.5.5",
							},
							{
								IP: "1.2.3.4",
							},
							{
								IP: "1.1.1.1",
							},
						},
					},
				},
			},
		},
		{
			currentState: &Endpoint{
				Ports: []corev1.EndpointPort{
					{
						Port: 1234,
					},
				},
				ServiceName:      "TestService",
				ServiceNamespace: "TestNamespace",
			},
			desiredState: &Endpoint{
				IPs: []string{
					"1.1.1.1",
				},
				ServiceName:      "TestService",
				ServiceNamespace: "TestNamespace",
			},
			updateChange: &corev1.Endpoints{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "TestService",
					Namespace: "TestNamespace",
				},
				Subsets: []corev1.EndpointSubset{
					{
						Ports: []corev1.EndpointPort{
							{
								Port: 1234,
							},
						},
					},
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var err error

			var r *Resource
			{
				c := Config{
					G8sClient: g8sfake.NewSimpleClientset(),
					K8sClient: k8sfake.NewSimpleClientset(),
					Logger:    microloggertest.New(),
				}

				r, err = New(c)
				if err != nil {
					t.Fatal(err)
				}
			}

			updateChange, err := r.newUpdateChange(context.Background(), nil, tc.currentState, tc.desiredState)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tc.updateChange, updateChange) {
				t.Fatalf("\n\n%s\n", cmp.Diff(updateChange, tc.updateChange))
			}
		})
	}
}
