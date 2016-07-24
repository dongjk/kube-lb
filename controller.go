/*
Copyright 2015 The Kubernetes Authors All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"log"
	"os"
	"os/exec"
	"reflect"
	"text/template"

	"k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/util/flowcontrol"
)

func shellOut(cmd string) {
	out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to execute %v: %v, err: %v", cmd, string(out), err)
	}
}

func main() {
	var endpoints client.EndpointsInterface
	if kubeClient, err := client.NewInCluster(); err != nil {
		log.Fatalf("Failed to create client: %v.", err)
	} else {
		endpoints = kubeClient.Endpoints(api.NamespaceDefault)
	}

	rateLimiter := flowcontrol.NewTokenBucketRateLimiter(0.1, 1)
	epl := &api.EndpointsList{}

	shellOut("nginx")
	for {
		rateLimiter.Accept()
		currentEndpoints, err := endpoints.List(api.ListOptions{})
		if err != nil {
			log.Printf("Error retrieving endpoints: %v", err)
			continue
		}
		if reflect.DeepEqual(currentEndpoints.Items, epl.Items) {
			continue
		}
		epl = currentEndpoints

		tmpl, _ := template.ParseFiles("/tmp/conf.tpl")
		if w, err := os.Create("/etc/nginx/nginx.conf"); err != nil {
			log.Fatalf("Failed to open conf: %v", err)
		} else if err := tmpl.Execute(w, currentEndpoints); err != nil {
			log.Printf("Failed to write template %v", err)
			continue
		}
		shellOut("nginx -s reload")
	}
}
