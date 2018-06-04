package(default_visibility = ["//visibility:public"])

cc_library(
    name = "userinfo-lib",
    srcs = ["userinfo-server.cc"],
    hdrs = ["userinfo-server.h"],
    deps = [
        "@com_github_grpc_grpc//:grpc++",
        "//proto/userinfo:userinfo_cc_grpc",
    ],
)

cc_binary(
    name = "userinfo-bin",
    deps = [
        ":userinfo-lib",
        "//proto/userinfo:userinfo_cc_grpc",
    ],
)