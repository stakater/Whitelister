package utils

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type IpPermission struct {
	IpCidr      *string
	Description *string
	FromPort    *int64
	ToPort      *int64
	IpProtocol  *string
}

//GetLoadBalancerNameFromDNSName gets the name of load balancer from DNS name by splitting the dnsName on '-'
func GetLoadBalancerNameFromDNSName(dnsName string) string {
	return strings.Split(dnsName, "-")[0]
}

//IsEc2IpPermissionEqual Compares two ec2 ips to check if they are equal
func IsEc2IpPermissionEqual(ipPermission1 *ec2.IpPermission, ipPermission2 *ec2.IpPermission) bool {

	if *ipPermission1.FromPort == *ipPermission2.FromPort &&
		*ipPermission1.IpProtocol == *ipPermission2.IpProtocol &&
		*ipPermission1.ToPort == *ipPermission2.ToPort &&
		IsEc2IpRangeEqual(ipPermission1.IpRanges, ipPermission2.IpRanges) &&
		IsEc2Ipv6RangeEqual(ipPermission1.Ipv6Ranges, ipPermission2.Ipv6Ranges) {

		return true
	}
	return false
}

//IsEc2IpRangeEqual comparese to ec2.ipRanges to check if they are equal
func IsEc2IpRangeEqual(ipRanges1 []*ec2.IpRange, ipRanges2 []*ec2.IpRange) bool {
	if len(ipRanges1) != len(ipRanges2) {
		return false
	}
	var ipRangeExists bool

	for _, ipRange1 := range ipRanges1 {
		ipRangeExists = false
		for _, ipRange2 := range ipRanges2 {
			if IsStringEqual(ipRange1.CidrIp, ipRange2.CidrIp) &&
				IsStringEqual(ipRange1.Description, ipRange2.Description) {
				ipRangeExists = true
				break
			}
		}
		if !ipRangeExists {
			return false
		}
	}

	return true
}

//IsEc2Ipv6RangeEqual comparese to ec2.ipv6Ranges to check if they are equal
func IsEc2Ipv6RangeEqual(ipv6Ranges1 []*ec2.Ipv6Range, ipv6Ranges2 []*ec2.Ipv6Range) bool {

	if len(ipv6Ranges1) != len(ipv6Ranges2) {
		return false
	}
	var ipv6RangeExists bool

	for _, ipv6Range1 := range ipv6Ranges1 {
		ipv6RangeExists = false
		for _, ipv6Range2 := range ipv6Ranges2 {
			if IsStringEqual(ipv6Range1.CidrIpv6, ipv6Range2.CidrIpv6) &&
				IsStringEqual(ipv6Range1.Description, ipv6Range2.Description) {
				ipv6RangeExists = true
				break
			}
		}
		if !ipv6RangeExists {
			return false
		}
	}

	return true
}

// IsStringEqual Compares two String pointers with checks for null pointers
func IsStringEqual(val1 *string, val2 *string) bool {
	if val1 == nil && val2 != nil {
		return false
	} else if val1 != nil && val2 == nil {
		return false
	} else {
		return *val1 == *val2
	}
}
