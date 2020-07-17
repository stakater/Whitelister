# Configuration

Sample config map looks like this

```yaml
syncInterval: 10s
filter:
  labelName: app
  labelValue: internal-ingress
ipProviders:
  - name: kubernetes
    params:
      FromPort: 0
      ToPort: 65535
      IpProtocol: tcp
  - name: git
    params:
      AccessToken: "ACCESS_TOKEN"
      URL: "http://github.com/example.git"
      Config: "config.yaml"
provider:
  name: "aws"
  params:
    KeepRuleDescriptionPrefix: "Important: "
    Region: "us-west-2"
    RemoveRule: true
    RoleArn: "role-arn"
```

|Key |Status |Description|
|----|-------|-----------|
|syncInterval| required |The interval after which whitelister syncs the Ip Providers input with the security group. Sync interval is a positive sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".|
|filter.labelName| required | Label Name on which to filter ingresses that will provide the load balancers whose security groups need to be modified|
|filter.labelValue| required |Label Value on which to filter ingresses that will provide the load balancers whose security groups need to be modified|
|ipProviders| required, Min length = 1 |List of IP Providers.|
|ipProviders[].name| required |Name of the IP Provider e.g "kubernetes"|
|ipProviders[].params| required |Map to be passed to the IP Provider|
|provider| required |Cloud provider that where the servers are hosted
|provider[].name| required |Name of Cloud Provider e.g "aws"|
|provider[].params| required |Map to be passed to the Cloud Provider|

## Ip Providers

Whitelister supports the following IP Providers

1. [Kubernetes](ipProviders/kubernetes.md)
2. [GitHub](ipProviders/github.md)

## Providers

Whitelister supports the following Providers

1. [Amazon Web Services](providers/aws.md)