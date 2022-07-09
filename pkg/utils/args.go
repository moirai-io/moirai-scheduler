package utils

import (
	"encoding/json"
	"fmt"

	configv1beta3 "github.com/rudeigerc/moirai/apis/config/v1beta3"
	"k8s.io/apimachinery/pkg/runtime"
)

// ParsePluginArgs parses the Moirai scheduler plugin arguments.
func ParsePluginArgs(obj runtime.Object, args *configv1beta3.MoiraiArgs) error {
	unknown, ok := obj.(*runtime.Unknown)
	if !ok {
		return fmt.Errorf("unable to parse obj from runtime.Object to runtime.Unknown")
	}

	if unknown.ContentType != "application/json" {
		return fmt.Errorf("unable to parse content type: %v", unknown.ContentType)
	}

	if err := json.Unmarshal(unknown.Raw, args); err != nil {
		return fmt.Errorf("unable to parse args: %v", err)
	}
	return nil
}
