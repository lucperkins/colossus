java_library(
    name = "data_java_lib",
    srcs = ["DataHandler.java"],
    visibility = ["//visibility:public"],
    deps = [
        "//proto/data:data_java_grpc",
        "//proto/data:data_java_proto",
        "@io_grpc_grpc_java//core",
        "@io_grpc_grpc_java//netty",
        "@io_grpc_grpc_java//stub",
    ],
)

java_binary(
    name = "data_java_bin",
    main_class = "colossus.DataHandler",
    visibility = ["//visibility:public"],
    runtime_deps = [
        ":data_java_lib",
    ],
)
