package kubeclient


import (
	"io/ioutil"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"
)

type KubeClient struct {
	Host string
	ca string
	cert string
	key string
	Client *kubernetes.Clientset
}

func NewKubeClientWithFile(host, ca, cert, key string) (*KubeClient, error) {
	client := &KubeClient{
		Host: host,
		ca: ca,
		cert: cert,
		key: key,
	}

	err := client.load()
	return client, err
}

func (this *KubeClient) load() error {

	ca, err := ioutil.ReadFile(this.ca)
	if err != nil {
		return err
	}

	cert, err := ioutil.ReadFile(this.cert)
	if err != nil {
		return err
	}

	key, err := ioutil.ReadFile(this.key)
	if err != nil {
		return err
	}

	cfg := &rest.Config{
		Host: this.Host,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: ca,
			CertData: cert,
			KeyData: key,
		},
	}

	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return err
	}

	this.Client = client
	return nil
}

func NewKubeClient(host, ca, cert, key string) (*KubeClient, error) {
	kube := &KubeClient{}
	cfg := &rest.Config{
		Host: host,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: []byte(ca),
			CertData: []byte(cert),
			KeyData: []byte(key),
		},
	}

	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	kube.Client = client
	return kube, nil
}
