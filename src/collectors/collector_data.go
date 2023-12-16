// Copyright 2021 Akamai Technologies, Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collectors

import (
	"errors"

	"github.com/prometheus/common/log"
)

const (
	HoursInDay            = 24
	trafficReportInterval = 5 // mins
)

var (
	defaultItemConfig          = ItemConfig{}
	defaultMarketHistoryConfig = MarketHistoryConfig{}
	DefaultEveSIExporter       = EveSIExporter{
		Items:         make([]*ItemConfig, 0),
		MarketHistory: make([]*MarketHistoryConfig, 0),
	}
)

type EveSIExporter struct {
	Items         []*ItemConfig          `yaml:"properties,omitempty"`
	MarketHistory []*MarketHistoryConfig `yaml:"datacenters,omitempty"`
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (c *EveSIExporter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultDomainTraffic
	type plain DomainTraffic
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	log.Debugf("Domain: [%v]", *c)
	if c.Name == "" {
		return errors.New("required domain name is empty")
	}
	if (c.Properties == nil || len(c.Properties) < 1) && (c.Datacenters == nil || len(c.Datacenters) < 1) && (c.Liveness == nil || len(c.Liveness) < 1) {
		return errors.New("No property, datacenter or liveness configs to collect")
	}

	return nil
}

type ItemConfig struct {
	Name   string `yaml:"item_name"`
	ItemID []int  `yaml:"item,omitempty"`
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (p *ItemConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*p = defaultItemConfig
	type plain ItemConfig
	if err := unmarshal((*plain)(p)); err != nil {
		return err
	}
	log.Debugf("Item: [%v]", *p)
	if p.Name == "" {
		return errors.New("required Item name is empty")
	}

	return nil
}

type MarketHistoryConfig struct {
	RegionName []string `yaml:"region_name"`
	RegionID   []int    `yaml:"region_id,omitempty"`
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (p *MarketHistoryConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*p = defaultMarketHistoryConfig
	type plain MarketHistoryConfig
	if err := unmarshal((*plain)(p)); err != nil {
		return err
	}
	log.Debugf("Region: [%v]", *p)
	if p.RegionName[] == "" {
		return errors.New("required region name is empty")
	}

	return nil
}
