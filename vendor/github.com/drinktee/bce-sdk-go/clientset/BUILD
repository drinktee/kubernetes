load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["clientset.go"],
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/drinktee/bce-sdk-go/bcc:go_default_library",
        "//vendor/github.com/drinktee/bce-sdk-go/bce:go_default_library",
        "//vendor/github.com/drinktee/bce-sdk-go/blb:go_default_library",
        "//vendor/github.com/drinktee/bce-sdk-go/cce:go_default_library",
        "//vendor/github.com/drinktee/bce-sdk-go/eip:go_default_library",
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
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
