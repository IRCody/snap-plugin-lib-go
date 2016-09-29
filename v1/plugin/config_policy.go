/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2016 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package plugin

import (
	"strings"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin/rpc"
)

type ConfigPolicy struct {
	IntegerRules map[string]intRuleKey
	BoolRules    map[string]boolRuleKey
	StringRules  map[string]stringRuleKey
	FloatRules   map[string]floatRuleKey
}

func NewConfigPolicy() *ConfigPolicy {
	return &ConfigPolicy{
		IntegerRules: map[string]intRuleKey{},
		BoolRules:    map[string]boolRuleKey{},
		StringRules:  map[string]stringRuleKey{},
		FloatRules:   map[string]floatRuleKey{},
	}
}

type intRuleKey struct {
	rule integerRule
	key  []string
}

type boolRuleKey struct {
	rule boolRule
	key  []string
}

type floatRuleKey struct {
	rule floatRule
	key  []string
}

type stringRuleKey struct {
	rule stringRule
	key  []string
}

// AddIntRule adds a given integerRule to the IntegerRules map.
// This will overwrite any existing entry.
func (c *ConfigPolicy) AddIntRule(key []string, in integerRule) {
	k := strings.Join(key, ".") // Method used on daemon side in ctree
	c.IntegerRules[k] = intRuleKey{rule: in, key: key}
}

// AddBoolRule adds a given boolRule to the BoolRules map.
// This will overwrite any existing entry.
func (c *ConfigPolicy) AddBoolRule(key []string, in boolRule) {
	k := strings.Join(key, ".") // Method used in daemon/ctree
	c.BoolRules[k] = boolRuleKey{rule: in, key: key}
}

// AddFloatRule adds a given floatRule to the FloatRules map.
// This will overwrite any existing entry.
func (c *ConfigPolicy) AddFloatRule(key []string, in floatRule) {
	k := strings.Join(key, ".") // Method used in daemon/ctree
	c.FloatRules[k] = floatRuleKey{rule: in, key: key}
}

// AddStringRule adds a given stringRule to the StringRules map.
// This will overwrite any existing entry.
func (c *ConfigPolicy) AddStringRule(key []string, in stringRule) {
	k := strings.Join(key, ".") // Method used in daemon/ctree
	c.StringRules[k] = stringRuleKey{rule: in, key: key}
}

func newGetConfigPolicyReply(cfg ConfigPolicy) *rpc.GetConfigPolicyReply {
	ret := &rpc.GetConfigPolicyReply{
		BoolPolicy:    map[string]*rpc.BoolPolicy{},
		FloatPolicy:   map[string]*rpc.FloatPolicy{},
		IntegerPolicy: map[string]*rpc.IntegerPolicy{},
		StringPolicy:  map[string]*rpc.StringPolicy{},
	}

	for k, v := range cfg.IntegerRules {
		r := &rpc.IntegerRule{
			Required:   v.rule.Required,
			Default:    v.rule.Default,
			HasDefault: v.rule.HasDefault,
			Minimum:    v.rule.Minimum,
			HasMin:     v.rule.HasMin,
			Maximum:    v.rule.Maximum,
			HasMax:     v.rule.HasMax,
		}

		if ret.IntegerPolicy[k] == nil {
			ret.IntegerPolicy[k] = &rpc.IntegerPolicy{
				Rules: map[string]*rpc.IntegerRule{},
				Key:   v.key,
			}
		}
		ret.IntegerPolicy[k].Rules[v.rule.Key] = r
	}

	for k, v := range cfg.FloatRules {
		r := &rpc.FloatRule{
			Required:   v.rule.Required,
			Default:    v.rule.Default,
			HasDefault: v.rule.HasDefault,
			Minimum:    v.rule.Minimum,
			HasMin:     v.rule.HasMin,
			Maximum:    v.rule.Maximum,
			HasMax:     v.rule.HasMax,
		}

		if ret.FloatPolicy[k] == nil {
			ret.FloatPolicy[k] = &rpc.FloatPolicy{
				Rules: map[string]*rpc.FloatRule{},
				Key:   v.key,
			}
		}
		ret.FloatPolicy[k].Rules[v.rule.Key] = r
	}

	for k, v := range cfg.StringRules {
		r := &rpc.StringRule{
			Required:   v.rule.Required,
			Default:    v.rule.Default,
			HasDefault: v.rule.HasDefault,
		}

		if ret.StringPolicy[k] == nil {
			ret.StringPolicy[k] = &rpc.StringPolicy{
				Rules: map[string]*rpc.StringRule{},
				Key:   v.key,
			}
		}
		ret.StringPolicy[k].Rules[v.rule.Key] = r
	}

	for k, v := range cfg.BoolRules {
		r := &rpc.BoolRule{
			Required:   v.rule.Required,
			Default:    v.rule.Default,
			HasDefault: v.rule.HasDefault,
		}

		if ret.BoolPolicy[k] == nil {
			ret.BoolPolicy[k] = &rpc.BoolPolicy{
				Rules: map[string]*rpc.BoolRule{},
				Key:   v.key,
			}
		}
		ret.BoolPolicy[k].Rules[v.rule.Key] = r
	}

	return ret
}
