package builder

import (
	"fmt"

	mariadbv1alpha1 "github.com/mariadb-operator/mariadb-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func MariaDBPort(svc *corev1.Service) (*v1.ServicePort, error) {
	for _, p := range svc.Spec.Ports {
		if p.Name == mariaDbPortName {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("Service port not found")
}

type ServiceOpts struct {
	Labels      map[string]string
	Annotations map[string]string
	Type        corev1.ServiceType
}

func (b *Builder) BuildService(mariadb *mariadbv1alpha1.MariaDB, key types.NamespacedName,
	opts ServiceOpts) (*corev1.Service, error) {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        key.Name,
			Namespace:   key.Namespace,
			Labels:      opts.Labels,
			Annotations: opts.Annotations,
		},
		Spec: corev1.ServiceSpec{
			Ports:    buildPorts(mariadb),
			Selector: opts.Labels,
			Type:     opts.Type,
		},
	}
	if err := controllerutil.SetControllerReference(mariadb, svc, b.scheme); err != nil {
		return nil, fmt.Errorf("error setting controller reference to Service: %v", err)
	}
	return svc, nil
}

func buildPorts(mariadb *mariadbv1alpha1.MariaDB) []v1.ServicePort {
	ports := []v1.ServicePort{
		{
			Name: mariaDbPortName,
			Port: mariadb.Spec.Port,
		},
	}

	if mariadb.Spec.Metrics != nil {
		metricsPort := v1.ServicePort{
			Name: metricsPortName,
			Port: metricsPort,
		}
		ports = append(ports, metricsPort)
	}

	return ports
}
