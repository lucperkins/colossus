load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "github.com/lucperkins/colossus/auth",
    visibility = ["//visibility:public"],
    deps = [
        "//proto/auth:go_default_library",
        "@com_github_go_redis_redis//:go_default_library",
        "@com_github_grpc_ecosystem_go_grpc_prometheus//:go_default_library",
        "@com_github_prometheus_client_golang//prometheus:go_default_library",
        "@com_github_prometheus_client_golang//prometheus/promhttp:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
    ],
)

go_binary(
    name = "auth_bin",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

go_binary(
    name = "auth_linux_bin",
    embed = [":go_default_library"],
    goarch = "amd64",
    goos = "linux",
    pure = "on",
    visibility = ["//visibility:public"],
)
