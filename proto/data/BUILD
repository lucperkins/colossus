load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("@io_grpc_grpc_java//:java_grpc_library.bzl", "java_grpc_library")

proto_library(
    name = "data_proto",
    srcs = ["data.proto"],
    visibility = ["//visibility:public"],
)

java_proto_library(
    name = "data_java_proto",
    visibility = ["//visibility:public"],
    deps = [":data_proto"],
)

java_grpc_library(
    name = "data_java_grpc",
    srcs = ["data_proto"],
    visibility = ["//visibility:public"],
    deps = [":data_java_proto"],
)

go_proto_library(
    name = "data_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_grpc"],
    importpath = "github.com/lucperkins/colossus/proto/data",
    proto = ":data_proto",
    visibility = ["//visibility:public"],
)

go_library(
    name = "go_default_library",
    embed = [":data_go_proto"],
    importpath = "github.com/lucperkins/colossus/proto/data",
    visibility = ["//visibility:public"],
)
