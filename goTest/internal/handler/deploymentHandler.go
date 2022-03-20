package handler

import (
	"go-docker/internal/model"
	"go-docker/internal/controller"
	v1beta1 "k8s.io/api/admission/v1beta1"
	"io/ioutil"
	"net/http"
	appsv1 "k8s.io/api/apps/v1"
	"log"
	"fmt"
	"encoding/json"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 이 부분을 라우터로 이동?
func NewHandler() http.Handler {
	log.Println("[IN] newHandler")
	mux := http.NewServeMux()
	mux.HandleFunc("/mutate", mutateHandler)

	return mux
}

//mutating Handler
func mutateHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("[IN] mutateHandler")
	defer r.Body.Close()

	admReview := v1beta1.AdmissionReview{}

	body, err := ioutil.ReadAll(r.Body) 
	if err != nil {
		FailedResponse(admReview, err)
		return
	}

	if err := json.Unmarshal(body, &admReview); err != nil {
		
		log.Printf("@@@@admReview error: %s\n", err)
		FailedResponse(admReview, err)
		return
	}
	log.Printf("admReview: %s\n", admReview)

	ar := admReview.Request
	log.Printf("ar.UID: %s\n", ar.UID)

	var deployment appsv1.Deployment
	if ar != nil {
		if err := json.Unmarshal(ar.Object.Raw, &deployment); err != nil {
			log.Printf("@@@@Couldn't deployment unmarshall raw object: %s\n", err)
			return
		}
	}

	patchList, err := controller.RequestMutate(deployment) //response객체만들때 필요한 값 가져옴
	if err != nil {
		FailedResponse(admReview, err)
		return
	}

	admReviewResult, err := SuccessResponse(admReview, patchList)
	if err != nil {
		FailedResponse(admReview, err)
		return
	}

	responseAdminReview := []byte{}
	responseAdminReview, err = json.Marshal(admReviewResult)
	log.Printf("@@@@@@@responseBody: %s\n", responseAdminReview)

	
	//w.WriteHeader(http.StatusOK)
	//w.Write(responseAdminReview)
	fmt.Fprint(w, responseAdminReview)
}




func SuccessResponse(admissionReview v1beta1.AdmissionReview, patch []model.JSONPatchEntry) (*v1beta1.AdmissionReview, error) {
	patchBytes, err := json.Marshal(&patch)
	if err != nil {
		return nil, err
	}

	return &v1beta1.AdmissionReview{
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
	}, nil
}

func FailedResponse(admissionReview v1beta1.AdmissionReview, err error) *v1beta1.AdmissionReview {
	return &v1beta1.AdmissionReview{
		TypeMeta: admissionReview.TypeMeta,
		Response: &v1beta1.AdmissionResponse{
			UID:     admissionReview.Request.UID,
			Allowed: false,
			Result:  &metav1.Status{Message: err.Error()},
		},
	}
}
