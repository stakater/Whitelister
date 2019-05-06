# GitHub

Github can be used as IP provider to whitelister. Github IP provider can be used to automatically add the IP addresses of all nodes to the security group of selected ingress.

## Configuration

Github Ip Provider supports the following configuration options

|Key       |Status  |Description|
|----------|--------|-----------|
|AccessToken |required|Access token generated from Github account.| 
|URL   |required|URL of the repository.|
|Config|optional|path of the config file within the repository (by default "config.yaml").|