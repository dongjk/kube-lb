all: push

TAG = 0.1
PREFIX = gcr.io/google_containers/kube-lb

controller: controller.go
	GOOS=linux go build -a -o controller ./controller.go

container: controller
	docker build -t $(PREFIX):$(TAG) .

clean:
	rm -f controller