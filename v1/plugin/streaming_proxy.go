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
	"golang.org/x/net/context"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin/rpc"
)

// TODO(danielscottt): plugin panics

type StreamProxy struct {
	pluginProxy

	plugin Streamer
}

func (c *StreamProxy) StreamMetrics(arg *rpc.MetricsArg, stream rpc.StreamCollector_StreamMetricsServer) error {

	metrics := []Metric{}

	for _, mt := range arg.Metrics {
		metric := fromProtoMetric(mt)
		metrics = append(metrics, metric)
	}
	ch, err := c.plugin.StreamMetrics(metrics)
	if err != nil {
		return err
	}
	for r := range ch {
		mts := []*rpc.Metric{}
		for _, mt := range r {
			metric, err := toProtoMetric(mt)
			if err != nil {
				return err
			}
			mts = append(mts, metric)
		}
		reply := &rpc.MetricsReply{Metrics: mts}
		if err := stream.Send(reply); err != nil {
			return err
		}
	}
	return nil
}

func (c *StreamProxy) SetConfig(context.Context, *rpc.ConfigMap) (*rpc.ErrReply, error) {
	return nil, nil
}

func (c *StreamProxy) GetMetricTypes(ctx context.Context, arg *rpc.GetMetricTypesArg) (*rpc.MetricsReply, error) {
	cfg := fromProtoConfig(arg.Config)

	r, err := c.plugin.GetMetricTypes(cfg)
	if err != nil {
		return nil, err
	}
	metrics := []*rpc.Metric{}
	for _, mt := range r {
		// We can ignore this error since we are not returning data from
		// GetMetricTypes.
		metric, _ := toProtoMetric(mt)
		metrics = append(metrics, metric)
	}
	reply := &rpc.MetricsReply{
		Metrics: metrics,
	}
	return reply, nil
}
