package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_rust//rust:rust.bzl", "rust_binary", "rust_library")

rust_library(
    name = "userinfo_lib",
    srcs = [
        "main.rs",
    ],
)

rust_binary(
    name = "userinfo_bin",
    srcs = [
        "main.rs",
    ],
    deps = [
        ":userinfo_lib",
    ]
)