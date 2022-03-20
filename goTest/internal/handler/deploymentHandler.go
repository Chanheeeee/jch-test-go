package handler

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go-docker/internal/controller"
	"go-docker/internal/model"
	"io/ioutil"
	"k8s.io/api/admission/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

type Handler struct {
	controller *controller.Controller
}

//mutating Handler
//log.Println("[IN] mutateHandler")
//defer r.Body.Close()
func (h Handler) MutateHandler(w http.ResponseWriter, r *http.Request) {
	admReview := v1beta1.AdmissionReview{}

	deployment, err := getDeploymentFromBody(r, &admReview)
	if err != nil {
		failedResponse(w, admReview, err)
	}

	patchList, err := h.controller.Mutate(deployment) //response객체만들때 필요한 값 가져옴
	if err != nil {
		failedResponse(w, admReview, err)
		return
	}

	successResponse(w, admReview, patchList)
}

func getDeploymentFromBody(r *http.Request, admReview *v1beta1.AdmissionReview) (*appsv1.Deployment, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &admReview); err != nil {
		return nil, err
	}

	var deployment appsv1.Deployment
	if admReview.Request != nil {
		if err := json.Unmarshal(admReview.Request.Object.Raw, &deployment); err != nil {
			log.Printf("@@@@Couldn't deployment unmarshall raw object: %s\n", err)
			return nil, err
		}
	}

	return &deployment, nil
}

func successResponse(w http.ResponseWriter, admissionReview v1beta1.AdmissionReview, patch []model.JSONPatchEntry) {
	patchBytes, err := json.Marshal(&patch)
	if err != nil {
		log.Errorf("Failed Marshal: %v", err)
	}

	admReview := &v1beta1.AdmissionReview{
		TypeMeta: admissionReview.TypeMeta,
		Response: &v1beta1.AdmissionResponse{
			UID:       admissionReview.Request.UID,
			Allowed:   true,
			Patch:     patchBytes,
			PatchType: func() *v1beta1.PatchType {
				pt := v1beta1.PatchTypeJSONPatch
				return &pt
			}(),
			Result:	   &metav1.Status{Status: "Success",},
		},
	}

	admReviewResult, err := json.Marshal(admReview)
	if err != nil {
		log.Errorf("Failed Marshal: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(admReviewResult)
}

func failedResponse(w http.ResponseWriter, admissionReview v1beta1.AdmissionReview, err error) {
	log.Printf("Error handling webhook request: %v", err)

	admReview := &v1beta1.AdmissionReview{
		TypeMeta: admissionReview.TypeMeta,
		Response: &v1beta1.AdmissionResponse{
			UID:     admissionReview.Request.UID,
			Allowed: false,
			Result:  &metav1.Status{Message: err.Error()},
		},
	}

	admReviewResult, err := json.Marshal(admReview)
	if err != nil {
		log.Errorf("Failed Marshal: %v", err)
	}

	w.WriteHeader(http.StatusInternalServerError)
	_, writeErr := w.Write(admReviewResult)
	if writeErr != nil {
		log.Printf("Could not write response: %v", writeErr)
	}
}
