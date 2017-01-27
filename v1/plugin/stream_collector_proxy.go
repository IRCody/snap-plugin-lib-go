package plugin

import (
	"time"

	"golang.org/x/net/context"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin/rpc"
)

type StreamProxy struct {
	pluginProxy

	plugin StreamCollector
}

func (c *StreamProxy) StreamMetrics(stream rpc.StreamCollector_StreamMetricsServer) error {

	inChan := make(chan []Metric)
	outChan, err := c.plugin.StreamMetrics(inChan)
	if err != nil {
		return err
	}
	go streamRecv(c.plugin, inChan, stream)
	for r := range outChan {
		mts := []*rpc.Metric{}
		for _, mt := range r {
			metric, err := toProtoMetric(mt)
			if err != nil {
				return err
			}
			mts = append(mts, metric)
		}
		reply := &rpc.CollectReply{
			Metrics_Reply: &rpc.MetricsReply{Metrics: mts},
		}
		if err := stream.Send(reply); err != nil {
			return err
		}
	}
	return nil
}

func streamRecv(
	plugin StreamCollector,
	ch chan []Metric,
	s rpc.StreamCollector_StreamMetricsServer) {

	for {
		s, err := s.Recv()
		if err != nil {
			// TODO(CDR):handle this error
			return
		}
		if s != nil {
			plugin.SetMaxBuffer(s.MaxMetricsBuffer)
			plugin.SetMaxCollectDuration(time.Duration(s.MaxCollectDuration))
			plugin.SetConfig(s.Other)
			if s.Metrics_Arg != nil {
				metrics := []Metric{}
				for _, mt := range s.Metrics_Arg.Metrics {
					metric := fromProtoMetric(mt)
					metrics = append(metrics, metric)
				}
				ch <- metrics
			}
		}
	}
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
