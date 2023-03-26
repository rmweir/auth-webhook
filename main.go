package main

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"os"

	cattlerbacv1 "github.com/rmweir/role-keeper/api/v1"
	v1 "k8s.io/api/authorization/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(cattlerbacv1.AddToScheme(scheme))
}

func main() {
	restConfig, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		panic(err)
	}

	c, err := client.New(restConfig, client.Options{Scheme: scheme})
	if err != nil {
		panic(err)
	}

	h := handler{
		client: c,
	}

	http.HandleFunc("/", h.ServeHTTP)
	port := os.Getenv("WEBHOOK_PORT")
	/*
		cert := os.Getenv("WEBHOOK_CERT_PATH")
		key := os.Getenv("WEBHOOK_KEY_PATH")
		err := http.ListenAndServeTLS(fmt.Sprintf(":%s", port), cert, key, nil)
	*/

	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		panic(err)
	}
}

type handler struct {
	client client.Client
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sar := v1.SubjectAccessReview{}

	if err := json.NewDecoder(r.Body).Decode(&sar); err != nil {
		return
	}

	var sr cattlerbacv1.SubjectRegistrar
	var err error
	if err = h.client.Get(context.Background(), types.NamespacedName{Name: sar.Spec.User, Namespace: "default"}, &sr); err != nil {
		return
	}

	// sar.Spec.ResourceAttributes
	sar.Status.Allowed = true
	fmt.Printf("ok for: %s\n", sar.String())
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(&sar); err != nil {
		return
	}

	return
}
