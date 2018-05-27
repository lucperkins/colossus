# Colossus — An example microservices architecture for Kubernetes

This is an example project that combines several hip technologies that I like:

* The [Bazel](https://bazel.build) build tool
* [Kubernetes](https://kubernetes.io) and [Minikube](https://kubernetes.io/docs/getting-started-guides/minikube/)
* [Protocol Buffers](https://developers.google.com/protocol-buffers/) and [gRPC](https://grpc.io)
* [Docker](https://docker.com)

Colossus is basically a microservice architecture consisting of three services:

Service | Where the code lives | Language
:-------|:---------------------|:--------
An HTTP service that takes web requests. This is the entry point into the backend services. | [`web`](web) | Go
An authentication/authorization service | [`auth`](auth) | Go
A "data" service that handles data requests | [`data`](data) | Java

What do these services actually do?

* The web service requires a `Password` header and a `String` header. The `Password` header is used for authentication and the `String` header is used as a data input. You need to make `POST` requests to the `/string` endpoint.
* The auth service verifies the password passed to the web service. There is only one password that actually works: `tonydanza`. Use any other password and you'll get a `405 Unauthorized` HTTP error.
* The data service handles words or strings that you pass to the web service using the `String` header. The data service simply capitalizes whatever you pass via that header and returns it.

> Wait a second, these services don't do anything meaningful! Nope, they sure don't. But that's okay because the point of this project is to show you how to get the basic (yet not-at-all-trivial) plumbing to work. Colossus is a **boilerplate project** that's meant as a springboard to more complex and meaningful projects.

## What's the point?

Getting all of these technologies to work together was a real challenge. I had to dig through countless GitHub issues and dozens of example projects to make all these things work together.

## Running Colossus locally

In order to run Colossus locally, you'll need to run a local Docker registry. If your Docker daemon is started up, you can run the local registry like this:

```bash
$ docker run -d -p 5000:5000 --restart=always --name registry registry:2

# Alternatively
$ make docker-registry
```

Once the registry is running, you'll need to start up [Minikube](https://kubernetes.io/docs/getting-started-guides/minikube/) in conjunction with an insecure registry (i.e. the Docker registry running locally):

```bash
$ minikube start --insecure-registry localhost:5000

# Alternatively
$ make minikube-start
```

Once Minikube is up and running (use `minikube status` to check), you'll need to enable the ingress add-on:

```bash
$ minikube addons enable ingress

# Alternatively
$ make minikube-setup
```

Now Minikube is all set. The one required dependency for Colossus is a Redis cluster. To run a Redis cluster in Kubernetes-on-Minikube (configuration in [`k8s/redis.yaml`](k8s/redis.yaml)):

```bash
$ kubectl apply -f k8s/redis.yaml
```

Once that's up and running (you can check using `kubectl get pods -w`), you can deploy Colossus using one command:

```bash
$ make deploy
```

Okay, that's suspiciously magical so I'll break it down into pieces. `make deploy` will do the following:

1. Build Docker images for each service using Bazel
1. Upload those images to the local Docker registry (in each case the image will be run using the `--norun` flag, which will upload the images without actually running them)
1. Apply the Kubernetes configuration in [`k8s/colossus.yaml`](k8s/colossus.yaml), which has Kubernetes `Service` and `Deployment` configurations for each of the three services (each of which runs on three instances) as well as an `Ingress` configuration for access to the HTTP service

Run `kubectl get pods` and if all of the pods have the status `Running` then Colossus is ready to take requests!

## Making requests

In order to access the web service, you'll need to get an IP address for Minikube. I recommend setting it as an environment variable:

```bash
$ export MINIKUBE_IP=$(minikube ip)
```

Now let's make a request to our web service:

```bash
$ curl -i $MINIKUBE_IP
```