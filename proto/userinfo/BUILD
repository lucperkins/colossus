load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")

proto_library(
    name = "userinfo_proto",
    srcs = ["userinfo.proto"],
    visibility = ["//visibility:public"],
)

go_proto_library(
    name = "userinfo_go_grpc",
    compilers = ["@io_bazel_rules_go//proto:go_grpc"],
    importpath = "github.com/lucperkins/colossus/proto/userinfo",
    proto = ":userinfo_proto",
    visibility = ["//visibility:public"],
)

cc_proto_library(
    name = "userinfo_cc_proto",
    visibility = ["//visibility:public"],
    deps = [":userinfo_proto"],
)

cc_grpc_library(
    name = "userinfo_cc_grpc",
    srcs = ["userinfo.proto"],
    proto_only = False,
    use_external = True,
    visibility = ["//visibility:public"],
    well_known_protos = True,
    deps = [],
)

go_library(
    name = "go_default_library",
    embed = [":userinfo_go_grpc"],
    importpath = "github.com/lucperkins/colossus/proto/userinfo",
    visibility = ["//visibility:public"],
)
