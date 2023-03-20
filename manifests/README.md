```shell
kubectl create secret tls mangrove-tls --key mangrove-tls/cluster.peopledata.org.cn.key --cert mangrove-tls/cluster.peopledata.org.cn_bundle.crt -n kube-did
```

```shell
kubectl create secret docker-registry harbor-registry --docker-server=harbor.peopledata.org.cn --docker-username=admin --docker-password=HarborDiD12345 -n kube-did
```
