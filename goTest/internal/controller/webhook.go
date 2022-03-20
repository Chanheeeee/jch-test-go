package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-docker/internal/model"
	"log"
	"strconv"
	"strings"
	//v1beta1 "k8s.io/api/admission/v1beta1" //AdmissionReview, AdmissionRequest,  PatchTypeJSONPatch, AdmissionResponse
	//corev1 "k8s.io/api/core/v1" //Pod
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1" //status(AdmissionResponse.Result)
)

const (
	MAX_REPLICA   = 3
	MIN_CPU_VALUE = 200
	MAX_CPU_VALUE = 500
)

type Controller struct{}

func (c Controller) Mutate(deployment *appsv1.Deployment) ([]model.JSONPatchEntry, error) { //([]byte, error)
	return checkDeployment(*deployment)
}

func makePatch(opValue string, path string, value []byte) model.JSONPatchEntry {

	return model.JSONPatchEntry{OP: opValue, Path: path, Value: value}
}

func checkDeployment(deployment appsv1.Deployment) ([]model.JSONPatchEntry, error) {

	var patchList []model.JSONPatchEntry

	//replica Cechking
	replicaBytes, err := checkReplicas(&deployment)
	if err != nil {
		log.Printf("@@@@marshall replicas: %s\n", err)
		return nil, err
	}

	log.Printf("@@@@@@@replicaBytes: %s\n", replicaBytes)
	patchList = append(patchList, makePatch("replace", "/spec/replicas", replicaBytes))

	//cpu Resource Checking
	containersBytes, err := checkCPUResource(&deployment)
	if err != nil {
		log.Printf("@@@@check container cpu error: %s\n", err)

		return nil, err
	}
	log.Printf("@@@@@@@containersBytes: %s\n", containersBytes)
	patchList = append(patchList, makePatch("replace", "/spec/template/spec/containers", containersBytes))

	return patchList, nil
}

func checkReplicas(deployment *appsv1.Deployment) ([]byte, error) {
	replicaCount := int32(MAX_REPLICA)
	if *deployment.Spec.Replicas > replicaCount {
		deployment.Spec.Replicas = &replicaCount
	}
	return json.Marshal(&deployment.Spec.Replicas)
}

func checkCPUResource(deployment *appsv1.Deployment) ([]byte, error) {

	for i, container := range deployment.Spec.Template.Spec.Containers {

		cpuToString, _ := json.Marshal(container.Resources.Limits.Cpu())

		if ok := strings.Contains(string(cpuToString[:]), "m"); ok { //isZero() && AsInt64

			var limitCPU = strings.Replace(string(cpuToString[:]), "m", "", 1)
			cpuValue, _ := strconv.ParseInt(limitCPU, 10, 64)
			log.Printf("@@@@cpuValue: %s\n", cpuValue) //이상한값들어가눈디...?

			value, err := checkEachCPU(cpuValue)
			if err != nil {
				log.Printf("@@@@limit check: %s\n", err)
				return nil, err
			}

			res := v1.ResourceList{}
			res[v1.ResourceCPU] = *resource.NewMilliQuantity(value, resource.DecimalSI)
			deployment.Spec.Template.Spec.Containers[i].Resources.Limits = res
		}

	}
	return json.Marshal(&deployment.Spec.Template.Spec.Containers)
}

func checkEachCPU(cpu int64) (int64, error) {

	switch {
	case cpu > MAX_CPU_VALUE:
		log.Printf("500보다 더 큰값이라 실패")
		return 0, errors.New(fmt.Sprintf("usage of CPU > %dm", MAX_CPU_VALUE))

	case cpu >= MIN_CPU_VALUE && cpu <= MAX_CPU_VALUE:
		log.Printf("적당한 값")
		return cpu, nil

	case cpu < MIN_CPU_VALUE:
		log.Printf("최솟값 200보다 작은 값이 들어왔어")
		return 200, nil

	default:
		log.Printf("default")
		return 200, nil
	}
}
