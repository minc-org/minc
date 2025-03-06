package kubeconfig

import (
	"fmt"
	"os"

	"github.com/minc-org/minc/pkg/log"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func UpdateKubeConfig(config []byte) error {
	kubeConfigPath := clientcmd.RecommendedHomeFile
	if os.Getenv("KUBECONFIG") != "" {
		kubeConfigPath = os.Getenv("KUBECONFIG")
	}
	log.Debug(fmt.Sprintf("Updating kubeconfig at %s", kubeConfigPath))
	defaultConfig, err := clientcmd.LoadFromFile(kubeConfigPath)
	if err != nil {
		defaultConfig = api.NewConfig()
	}
	uShiftConfig, err := clientcmd.Load(config)
	if err != nil {
		return err
	}
	mergedConfig := mergeConfigs(defaultConfig, uShiftConfig)

	err = clientcmd.WriteToFile(*mergedConfig, kubeConfigPath)
	if err != nil {
		return err
	}
	return nil
}

// mergeConfigs merges two kubeconfig files and returns a single merged configuration.
func mergeConfigs(config1, config2 *api.Config) *api.Config {
	mergedConfig := config1.DeepCopy() // Copy config1 as base

	// Merge clusters
	for name, cluster := range config2.Clusters {
		mergedConfig.Clusters[name] = cluster
	}

	// Merge contexts
	for name, context := range config2.Contexts {
		mergedConfig.Contexts[name] = context

	}

	// Merge AuthInfos (users)
	for name, auth := range config2.AuthInfos {
		mergedConfig.AuthInfos[name] = auth
	}

	// Optionally, update the current context if not set
	if mergedConfig.CurrentContext == "" {
		mergedConfig.CurrentContext = config2.CurrentContext
	}

	return mergedConfig
}
