package kube_api

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	certPath  = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	tokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"
)

func New() (*Api, error) {
	api := Api{
		Address: "https://" + os.Getenv("KUBERNETES_SERVICE_HOST"),
	}

	err := api.setToken()
	if err != nil {
		return nil, err
	}

	err = api.setCertPool()
	if err != nil {
		return nil, err
	}

	return &api, api.createPodClient()
}

func (a *Api) setToken() error {
	readToken, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return fmt.Errorf("cannot read token: %w", err)
	}
	a.token = "Bearer " + string(readToken)
	return nil
}

func (a *Api) setCertPool() error {
	caCert, err := ioutil.ReadFile(certPath)
	if err != nil {
		return fmt.Errorf("cannot get cert, %w", err)
	}
	a.caCertPool = x509.NewCertPool()
	a.caCertPool.AppendCertsFromPEM(caCert)
	return nil
}

func (a *Api) createPodClient() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return fmt.Errorf("failed to create in-cluster config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create clientset: %w", err)
	}
	a.podClient = clientset.CoreV1().Pods(corev1.NamespaceDefault)
	return nil
}

/*func (a *Api) RunWorker() error {
	httpcli := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: a.caCertPool,
			},
		},
	}
	return nil
}*/
