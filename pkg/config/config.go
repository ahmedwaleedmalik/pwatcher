package config

import (
	"os"
)

const (
	podFilterEnvVariable       string = "POD_FILTER_KEY"
	namespaceFilterEnvVariable string = "NAMESPACE_FILTER_KEY"
)

var (
	PodFilterKey       string = getPodFilterKey()
	NamespaceFilterKey string = getNamespaceFilterKey()
)

func getPodFilterKey() string {
	podFilterKey, _ := os.LookupEnv(podFilterEnvVariable)
	return podFilterKey
}
func getNamespaceFilterKey() string {
	namespaceFilterKey, _ := os.LookupEnv(namespaceFilterEnvVariable)
	return namespaceFilterKey
}
