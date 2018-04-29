load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["root.go"],
    importpath = "github.com/lucperkins/colossus/cli/cmd",
    visibility = ["//visibility:public"],
    deps = ["//vendor/github.com/spf13/cobra:go_default_library"],
)
