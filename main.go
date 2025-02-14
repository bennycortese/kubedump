package main

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func cacheNodeOsImage(osImageString string) {
	fmt.Println(osImageString)
	// TODO - figure out what to put here
}

func updateNodes(clientset kubernetes.Interface) error {
	ctx := context.Background()

	nodesList, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})

	if err != nil {
		return err
	}

	if len(nodesList.Items) == 0 {
		fmt.Println("Oh no!")
	} else {
		for _, node := range nodesList.Items {
			cacheNodeOsImage(node.Status.NodeInfo.OSImage)
		}
	}

	return nil
}

func main() {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", "/home/bennystream/.kube/config")
	if err != nil {
		panic(err)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	if err := do(clientset); err != nil {
		panic(err)
	}
}

func do(clientset kubernetes.Interface) error {

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx",
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"jon": "washere",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"jon": "washere",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:1.14.2",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	updateNodes(clientset)

	ctx := context.Background()

	deploymentsList, err := clientset.AppsV1().Deployments("default").List(ctx, metav1.ListOptions{})

	if err != nil {
		return err
	}

	if len(deploymentsList.Items) == 0 {
		_, err := clientset.AppsV1().Deployments("default").Create(ctx, dep, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}

	if len(deploymentsList.Items) > 0 {
		created := &deploymentsList.Items[0]
		replicas := int32(5)
		created.Spec.Replicas = &replicas

		_, err = clientset.AppsV1().Deployments("default").Update(ctx, created, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		fmt.Println("updated")
	}
	fmt.Println(len(deploymentsList.Items))

	// time.Sleep(1 * time.Second)

	// defer func() {
	// 	err := clientset.AppsV1().Deployments("default").Delete(ctx, dep.Name, metav1.DeleteOptions{})
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println("deleted")
	// }()

	/*
		created, err = clientset.AppsV1().Deployments("default").Get(ctx, dep.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		replicas := int32(2)
		created.Spec.Replicas = &replicas

		_, err = clientset.AppsV1().Deployments("default").Update(ctx, created, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		fmt.Println("updated")
	*/

	// err = wait.PollImmediate(100*time.Millisecond, 5*time.Second, func() (bool, error) {
	// 	created, err := clientset.AppsV1().Deployments("default").Get(ctx, dep.Name, metav1.GetOptions{})
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return false, nil
	// 	}
	// 	return created.Status.ReadyReplicas == replicas, nil
	// })
	// if err != nil {
	// 	return err
	// }
	// fmt.Println("ready")

	return nil
}
