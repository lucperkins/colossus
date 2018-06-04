load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")

proto_library(
    name = "userinfo_proto",
    srcs = ["userinfo.proto"],
)

go_proto_library(
    name = "userinfo_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_grpc"],
    importpath = "github.com/lucperkins/colossus/proto/userinfo",
    proto = ":userinfo_proto",
    visibility = ["//visibility:public"],
)

cc_proto_library(
    name = "userinfo_cc_proto",
    deps = [":userinfo_proto"],
    visibility = ["//visibility:public"],
)

cc_grpc_library(
    name = "userinfo_cc_grpc",
    srcs = ["userinfo.proto"],
    proto_only = False,
    use_external = True,
    well_known_protos = True,
    deps = [],
    visibility = ["//visibility:public"],
)