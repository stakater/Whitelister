package azure

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/network/mgmt/network"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/sirupsen/logrus"
	"github.com/stakater/Whitelister/internal/pkg/utils"
)

func (a *Azure) updateSecurityRules(securityGroup network.SecurityGroup, ipPermissions []utils.IpPermission) error {
	if a.RemoveRule {
		err := deleteRedundantSecurityRules(a, securityGroup, ipPermissions)
		if err != nil {
			logrus.Error("Error deleting security rules")
			return err
		}
	}

	err := createNewSecurityRules(a, securityGroup, ipPermissions)
	if err != nil {
		logrus.Error("Error creating security rules")
		return err
	}
	return nil
}

func deleteRedundantSecurityRules(a *Azure, securityGroup network.SecurityGroup, ipPermissions []utils.IpPermission) error {

	for _, existingSecurityRule := range *securityGroup.SecurityGroupPropertiesFormat.SecurityRules {

		logrus.Infof("Found security group %s", *existingSecurityRule.Name)
		if !strings.HasPrefix(*existingSecurityRule.Name, a.KeepRuleDescriptionPrefix) {
			err := deleteRuleNotInProvidedIPList(a, ipPermissions, existingSecurityRule, *securityGroup.Name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func deleteRuleNotInProvidedIPList(a *Azure, ipPermissions []utils.IpPermission, existingSecurityRule network.SecurityRule, resourceName string) error {

	for _, ipPermission := range ipPermissions {

		networkProtocol, err := getProtocol(*ipPermission.IpProtocol)
		if err != nil {
			return err
		}
		fromPortStr := strconv.FormatInt(*ipPermission.FromPort, 10)
		toPortStr := strconv.FormatInt(*ipPermission.ToPort, 10)

		if doesNotMatchProvidedRule(ipPermission, existingSecurityRule, fromPortStr, toPortStr, networkProtocol) {

			logrus.Infof("deleting resource %s", *existingSecurityRule.Name)
			err := deleteRule(a, resourceName, *existingSecurityRule.Name)
			if err != nil {
				logrus.Error("Error deleting the security rule " + *existingSecurityRule.Name)
				return err
			}
		}
	}
	return nil
}

func createNewSecurityRules(a *Azure, securityGroup network.SecurityGroup, ipPermissions []utils.IpPermission) error {
	for _, ipPermission := range ipPermissions {
		fromPortStr := strconv.FormatInt(*ipPermission.FromPort, 10)
		toPortStr := strconv.FormatInt(*ipPermission.ToPort, 10)

		for _, ipRange := range *&ipPermission.IpRanges {
			ipDescription := ipRange.Description
			ipCidr := ipRange.IpCidr
			networkProtocol, err := getProtocol(*ipPermission.IpProtocol)
			if err != nil {
				return err
			}

			if len(*securityGroup.SecurityGroupPropertiesFormat.SecurityRules) == 0 {
				logrus.Infof("creating resource %s", *ipDescription)
				err := createSecurityRule(a, *securityGroup.Name, ipDescription, fromPortStr, toPortStr, ipCidr, *ipPermission.IpProtocol)
				if err != nil {
					logrus.Error("Error adding security rule for azure")
					return err
				}
			} else {
				for _, existingSecurityRule := range *securityGroup.SecurityGroupPropertiesFormat.SecurityRules {
					if doesNotMatchProvidedRule(ipPermission, existingSecurityRule, fromPortStr, toPortStr, networkProtocol) {
						logrus.Infof("creating resource %s", *existingSecurityRule.Name)
						err := createSecurityRule(a, *securityGroup.Name, ipDescription, fromPortStr, toPortStr, ipCidr, *ipPermission.IpProtocol)
						if err != nil {
							logrus.Error("Error adding security rule for azure")
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

func createSecurityRule(a *Azure, resourceName string, ipDescription *string, fromPortStr string, toPortStr string, ipRange *string, ipProtocol string) error {

	networkProtocol, err := getProtocol(ipProtocol)
	if err != nil {
		return err
	}

	future, err := a.securityRulesClient.CreateOrUpdate(
		context.TODO(),
		a.ResourceGroupName,
		resourceName,
		*ipDescription,
		network.SecurityRule{
			SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
				Access:                   network.SecurityRuleAccessAllow,
				DestinationAddressPrefix: to.StringPtr("*"),
				DestinationPortRange:     to.StringPtr(fromPortStr + "-" + toPortStr),
				Direction:                network.SecurityRuleDirectionInbound,
				Description:              ipDescription,
				Priority:                 to.Int32Ptr(100),
				Protocol:                 *networkProtocol,
				SourceAddressPrefix:      ipRange,
				SourcePortRange:          to.StringPtr("*"),
			},
		})
	if err != nil {
		return err
	}
	err = future.WaitForCompletionRef(context.TODO(), a.securityRulesClient.Client)
	if err != nil {
		return err
	}
	return nil
}

func getProtocol(ipProtocol string) (*network.SecurityRuleProtocol, error) {
	var networkProtocol network.SecurityRuleProtocol
	switch ipProtocol {
	case "tcp":
		networkProtocol = network.SecurityRuleProtocolTCP
	case "udp":
		networkProtocol = network.SecurityRuleProtocolUDP
	case "*":
		networkProtocol = network.SecurityRuleProtocolAsterisk
	default:
		return nil, errors.New("ip protocol " + ipProtocol + " unidentified")
	}
	return &networkProtocol, nil
}

func deleteRule(a *Azure, resourceName string, ipDescription string) error {
	futureT, err := a.securityRulesClient.Delete(context.TODO(), a.ResourceGroupName, resourceName, ipDescription)
	if err != nil {
		return err
	}
	err = futureT.WaitForCompletionRef(context.TODO(), a.securityRulesClient.Client)
	if err != nil {
		return err
	}
	return nil
}

func doesNotMatchProvidedRule(ipPermission utils.IpPermission, securityRule network.SecurityRule, fromPortStr string, toPortStr string, networkProtocol *network.SecurityRuleProtocol) bool {
	for _, ipRange := range *&ipPermission.IpRanges {

		if *securityRule.Name != *ipRange.Description ||
			securityRule.Access != network.SecurityRuleAccessAllow ||
			*securityRule.DestinationAddressPrefix != "*" ||
			*securityRule.DestinationPortRange != fromPortStr+"-"+toPortStr ||
			securityRule.Direction != network.SecurityRuleDirectionInbound ||
			*securityRule.Description != *ipRange.Description ||
			securityRule.Protocol != *networkProtocol ||
			*securityRule.SourceAddressPrefix != *ipRange.IpCidr ||
			*securityRule.SourcePortRange != "*" {
			return true
		}
	}

	return false
}
