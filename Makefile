BAZEL = /usr/local/bin/bazel
KCTL  = kubectl

clean:
	$(BAZEL) clean

.PHONY: build
build: clean gazelle
	$(BAZEL) build //...

gazelle-repos:
	$(BAZEL) run //:gazelle -- update-repos -from_file=Gopkg.lock

gazelle: gazelle-repos
	$(BAZEL) run //:gazelle

docker-registry:
	docker run -d -p 5000:5000 --restart=always --name registry registry:2

minikube-start:
	minikube start --insecure-registry localhost:5000

minikube-setup:
	minikube addons enable ingress

docker-local-push: gazelle
	$(BAZEL) run //:colossus-web -- --norun
	$(BAZEL) run //:colossus-auth -- --norun
	$(BAZEL) run //:colossus-data -- --norun

deploy: docker-local-push
	$(KCTL) apply -f k8s/colossus.yaml 

teardown:
	$(KCTL) delete svc,deployment,ing --all

restart: teardown deploy

busybox-run:
	$(KCTL) run curl --image=radial/busyboxplus:curl -i --tty