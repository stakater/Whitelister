# ![](assets/web/whitelister-round-100px.png) Whitelister

[![Go Report Card](https://goreportcard.com/badge/github.com/stakater/whitelister?style=flat-square)](https://goreportcard.com/report/github.com/stakater/whitelister)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/stakater/whitelister)
[![Release](https://img.shields.io/github/release/stakater/whitelister.svg?style=flat-square)](https://github.com/stakater/whitelister/releases/latest)
[![GitHub tag](https://img.shields.io/github/tag/stakater/whitelister.svg?style=flat-square)](https://github.com/stakater/whitelister/releases/latest)
[![Docker Pulls](https://img.shields.io/docker/pulls/stakater/whitelister.svg?style=flat-square)](https://hub.docker.com/r/stakater/whitelister/)
[![Docker Stars](https://img.shields.io/docker/stars/stakater/whitelister.svg?style=flat-square)](https://hub.docker.com/r/stakater/whitelister/)
[![MicroBadger Size](https://img.shields.io/microbadger/image-size/stakater/whitelister.svg?style=flat-square)](https://microbadger.com/images/stakater/whitelister)
[![MicroBadger Layers](https://img.shields.io/microbadger/layers/stakater/whitelister.svg?style=flat-square)](https://microbadger.com/images/stakater/whitelister)
[![license](https://img.shields.io/github/license/stakater/whitelister.svg?style=flat-square)](LICENSE)
[![Get started with Stakater](https://stakater.github.io/README/stakater-github-banner.png)](http://stakater.com/?utm_source=whitelister&utm_medium=github)

## Problem

We would like to restrict access to servers to specific addresses only and keep audit trail of IP addresses allowed. 
A tool to manage access to servers based on IP addresses and ports.

## Solution

Whitelister can be used to manage security group rules that control access to servers on different ports.

At the moment it supports Kubernetes as an IP Provider and Aws as the cloud provider.
You can read more about the configuration options [here](docs/config.md)

## Deploying to Kubernetes

You can deploy Whitelister by following methods

### Vanilla Manifests

You can apply vanilla manifests by running the following command

```bash
kubectl apply -f https://raw.githubusercontent.com/stakater/Whitelister/master/deployments/kubernetes/whitelister.yaml
```

Whitelister gets deployed in `default` namespace and searches for ingresses in all namespaces with label name `whitelister` and label value `true`. You will have to modify this file and add your credentials to access the cloud provider.

### Helm Charts

Alternatively if you have configured helm on your cluster, you can add Whitelister to helm from our public chart repository and deploy it via helm using below mentioned commands

```bash
helm repo add stakater https://stakater.github.io/stakater-charts

helm repo update

helm install stakater/whitelister
```

**Note:**  By default whitelister is installed in default namespace. To run in a specific namespace, please run following command. It will install Whitelister in `test` namespace.

```bash
helm install stakater/whitelister --namespace test
```

## Use Case

- Let's say that using [Scaler](https://github.com/stakater/scaler), you stop/destroy your dev severs at night and start/create them again in the morning for cost saving. This can cause your servers to have new IP addresses every day and it can be a tedious job to add IP addresses of multiple Servers to multiple Security Groups everyday.

  This is where Whitelister comes to the rescue. It can automatically detect new server creation by fetching nodes from Kubernetes and then modify security group of selected Ingress to add Ip addresses of new servers and optionally delete IP Addresses for old servers.
- Let's say you have a large team and you don't want to share cloud access with everyone. Now in order to still allow them to be able to access IP restricted servers, you can use a GIT repo which maintains a list of allowed IP addresses and users can add and remove the IP addresses there.

## Help

### Documentation

You can find more documentation [here](docs)

### Have a question?

File a GitHub [issue](https://github.com/stakater/Whitelister/issues), or send us an [email](mailto:stakater@gmail.com).

### Talk to us on Slack

Join and talk to us on Slack for discussing Whitelister

[![Join Slack](https://stakater.github.io/README/stakater-join-slack-btn.png)](https://slack.stakater.com/)
[![Chat](https://stakater.github.io/README/stakater-chat-btn.png)](https://stakater-community.slack.com/messages/C018GBSABM3)

## Contributing

### Bug Reports & Feature Requests

Please use the [issue tracker](https://github.com/stakater/Whitelister/issues) to report any bugs or file feature requests.

### Developing

PRs are welcome. In general, we follow the "fork-and-pull" Git workflow.

 1. **Fork** the repo on GitHub
 2. **Clone** the project to your own machine
 3. **Commit** changes to your own branch
 4. **Push** your work back up to your fork
 5. Submit a **Pull request** so that we can review your changes

NOTE: Be sure to merge the latest from "upstream" before making a pull request!

## Changelog

View our closed [Pull Requests](https://github.com/stakater/Whitelister/pulls?q=is%3Apr+is%3Aclosed).

## License

Apache2 © [Stakater](http://stakater.com)

## About

`Whitelister` is maintained by [Stakater][website]. Like it? Please let us know at <hello@stakater.com>

See [our other projects][community]
or contact us in case of professional services and queries on <hello@stakater.com>

  [website]: http://stakater.com/
  [community]: https://github.com/stakater/
