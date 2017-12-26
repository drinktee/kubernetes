load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "core.go",
        "error.go",
        "request.go",
        "response.go",
        "types.go",
    ],
    visibility = ["//visibility:public"],
    deps = ["//vendor/github.com/drinktee/bce-sdk-go/util:go_default_library"],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
