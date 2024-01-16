package streams

import (
	"fmt"

	"go.opentelemetry.io/collector/pdata/pmetric"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/deltatocumulativeprocessor/internal/metrics"
)

func NewTracker(aggr Aggregator) metrics.Aggregator {
	tracker := &Tracker{
		series: make(map[Ident]Meta),
		aggr:   aggr,
	}
	return metrics.NewLock(tracker)
}

type Tracker struct {
	series map[Ident]Meta
	aggr   Aggregator
}

func (t *Tracker) Consume(m metrics.Metric) {
	Samples(m, func(meta Meta, dp pmetric.NumberDataPoint) {
		id := meta.Identity()
		t.series[id] = meta
		t.aggr.Aggregate(id, dp)
	})
}

func (t *Tracker) Export() metrics.Map {
	var mm metrics.Map

	status := make(map[metrics.Ident]struct{})
	done := struct{}{}
	for id, meta := range t.series {
		if _, done := status[id.metric]; done {
			continue
		}

		res, sc, m := mm.For(id.metric)
		meta.metric.Resource().CopyTo(*res)
		meta.metric.Scope().CopyTo(*sc)
		meta.metric.CopyTo(*m)

		status[id.metric] = done
	}

	for id, meta := range t.series {
		_, _, m := mm.For(id.metric)

		fmt.Println("id is: %+v\n", id)
		fmt.Println("id.metric is: %+v\n", id.metric)
		fmt.Println("m is: %+v", m)
		fmt.Println("mm is: %+v", mm)

		// if state := reflect.ValueOf(m).FieldByName("state"); state.Elem().Int() != 0 {
		// 	panic("invalid access to shared data")
		// } // WORKS

		if mType := (*m).Type(); mType != pmetric.MetricTypeSum {
			panic(fmt.Sprintf("wrong type: %v", mType))
		} // PANICS

		dp := m.Sum().DataPoints().AppendEmpty()
		meta.attrs.CopyTo(dp.Attributes())
		t.aggr.Value(id).CopyTo(dp)
	}

	return mm
}
