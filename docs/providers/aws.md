# Amazon Web Services (AWS)

Amazon can be used as a cloud provider where your servers reside. The current provider can add a list of IP rules to the security group and optionally remove previously added security rules. If you wish to keep some of the hard coded rules then you can add a certain prefix to their description and Whitelister will not remove them.

## Configuration

Aws provider supports the following configuration

|Key       |Status  |Description|
|----------|--------|-----------|
|RoleArn   |required|Arn of the role that the whitelister should assume. For enhanced security, this is the only mode allowed right now.|
|Region    |required|Aws Region in which the security group reside|
|RemoveRule|required|Whether to remove un-recognized rules or not. Accepts `true` or `false`|
|KeepRuleDescriptionPrefix|optional|A string value, which when found as a prefix in the description of a security rule then the security rule is not removed|

## Permissions needed for the role

The role whose ARN is specified above should have the following permissions specified in its policy:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "Whitelister-role-policy-rule-1",
            "Effect": "Allow",
            "Action": [
                "ec2:RevokeSecurityGroupIngress",
                "elasticloadbalancing:DescribeLoadBalancers",
                "ec2:AuthorizeSecurityGroupIngress",
                "ec2:DescribeSecurityGroups",
                "ec2:UpdateSecurityGroupRuleDescriptionsIngress"
            ],
            "Resource": "*"
        }
    ]
}
```
