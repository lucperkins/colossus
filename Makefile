BAZEL           = bazel
KCTL            = kubectl
REDIS_POD       = $(shell $(KCTL) get pods -l app=redis -o jsonpath='{.items[0].metadata.name}')
REDIS_CLI_EXEC  = $(KCTL) exec -it $(REDIS_POD) -- redis-cli

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
	$(BAZEL) run //:colossus-web-push
	$(BAZEL) run //:colossus-auth-push
	$(BAZEL) run //:colossus-data-push

k8s-redis-deploy:
	$(KCTL) apply -f k8s/redis.yaml

k8s-colossus-deploy:
	$(KCTL) apply -f k8s/colossus.yaml

redis-set-password:
	$(REDIS_CLI_EXEC) SET password tonydanza

redis-get-password:
	$(REDIS_CLI_EXEC) GET password

deploy: docker-local-push k8s-colossus-deploy

restart-colossus:
	$(KCTL) delete -f k8s/colossus.yaml
	$(KCTL) apply -f k8s/colossus.yaml

teardown:
	$(KCTL) delete svc,deployment,ing --all
	$(KCTL) delete po/busybox

restart: teardown deploy

busybox-run:
	$(KCTL) run curl --image=radial/busyboxplus:curl -i --tty
