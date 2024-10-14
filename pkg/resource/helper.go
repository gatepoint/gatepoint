package resource

import (
	"fmt"

	"github.com/gatepoint/gatepoint/api/label"
	"k8s.io/apimachinery/pkg/fields"
)

func ManageByGatepointLabel() string {
	return FmtLabelKeyValue(label.IoKubernetesAppManagedBy, "gatepoint")
}

func FmtLabelKeyValue(instance label.Instance, v string) string {
	return fields.ParseSelectorOrDie(fmt.Sprintf("%s=%s", instance.Name, v)).String()
}

func ManageByGatepointLabelMap() map[string]string {
	return map[string]string{
		label.IoKubernetesAppManagedBy.Name: "gatepoint",
	}
}
