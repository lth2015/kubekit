package appmachinery

import (
	"strings"
	"k8s.io/client-go/kubernetes"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"fmt"
)

type Owner struct {
	Kind string `json:"kind"`
	Name string `json:""`
	Controller *bool `json:"controller,omitempty"`
	ApiVersion string `json:"apiVersion"`
}

func GetOwner(clientset *kubernetes.Clientset, namespace, name string) (*Owner, error) {
	pod, err := clientset.CoreV1().Pods(namespace).Get(name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var owner *Owner
	for _, ref := range pod.OwnerReferences {
		switch ref.Kind {
		case "ReplicaSet":
			owner, err = getDeploymentByReplicaSet(clientset, pod.Namespace, ref.Name)
		case "StatefulSet":
			owner, err = getOwner(ref)
		case "DaemonSet":
			owner, err = getOwner(ref)
		case "Job":
			owner, err = getCronJobByJob(clientset, pod.Namespace, ref.Name)
		default:
			owner, err = getUnknownOwner()
		}
	}

	return owner, err
}

// Pod controlled by ReplicaSet
func getDeploymentByReplicaSet(clientset *kubernetes.Clientset, namespace, name string) (*Owner, error) {
	rs, err := clientset.AppsV1().ReplicaSets(namespace).Get(name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	owner := &Owner{}
	for _, ref := range rs.OwnerReferences {
		if strings.EqualFold(ref.Kind, "Deployment") {
			owner.Kind = ref.Kind
			owner.Name = ref.Name
			owner.ApiVersion = ref.APIVersion
			owner.Controller = ref.Controller
			return owner, nil
		}
	}

	return owner, fmt.Errorf("ReplicaSet %s/%s has no controller of deployment", namespace, name)
}

// Pod controlled by StatefulSet/DaemonSet
func getOwner(reference meta_v1.OwnerReference) (*Owner, error) {
	owner := &Owner{}
	owner.Kind = reference.Kind
	owner.Name = reference.Name
	owner.Controller = reference.Controller
	owner.ApiVersion = reference.APIVersion
	return owner, nil
}

// Pod controlled by Job
func getCronJobByJob(clientset *kubernetes.Clientset, namespace, name string) (*Owner, error) {
	job, err := clientset.BatchV1().Jobs(namespace).Get(name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	owner := &Owner{}
	for _, ref := range job.OwnerReferences {
		if strings.EqualFold(ref.Kind, "CronJob") {
			owner.Kind = ref.Kind
			owner.Name = ref.Name
			owner.Controller = ref.Controller
			owner.ApiVersion = ref.APIVersion
			return owner, nil
		}
	}
	return owner, fmt.Errorf("Job %s/%s has no controller of cronjob", namespace, name)
}

// Unknown controller
func getUnknownOwner(reference meta_v1.OwnerReference) (*Owner, error) {
	owner := &Owner{}
	b := false
	if owner.Kind = reference.Kind; strings.EqualFold(owner.Kind, "") {
		owner.Kind = "Unknown"
	}

	if owner.Name = reference.Name; strings.EqualFold(owner.Name, "") {
		owner.Name = "Unknown"
	}

	if owner.ApiVersion = reference.APIVersion; strings.EqualFold(owner.ApiVersion, "") {
		owner.ApiVersion = "Unknown"
	}
	owner.Controller = &b
	return owner, fmt.Errorf("Unknown resource of Unknown controller")
}
