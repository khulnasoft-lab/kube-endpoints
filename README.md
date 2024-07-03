## kube-endpoints

Adds health check functionality for maintaining fixed endpoints when accessing external services from within a cluster.
## Notes:
Version v0.1.1 used the sealyun.com domain.
From version v0.2.0 onwards, all domains are khulnasoft-lab.github.io.
Version v0.2.1 adjusted the Hosts configuration. Pay attention during the upgrade.
You can also manually execute the script with the namespace set to xxx:

```
for cep in $(kubectl get cep -n xxx  -o jsonpath={.items[*].metadata.name});do kubectl patch cep -n xxx --type='json' -p='[{"op": "replace", "path": "/metadata/finalizers", "value":[]}]'  $cep;done
```
It is recommended to back up resources before upgrading, delete the existing custom resource (cr), and then recreate it.

## Background

In practice, there is often a need for services in one K8s cluster to access services in another cluster or external services. The usual solution is to manually maintain fixed Services and Endpoints or hard-code IPs in the application configuration. In these cases, there is no health check functionality for external services, making high availability impossible. If high availability is needed, an external high-availability load balancer (LB) is generally introduced, but this adds complexity and many companies are not equipped to introduce it, making it not the optimal solution.
As we know, the main function of Kube-Proxy is to maintain cluster Services and Endpoints and create corresponding IPVS rules on the respective hosts, enabling access via ClusterIP within Pods.
This led to the idea: write a controller that maintains a CRD to automatically create the necessary Service and Endpoint for external services, performing health checks on the external service data (IP
list) within the created Endpoint, and removing the data if the health check fails.

## Introduction

kube-endpoints is a cloud-native, high-reliability, high-performance layer 4 load balancer for accessing external services from within K8s services, with health check functionality.

Features
More cloud-native
Declarative API: The health check definition method is consistent with Kubelet, with familiar syntax and usage
High reliability: Uses native Service and Endpoint resources, avoiding reinventing the wheel
High performance and stability: Uses native IPVS high-performance layer 4 load balancing
Core Advantages
Completely uses native K8s Service and Endpoint resources with no custom IPVS strategies, leveraging K8s Service capabilities for high reliability
Managed via a controller for a CRD resource ClusterEndpoint (abbreviated as cep), eliminating the need to manually manage Service and Endpoint resources
Fully compatible with existing custom Service and Endpoint resources, allowing seamless switching to kube-endpoints management
Uses native IPVS layer 4 load balancing without introducing Nginx, HAProxy, or other load balancers, reducing complexity while meeting high performance and stability requirements
Use Cases
Primarily used when Pods within a cluster need to access external services, such as databases or middleware. The health check capability of kube-endpoints can promptly remove problematic backend services, avoiding impact from single replica failures, and the status can be checked for backend service health and health check failures.

## Helm Installation

```
VERSION="0.2.1"
wget https://github.com/khulnsoft-lab/kube-endpoints/releases/download/v${VERSION}/kube-endpoints-${VERSION}.tgz
helm install -n kube-system kube-endpoints ./kube-endpoints-${VERSION}.tgz
```

## Sealos Installation
```
sealos run khulnsoft-lab/kube-endpoints:v0.2.1
```
Usage

```
apiVersion: khulnasoft-lab.github.io/v1beta1
kind: ClusterEndpoint
metadata:
  name: wordpress
  namespace: default
spec:
  periodSeconds: 10
  ports:
    - name: wp-https
      hosts:
        ## Hosts with the same port
        - 10.33.40.151
        - 10.33.40.152
      protocol: TCP
      port: 38081
      targetPort: 443
      tcpSocket:
        enable: true
      timeoutSeconds: 1
      failureThreshold: 3
      successThreshold: 1
    - name: wp-http
      hosts:
        ## Hosts with the same port
        - 10.33.40.151
        - 10.33.40.152
      protocol: TCP
      port: 38082
      targetPort: 80
      httpGet:
        path: /healthz
        scheme: http
      timeoutSeconds: 1
      failureThreshold: 3
      successThreshold: 1      
    - name: wp-udp
      hosts:
        ## Hosts with the same port
        - 10.33.40.151
        - 10.33.40.152
      protocol: UDP
      port: 38003
      targetPort: 1234
      udpSocket:
        enable: true
        data: "This is flag data for UDP svc test"
      timeoutSeconds: 1
      failureThreshold: 3
      successThreshold: 1
    - name: wp-grpc
      hosts:
        ## Hosts with the same port
        - 10.33.40.151
        - 10.33.40.152
      protocol: TCP
      port: 38083
      targetPort: 8080
      grpc:
        enable: true
      timeoutSeconds: 1
      failureThreshold: 3
      successThreshold: 1
```
## Summary

The introduction of "kube-endpoints" addresses the issue of accessing external services from within a cluster without intruding on the product and while maintaining cloud-native characteristics. This approach will likely become standard practice in future development or operations, offering an elegant solution to certain problems from a development perspective.
