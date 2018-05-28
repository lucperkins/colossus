# Colossus — An example microservices architecture for Kubernetes

This is an example project that combines several cloud native technologies that I really like and have been meaning to get working in a meaningful way:

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
* The auth service verifies the password passed to the web service. There is only one password that actually works: `tonydanza`. Use any other password and you'll get a `401 Unauthorized` HTTP error.
* The data service handles words or strings that you pass to the web service using the `String` header. The data service simply capitalizes whatever you pass via that header and returns it.

> Wait a second, these services don't do anything meaningful! Nope, they sure don't. But that's okay because the point of this project is to show you how to get the basic (yet not-at-all-trivial) plumbing to work. Colossus is a **boilerplate project** that's meant as a springboard to more complex and meaningful projects.

## What's the point?

Getting all of these technologies to work together was a real challenge. I had to dig through countless GitHub issues and dozens of example projects to make all these things work together. I'm offering this repo as a starter pack for other people with a Bazel monorepo targeting Kubernetes.

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

# Alternatively
$ make k8s-redis-deploy
```

Once all five Redis pods are up and running (you can check using `kubectl get pods -w`), you need to set a password for the authentication service. To set the password as `tonydanza` (which the later curl examples assume):

```bash
$ REDIS_POD=$(kubectl get pods -l app=redis -o jsonpath='{.items[0].metadata.name}')
$ kubectl exec -it $REDIS_POD -- redis-cli -n 0 -h redis-cluster.default.svc.cluster.local SET password tonydanza
```

You can then verify that the password has been set throughout the cluster by running a `GET password` query from a different pod in the cluster:

```bash
$ kubectl exec -it $(kubectl get pods -l app=redis -o jsonpath='{.items[1].metadata.name}') -- redis-cli -n 0 GET password
"tonydanza"
```

Now that Redis is all set up, you can deploy Colossus using one command:

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
$ MINIKUBE_IP=$(minikube ip)
```

Now let's make a request to our web service:

```bash
$ curl -i -XPOST $MINIKUBE_IP/string
HTTP/1.1 401 Unauthorized
Server: nginx/1.13.12
Date: Sun, 27 May 2018 22:30:02 GMT
Content-Type: text/plain; charset=utf-8
Content-Length: 32
Connection: keep-alive
X-Content-Type-Options: nosniff

You cannot access this resource
```

Oops! We need to specify a password using the `Password` header. The password that we supply will be sent to the auth service for verification.

```bash
$ curl -i -XPOST -H Password:foo $MINIKUBE_IP/string
```

Oops! Denied again. Remember: the only password that works is `tonydanza`. Let's try this again:

```bash
$ curl -i -XPOST -H Password:tonydanza $MINIKUBE_IP/string
HTTP/1.1 400 Bad Request
Server: nginx/1.13.12
Date: Sun, 27 May 2018 22:33:31 GMT
Content-Type: text/plain; charset=utf-8
Content-Length: 50
Connection: keep-alive
X-Content-Type-Options: nosniff

You must specify a string using the String header
```

Oops! Forgot to specify a string using the `String` header, which means that our data service isn't even being access. Let's supply a string:

```bash
$ curl -i -XPOST -H Password:tonydanza -H String:"Hello, world" $MINIKUBE_IP/string
HTTP/1.1 200 OK
Server: nginx/1.13.12
Date: Sun, 27 May 2018 22:50:19 GMT
Content-Type: text/plain; charset=utf-8
Content-Length: 12
Connection: keep-alive

HELLO, WORLD%
```

Success! Our `Password` header is authenticating us via the auth service and the data service is handling our data request the way that we would expect. Colossus is a rousing success, folks 👍

## What's next

This is a humble start but I'd like to expand it a great deal in the future. In particular I'd like to add:

* A service mesh like [Conduit](https://conduit.io) or [Istio](https://istio.io/).
* A [Helm](https://helm.sh/) chart
* Some "real" services that do meaningful things, like interact with databases running on k8s or even a cloud service like Google BigTable.
* Real REST capabilities. Right now our web service doesn't do anything cool. At the very least it should provide some interesting CRUD operations.
* More languages. Right now Go and Java are pretty much the only languages that can be easily incorporated into a gRPC-plus-Bazel setup. I'm sure that support for Python, C++, and others is on the way, and I'll use those capabilities as the opportunity arises.
* Componentize service building. Right now each service is its own self-contained universe. I'd like to create a reusable service abstraction (or re-use abstractions built by others) for creating new services, especially more robust configuration management.
* Integration testing for specific services and the whole thing

The good news is that the hard part---especially getting Bazel to build the right things and Kubernetes to use a local image registry---is already behind me, so adding new services is fairly trivial.