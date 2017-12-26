load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "config.go",
        "factory.go",
        "object.go",
        "pool.go",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/jolestar/go-commons-pool/collections:go_default_library",
        "//vendor/github.com/jolestar/go-commons-pool/concurrent:go_default_library",
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
        "//vendor/github.com/jolestar/go-commons-pool/collections:all-srcs",
        "//vendor/github.com/jolestar/go-commons-pool/concurrent:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
