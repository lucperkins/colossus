BAZEL           = bazel
DEP             = dep
KCTL            = kubectl
SKAFFOLD        = skaffold
REDIS_POD       = $(shell $(KCTL) get pods -l app=redis -o jsonpath='{.items[0].metadata.name}')
REDIS_CLI_EXEC  = $(KCTL) exec -it $(REDIS_POD) -- redis-cli

dev:
	$(SKAFFOLD) dev

deploy:
	$(SKAFFOLD) deploy

teardown:
	$(SKAFFOLD) delete

run:
	$(SKAFFOLD) run

clean:
	$(BAZEL) clean --expunge

.PHONY: build
build: clean gazelle
	$(BAZEL) build //...

dep-ensure:
	$(DEP) ensure

gazelle-repos:
	$(BAZEL) run //:gazelle -- update-repos -from_file=Gopkg.lock

gazelle: gazelle-repos
	$(BAZEL) run //:gazelle

go-setup: dep-ensure gazelle

minikube-start:
	minikube start --memory 5120 --cpus 4

minikube-setup:
	minikube addons enable ingress

redis-set-password:
	$(REDIS_CLI_EXEC) SET password tonydanza

redis-get-password:
	$(REDIS_CLI_EXEC) GET password

busybox-run:
	$(KCTL) run curl --image=radial/busyboxplus:curl -i --tty
