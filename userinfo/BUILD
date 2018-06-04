package(default_visibility = ["//visibility:public"])

cc_library(
    name = "userinfo-lib",
    srcs = ["userinfo-server.cc"],
    hdrs = ["userinfo-server.h"],
    deps = [
        "//proto/userinfo:userinfo_cc_grpc",
        "@com_github_grpc_grpc//:grpc++",
    ],
)

cc_binary(
    name = "userinfo-bin",
    deps = [
        ":userinfo-lib",
    ],
)
