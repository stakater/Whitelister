# Whitelister

A tool to manage access to servers based on IP addresses and ports.

## Usage

Whitelister can be used to manage security group rules that control access to servers on different ports.

At the moment it supports Kubernetes as an IP Provider and Aws as the cloud provider.
You can read more about the configuration options [here](docs/config.md)

## Run

## Use Case

Let's say that using [Scaler](https://github.com/stakater/scaler), you stop/destroy your dev severs at night and start/create them again in the morning for cost saving. This can cause your servers to have new IP addresses every day and it can be a tedious job to add IP addresses of multiple Servers to multiple Security Groups everyday.

This is where Whitelister comes to the resuce. It can automatically detect new server creation by fetching nodes from Kubernetes and then modify security group of selected Ingress to add Ip addresses of new servers and optionally delete IP Addresses for old servers.