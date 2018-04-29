BAZEL = bazel

all: gazelle
	$(BAZEL) build //...

gazelle:
	$(BAZEL) run //:gazelle