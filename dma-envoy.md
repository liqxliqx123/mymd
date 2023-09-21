## dma-envoy

### 启动脚本



```shell
#! /bin/sh

IMAGE=nexus.wisers.com.cn/wisers-dma-base/dma-envoy
IMAGE_TAG=v1.18-latest
REPLICA_NUM=2

if [[ $1 == 'stage' ]]; then
    echo "env is stage"
    NAMESPACE=dma-common
    FILE_PATH=./stage
    HOST=https://api-stage-internal.wisersone.com
elif [[ $1 == 'us' ]]; then
    echo "env is us"
    NAMESPACE=dma-api
    FILE_PATH=./us
    HOST=https://api-us-internal.wisersone.com
    IMAGE=harbor.wisers.com/wisers-dma-hk/dma-envoy
    REPLICA_NUM=1
    LB_ANNOTATION="service.beta.kubernetes.io/aws-load-balancer-internal: \"10.3.0.0/16\""
elif [[ $1 == 'sg' ]]; then
    echo "env is sg"
    NAMESPACE=dma-api
    FILE_PATH=./sg
    HOST=https://api-sg-internal.wisersone.com
    IMAGE=harbor.wisers.com/wisers-dma-hk/dma-envoy
    REPLICA_NUM=1
    LB_ANNOTATION="service.beta.kubernetes.io/aws-load-balancer-internal: \"10.32.0.0/16\""
elif [[ $1 == 'qcloud' ]]; then
    echo "env is qcloud"
    NAMESPACE=dma-api
    FILE_PATH=./qcloud
    HOST=https://api-qcloud-internal.wisersone.com
    REPLICA_NUM=1
    LB_ANNOTATION="service.kubernetes.io/qcloud-loadbalancer-internal-subnetid: subnet-l0p4jztj"
else
    echo "env is prod"
    NAMESPACE=production
    FILE_PATH=./prod
    HOST=https://api-internal.wisersone.com
    LB_ANNOTATION="service.beta.kubernetes.io/aws-load-balancer-internal: \"10.18.0.0/16\""
fi

CONFIG_MAP_NAME=dma-envoy-fs
SERVICE_NAME=dma-envoy
LDS_FILE=$FILE_PATH/envoy-lds.yaml
CDS_FILE=$FILE_PATH/envoy-cds.yaml
BOOTSTRAP_FILE=envoy-bootstrap.yaml

update_envoy_config() {
    echo "--- Create configmap $CONFIG_MAP_NAME for envoy ---"
    kubectl -n $NAMESPACE delete configmap $CONFIG_MAP_NAME
    kubectl -n $NAMESPACE create configmap $CONFIG_MAP_NAME \
        --from-file $LDS_FILE \
        --from-file $CDS_FILE \
        --from-file $BOOTSTRAP_FILE
}

deploy_envoy() {
		#重建config map
    update_envoy_config
    
    #更新并执行deploy.yaml.template
    sed -e "s|\$IMAGE|$IMAGE|" \
        -e "s|\$IMAGE_TAG|$IMAGE_TAG|" \
        -e "s|\$LB_ANNOTATION|$LB_ANNOTATION|" \
        -e "s|\$REPLICA_NUM|$REPLICA_NUM|" deploy.yaml.template | kubectl apply -n $NAMESPACE -f -
    # 重启deployment
    kubectl -n $NAMESPACE rollout restart deployment $SERVICE_NAME
    # 如果5min内还没有完成则返回超时错误
    kubectl rollout status deployment $SERVICE_NAME --timeout=5m --namespace=$NAMESPACE
}

check_deploy() {
    echo "check envoy deployment"
    # TODO: add better check method
}

deploy_envoy
check_deploy
```



### deploy.yaml.template

创建deployment和2个Service

```
dma-envoy-public
dma-envoy
```



### envoy_cds.yaml

配置了envoy的cluster信息，cluster定义了一组上游服务主机，以及如何与这些主机进行通讯、负载均衡、监控检查等

```yaml
resources:
- '@type': type.googleapis.com/envoy.config.cluster.v3.Cluster
  connect_timeout: 1s
  load_assignment:
    cluster_name: dma-ilp-service-stage01
    endpoints:
    - lb_endpoints:
      - endpoint:
          address:
            socket_address:
              address: dma-ilp-service.stage01.svc.cluster.local
              port_value: 80
  max_requests_per_connection: 10
  name: dma-ilp-service-stage01
  type: STRICT_DNS
```

### envoy_lds.yaml

配置了envoy的listener信息, listener负责定义如何接收和处理流量

```yaml
- match:
    prefix: /stage01/ilp/
  route:
    cluster: dma-ilp-service-stage01
    prefix_rewrite: /
    retry_policy:
      retry_on: gateway-error
    timeout: 0s
```



```yaml
            - match:
                grpc: {}
                headers:
                - exact_match: dma-orca-server
                  name: service
                - exact_match: stage01
                  name: env
                prefix: /
              route:
                cluster: dma-orca-server-stage01-grpc
                retry_policy:
                  retry_on: gateway-error
                timeout: 0s
```

### enovy-bootstrap.yaml

```yaml
node:
  cluster: dma-cluster
  id: dma_id1
dynamic_resources:
  cds_config:
    path: /var/lib/envoy/envoy-cds.yaml
  lds_config:
    path: /var/lib/envoy/envoy-lds.yaml
admin:
  address:
    socket_address:
      address: 0.0.0.0
      port_value: 19000
```