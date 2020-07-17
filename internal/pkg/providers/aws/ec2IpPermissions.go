package aws

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stakater/Whitelister/internal/pkg/utils"
	"regexp"
)

func getEc2IpPermissions(ipPermissions []utils.IpPermission) []*ec2.IpPermission {

	var ec2IpPermissions []*ec2.IpPermission
	for _, ipPermission := range ipPermissions {
		ec2IpPermissions = append(ec2IpPermissions,
			(&ec2.IpPermission{}).
				SetIpProtocol(*ipPermission.IpProtocol).
				SetFromPort(*ipPermission.FromPort).
				SetToPort(*ipPermission.ToPort).
				SetIpRanges(getEc2IpRanges(ipPermission.IpRanges)),
		)
	}

	return ec2IpPermissions
}

func (a *Aws) filterIpPermissions(ipPermissions []*ec2.IpPermission) []*ec2.IpPermission {

	var filteredIpPermissions []*ec2.IpPermission

	for _, ipPermission := range ipPermissions {
		ipPermission.IpRanges = a.filterIpRanges(ipPermission.IpRanges)
		ipPermission.Ipv6Ranges = a.filterIpv6Ranges(ipPermission.Ipv6Ranges)
		//Must be checked otherwise all security rules are removed for a certain port range and protocol
		if len(ipPermission.IpRanges) != 0 || len(ipPermission.Ipv6Ranges) != 0 {
			filteredIpPermissions = append(filteredIpPermissions, ipPermission)
		}
	}

	if len(filteredIpPermissions) == 0 {
		return nil
	}

	return filteredIpPermissions
}

func (a *Aws) filterIpRanges(ipRanges []*ec2.IpRange) []*ec2.IpRange {

	reg, _ := regexp.Compile(a.KeepRuleDescriptionPrefix + ".*$")
	var filteredIpRanges []*ec2.IpRange

	for _, ipRange := range ipRanges {
		if ipRange.Description == nil || !reg.MatchString(*ipRange.Description) {
			filteredIpRanges = append(filteredIpRanges, ipRange)
		}
	}

	if len(filteredIpRanges) == 0 {
		return nil
	}

	return filteredIpRanges
}

func (a *Aws) filterIpv6Ranges(ipv6Ranges []*ec2.Ipv6Range) []*ec2.Ipv6Range {

	reg, _ := regexp.Compile(a.KeepRuleDescriptionPrefix + ".*$")
	var filteredIpv6Ranges []*ec2.Ipv6Range

	for _, ipv6Range := range ipv6Ranges {
		if ipv6Range.Description == nil || !reg.MatchString(*ipv6Range.Description) {
			filteredIpv6Ranges = append(filteredIpv6Ranges, ipv6Range)
		}
	}

	if len(filteredIpv6Ranges) == 0 {
		return nil
	}

	return filteredIpv6Ranges
}

func getEc2IpRanges(ipRanges []*utils.IpRange) []*ec2.IpRange {

	if ipRanges == nil {
		return nil
	}

	var ec2IpRanges []*ec2.IpRange

	for _, ipRange := range ipRanges {
		ec2IpRanges = append(ec2IpRanges, &ec2.IpRange{
			CidrIp:      ipRange.IpCidr,
			Description: ipRange.Description,
		})
	}
	return ec2IpRanges
}
