load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")

cc_grpc_library(
    name = "userinfo_cc_grpc",
    srcs = ["userinfo.proto"],
    deps = [],
    proto_only = False,
    well_known_protos = "@com_google_protobuf//:well_known_protos",
    use_external = True
)