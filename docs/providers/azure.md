# Azure

Azure can be used as a cloud provider where your servers reside. The current provider can add a list of IP rules to the security group and optionally remove previously added security rules. If you wish to keep some of the hard coded rules then you can add a certain prefix to their description and Whitelister will not remove them.

## Configuration

Azure provider supports the following configuration

|Key       |Status  |Description|
|----------|--------|-----------|
|SubscriptionID   |required|The subscription ID is a unique uuid string that identifies the Azure subscription|
|ClientID    |required|ID required to connect to Azure|
|ClientSecret    |required|Secret used for establishing connection with Azure|
|TenantID    |required|Unique identifier of the Azure active directory instance|
|RemoveRule|required|Whether to remove un-recognized rules or not. Accepts `true` or `false`|
|KeepRuleDescriptionPrefix|optional|A string value, which when found as a prefix in the description of a security rule then the security rule is not removed|