package config

import "fmt"

type FilterType int

const (
	loadBalancer FilterType = iota
	securityGroup
)

var LoadBalancerStr = "loadBalancer"
var SecurityGroupStr = "securityGroup"

func (filterType FilterType) String() string {
	filterTypes := [...]string{
		"loadBalancer",
		"securityGroup",
	}

	if filterType < loadBalancer || filterType > securityGroup {
		return "Unknown"
	}

	return filterTypes[filterType]
}

func toFilterType(s string) (u FilterType, err error) {
	switch s {
	case "loadBalancer":
		u = loadBalancer

	case "securityGroup":
		u = securityGroup

	default:
		err = fmt.Errorf("incorrect FilterType :%s provided", s)
	}
	return
}

func (filterType *FilterType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string

	if err := unmarshal(&value); err != nil {
		return err
	}

	FilterType, err := toFilterType(value)
	if err != nil {
		return err
	}
	*filterType = FilterType

	return nil
}
