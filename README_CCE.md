# 启动参数

## kube-apiserver
```
[Service]
ExecStart=/opt/kube/bin/kube-apiserver \
  --admission-control=NamespaceLifecycle,LimitRanger,ServiceAccount,PersistentVolumeLabel,DefaultStorageClass,ResourceQuota,DefaultTolerationSeconds \
  --advertise-address=100.64.230.11 \
  --external-hostname=100.64.230.11 \
  --allow-privileged=true \
  --apiserver-count=3 \
  --authorization-mode=RBAC,Node \
  --bind-address=100.64.230.11 \
  --client-ca-file=/etc/kubernetes/pki/ca.pem \
  --cloud-config=/etc/kubernetes/cloud.config \
  --cloud-provider=baidubce \
  --enable-swagger-ui=true \
  --etcd-cafile=/etc/etcd/ssl/ca.pem \
  --etcd-certfile=/etc/etcd/ssl/etcd.pem \
  --etcd-keyfile=/etc/etcd/ssl/etcd-key.pem \
  --etcd-servers=https://100.64.230.11:2379,https://100.64.230.10:2379,https://100.64.230.9:2379 \
  --experimental-bootstrap-token-auth=true \
  --feature-gates="Accelerators=true" \
  --insecure-port=0 \
  --logtostderr=true \
  --secure-port=6443 \
  --service-account-key-file=/etc/kubernetes/pki/ca-key.pem \
  --service-cluster-ip-range=172.18.0.0/16 \
  --storage-backend=etcd3 \
  --tls-cert-file=/etc/kubernetes/pki/apiserver.pem \
  --tls-private-key-file=/etc/kubernetes/pki/apiserver-key.pem \
  --v=4
```
1. 增加external-hostname

1. authorization-mode增加Node


## kubelet

```
[Service]
ExecStart=/opt/kube/bin/kubelet \
  --address=192.168.0.10 \
  --allow_privileged=true \
  --fail-swap-on=false \
  --client-ca-file=/etc/kubernetes/pki/ca.pem \
  --cloud-config=/etc/kubernetes/cloud.config \
  --cloud-provider=baidubce \
  --cluster-dns=172.18.0.10 \
  --cluster-domain=cluster.local \
  --docker-root=/var/lib/docker \
  --feature-gates="Accelerators=true" \
  --hostname_override=192.168.0.10 \
  --kubeconfig=/etc/kubernetes/kubelet.conf \
  --logtostderr=true \
  --network-plugin=kubenet \
  --non-masquerade-cidr=172.16.0.0/16 \
  --pod-infra-container-image=hub-readonly.baidubce.com/public/pause:2.0 \
  --pod-manifest-path=/etc/kubernetes/manifests \
  --v=4
```

1. 移除 --require-kubeconfig

1. 移除--apiservers

1. 增加 --fail-swap-on=false