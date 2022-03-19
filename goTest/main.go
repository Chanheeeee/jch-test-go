package main

import (
	"fmt"
	"net/http"
	"log"
	// "time"
)



// //방법2
// func barHandler (w http.ResponseWriter, r *http.Request) {
// 	//fmt.Fprint(w,r)
// 	name := r.URL.Query().Get("name")
// 	err := json.NewDecoder(r.Body).Decode(user) //r.Body에 json값 들어있음. user형태가 아닌경우엔 err가 나옴
// 	if err != nil { //에러인경우
// 		w.WriteHeader(http.StatusBadRequest)
// 		fmt.Fprint(w, err)
// 		return 
// 	}
// 	user.CreatedAt = time.Now()
// 	// -> 시간 변경된 user값을 다시 json으로 변환

// 	data, _ := json.Marshal(User) //json형태로 인코딩 (byte 배열과 , err) err는 _로 표시

// 	w.Header().Add("content-type", "application/json") //content-type이 json임을 명시
// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprint(w,string(data)) //byte array를 string으로 변환해서 전달

// 	fmt.Fprint(w, "Hello") //responsewriter에 해당 문자를 보내라
// }

// //방법3
// //json을 담을 스트럭쳐
// type User struct{
// 	FirstName 	string
// 	LastName 	string
// 	Email		string
// 	CreatedAt	time.Time
// }

// type FooHandler struct {}

// func (f *FooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	//request에 User Struct 형태의 json이 올텐데 이걸 얻자
// 	user := new(User) //json담을 struct 생성

// 	fmt.Fprint(w, "Hello FOO")
// }



func main() {
	// 어떤 경로에 해당하는 리퀘스트가 들어왔을때 어떤일 할건지(func에 해당하는 일) 핸들러 등록
	// w : response를 write할 수 있는 인자(보낼떄 씀)
	// r : 사용자가 요청한 리퀘스트 정보

	// ################ v1) router ##############
	// mux := http.NewServeMux()
	
	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprint(w, "Hello World") //responsewriter에 해당 문자를 보내라
	// })

	// // mux.HandleFunc("/bar", barHandler)

	// // mux.Handle("/foo", &FooHandler{})

	// http.ListenAndServe(":3000", mux) //listenPort



	// # v2) router
	// mux := http.NewServeMux()

	// mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprint(w, "Hello World") //responsewriter에 해당 문자를 보내라
	// })
	// // mux.HandleFunc("/", handleRoot)
	// // mux.HandleFunc("/mutate", handleMutate)

	// s := &http.Server{
	// 	Addr:           ":8443",
	// 	Handler:        mux,
	// }

	// log.Fatal(s.ListenAndServeTLS("./certs/mutating-test.pem", "./certs/mutating-test.key"))

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World") //responsewriter에 해당 문자를 보내라
	})
	err := http.ListenAndServeTLS(":8443", "./webhook-server-tls.crt", "./webhook-server-tls.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}