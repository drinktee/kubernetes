load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "client.go",
        "codec.go",
        "compress.go",
        "connection.go",
        "connectionpool.go",
        "data.pb.go",
        "doc.go",
        "header.go",
        "log.go",
        "rpcpackage.go",
        "server.go",
        "utility.go",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/baidu-golang/pbrpc/link:go_default_library",
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/github.com/golang/protobuf/proto:go_default_library",
        "//vendor/github.com/golang/snappy:go_default_library",
        "//vendor/github.com/jolestar/go-commons-pool:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [
        ":package-srcs",
        "//vendor/github.com/baidu-golang/pbrpc/link:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
