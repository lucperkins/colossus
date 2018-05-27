load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importmap = "vendor/google.golang.org/grpc",
    importpath = "github.com/lucperkins/colossus/auth",
    visibility = ["//visibility:public"],
    # keep
    deps = [
        "//proto/auth:go_default_library",
        "//vendor/github.com/go-redis/redis:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
    ],
)

go_binary(
    name = "auth",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

go_binary(
    name = "auth_linux",
    embed = [":go_default_library"],
    goarch = "amd64",
    goos = "linux",
    pure = "on",
    visibility = ["//visibility:public"],
)
