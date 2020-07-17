package config

import "fmt"

type FilterType int

const (
	LoadBalancer FilterType = iota
	SecurityGroup
)

func (filterType FilterType) String() string {
	filterTypes := [...]string{
		"LoadBalancer",
		"SecurityGroup",
	}

	if filterType < LoadBalancer || filterType > SecurityGroup {
		return "Unknown"
	}

	return filterTypes[filterType]
}

func toFilterType(s string) (u FilterType, err error) {
	switch s {
	case "LoadBalancer":
		u = LoadBalancer

	case "SecurityGroup":
		u = SecurityGroup

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
