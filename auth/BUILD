load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importmap = "vendor/google.golang.org/grpc",
    importpath = "github.com/lucperkins/colossus/auth",
    visibility = ["//visibility:public"],
    deps = [
        "//proto/auth:go_default_library",
        "@com_github_fluent_fluent_logger_golang//fluent:go_default_library",
        "@com_github_go_redis_redis//:go_default_library",
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
