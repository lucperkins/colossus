load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")

proto_library(
    name = "userinfo_proto",
    srcs = ["userinfo.proto"],
)

cc_proto_library(
    name = "userinfo_cc_proto",
    deps = [":userinfo_proto"],
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