# Kubernetes

Kubernetes can be used as IP provider to whitelister. Kubernetes IP provider can be used to automatically add the IP addresses of all nodes to the security group of selected ingress.

## Configuration

Kubernetes Ip Provider supports the following configuration options

|Key       |Status  |Description|
|----------|--------|-----------|
|From Port |required|The starting port of the port range to whitelist.|
|To Port   |required|The ending port of the port range to whitelist.|
|IpProtocol|required|The Ip Protocol on which to allow access on the specified port range.|