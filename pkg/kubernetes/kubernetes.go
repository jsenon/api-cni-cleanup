package kubernetes

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// K8sInternal Connect to Internal k8s Cluster
func K8sInternal() (client *kubernetes.Clientset, err error) {
	config, err := rest.InClusterConfig()
	log.Debug().Msg("Received config object k8s")
	if err != nil {
		log.Error().Msgf("Error config in cluster api kubernetes: ", err.Error())
		return nil, err
	}
	client, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Error().Msgf("Error creation clientset kubernetes: ", err.Error())
		return nil, err
	}
	return client, nil
}

// K8SExternal Connect to External k8s Cluster
func K8SExternal() (client *kubernetes.Clientset, err error) {
	kubeconfig := flag.String("kubeconfig", filepath.Join(homeDir(), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	flag.Parse()
	log.Debug().Msgf("Flag Kubeconfig: ", &kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	log.Debug().Msg("Received config object k8s")
	if err != nil {
		log.Error().Msgf("Error config external cluster api kubernetes: ", err.Error())
		return nil, err
	}
	// create the clientset
	client, err = kubernetes.NewForConfig(config)
	log.Debug().Msg("Received config Clientset")
	if err != nil {
		log.Error().Msgf("Error creation clientset kubernetes: ", err.Error())
		return nil, err
	}
	return client, nil
}

// homeDir set home directory
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}