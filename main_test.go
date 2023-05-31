package main

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestDo(t *testing.T) {
	clientset := fake.NewSimpleClientset()

	err := do(clientset)
	if err != nil {
		t.Error(err)
	}

	// check the replicas
	created, err := clientset.AppsV1().Deployments("default").Get(context.Background(), "nginx", metav1.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	if *created.Spec.Replicas != 5 {
		t.Error(":(")
	}
}
