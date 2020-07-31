package config

import "fmt"

type FilterType int

const (
	LoadBalancer FilterType = iota
	SecurityGroup
)

var loadBalancerStr = "LoadBalancer"
var securityGroupStr = "SecurityGroup"

func (filterType FilterType) String() string {
	filterTypes := [...]string{
		loadBalancerStr,
		securityGroupStr,
	}

	if filterType < LoadBalancer || filterType > SecurityGroup {
		return "Unknown"
	}

	return filterTypes[filterType]
}

func toFilterType(filterTypeStr string) (filterType FilterType, err error) {
	switch filterTypeStr {
	case loadBalancerStr:
		filterType = LoadBalancer

	case securityGroupStr:
		filterType = SecurityGroup

	default:
		err = fmt.Errorf("incorrect FilterType :%s provided", filterTypeStr)
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
