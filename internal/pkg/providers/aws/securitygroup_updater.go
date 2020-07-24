package aws

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"
	"github.com/stakater/Whitelister/internal/pkg/utils"
)

func (a *Aws) updateSecurityGroup(client *ec2.EC2, securityGroup *ec2.SecurityGroup,
	ipPermissions []*ec2.IpPermission) error {

	if a.RemoveRule {
		a.removeSecurityRules(client, securityGroup, ipPermissions)
	}
	addSecurityRules(client, securityGroup, ipPermissions)

	return nil
}

func addSecurityRules(client *ec2.EC2, securityGroup *ec2.SecurityGroup, ipPermissions []*ec2.IpPermission) {
	var ipPermissionExists bool
	var ipPermissionsToAdd []*ec2.IpPermission

	for _, ipPermission := range ipPermissions {
		ipPermissionExists = false
		for _, securityGroupIpPermission := range securityGroup.IpPermissions {
			if utils.IsEc2IpPermissionEqual(ipPermission, securityGroupIpPermission) {
				ipPermissionExists = true
				break
			}
		}
		if !ipPermissionExists {
			ipPermissionsToAdd = append(ipPermissionsToAdd, ipPermission)
		}
	}
	if len(ipPermissionsToAdd) > 0 {
		logrus.Infof("Adding security rules : %v for security group :%s", ipPermissionsToAdd, *securityGroup.GroupName)
		err := addSecurityGroupIngresses(client, securityGroup, ipPermissionsToAdd)
		if err != nil {
			logrus.Errorf("Error adding security rules for security group %s : %v", *securityGroup.GroupName, err)
		}
	} else {
		logrus.Infof("No security rules to add for security group : %s", *securityGroup.GroupName)
	}
}

func (a *Aws) removeSecurityRules(client *ec2.EC2, securityGroup *ec2.SecurityGroup, ipPermissions []*ec2.IpPermission) {
	var removeIpPermission bool
	var ipPermissionsToRemove []*ec2.IpPermission

	securityGroupFilteredIpPermissions := a.filterIpPermissions(securityGroup.IpPermissions)

	for _, securityGroupIpPermission := range securityGroupFilteredIpPermissions {
		removeIpPermission = true
		for _, ipPermission := range ipPermissions {
			if utils.IsEc2IpPermissionEqual(ipPermission, securityGroupIpPermission) {
				removeIpPermission = false
				break
			}
		}
		if removeIpPermission {
			ipPermissionsToRemove = append(ipPermissionsToRemove, securityGroupIpPermission)
		}
	}

	if len(ipPermissionsToRemove) > 0 {
		logrus.Infof("Removing security rules : %v for security group :%s", ipPermissionsToRemove, *securityGroup.GroupName)
		err := removeSecurityGroupIngresses(client, securityGroup, ipPermissionsToRemove)
		if err != nil {
			logrus.Errorf("Error removing security rules for security group %s : %v", *securityGroup.GroupName, err)
		}
	} else {
		logrus.Infof("No security rules to remove for security group : %s", *securityGroup.GroupName)
	}
}

func addSecurityGroupIngresses(client *ec2.EC2, securityGroup *ec2.SecurityGroup,
	ipPermissions []*ec2.IpPermission) error {

	_, err := client.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
		GroupId:       securityGroup.GroupId,
		IpPermissions: ipPermissions,
	})

	return err
}

func removeSecurityGroupIngresses(client *ec2.EC2, securityGroup *ec2.SecurityGroup,
	ipPermissions []*ec2.IpPermission) error {

	_, err := client.RevokeSecurityGroupIngress(&ec2.RevokeSecurityGroupIngressInput{
		GroupId:       securityGroup.GroupId,
		IpPermissions: ipPermissions,
	})

	return err
}
