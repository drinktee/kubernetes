load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "baidubce.go",
        "baidubce_instances.go",
        "baidubce_loadbalancer.go",
        "baidubce_routes.go",
        "baidubce_util.go",
        "baidubce_zones.go",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/cloudprovider:go_default_library",
        "//pkg/controller:go_default_library",
        "//vendor/github.com/drinktee/bce-sdk-go/bcc:go_default_library",
        "//vendor/github.com/drinktee/bce-sdk-go/bce:go_default_library",
        "//vendor/github.com/drinktee/bce-sdk-go/blb:go_default_library",
        "//vendor/github.com/drinktee/bce-sdk-go/cce:go_default_library",
        "//vendor/github.com/drinktee/bce-sdk-go/clientset:go_default_library",
        "//vendor/github.com/drinktee/bce-sdk-go/eip:go_default_library",
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/k8s.io/api/core/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/types:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["baidubce_test.go"],
    library = ":go_default_library",
    deps = [
        "//vendor/github.com/drinktee/bce-sdk-go/bce:go_default_library",
        "//vendor/github.com/drinktee/bce-sdk-go/clientset:go_default_library",
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
