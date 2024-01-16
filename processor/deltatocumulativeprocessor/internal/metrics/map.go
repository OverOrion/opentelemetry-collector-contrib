package metrics

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

type Map struct {
	resources map[ResourceIdent]pcommon.Resource
	scopes    map[ScopeIdent]pcommon.InstrumentationScope
	metrics   map[Ident]pmetric.Metric
}

func (mm *Map) For(id Ident) (*pcommon.Resource, *pcommon.InstrumentationScope, *pmetric.Metric) {
	if mm.resources == nil {
		mm.resources = make(map[ResourceIdent]pcommon.Resource)
	}

	if mm.scopes == nil {
		mm.scopes = make(map[ScopeIdent]pcommon.InstrumentationScope)

	}

	if mm.metrics == nil {
		mm.metrics = make(map[Ident]pmetric.Metric)
	}

	res, ok := mm.resources[id.ResourceIdent]
	if !ok {
		res = pcommon.NewResource()
		mm.resources[id.ResourceIdent] = res
	}

	sc, ok := mm.scopes[id.ScopeIdent]
	if !ok {
		sc = pcommon.NewInstrumentationScope()
		mm.scopes[id.ScopeIdent] = sc
	}

	m, ok := mm.metrics[id]
	if !ok {
		// metricsSlice := pmetric.NewMetricSlice()
		// metricsSlice.AppendEmpty()
		// metaMetric := metricsSlice.At(0)

		metaMetric := pmetric.NewMetric()

		m = metaMetric
		mm.metrics[id] = m
	}

	return &res, &sc, &m
}

func (mm Map) Merge() pmetric.Metrics {
	metrics := pmetric.NewMetrics()

	rms := make(map[ResourceIdent]pmetric.ResourceMetrics)
	for id, res := range mm.resources {
		rm := metrics.ResourceMetrics().AppendEmpty()
		res.CopyTo(rm.Resource())
		rms[id] = rm
	}

	sms := make(map[ScopeIdent]pmetric.ScopeMetrics)
	for id, sc := range mm.scopes {
		sm := rms[id.ResourceIdent].ScopeMetrics().AppendEmpty()
		sc.CopyTo(sm.Scope())
		sms[id] = sm
	}

	for id, m := range mm.metrics {
		metric := sms[id.ScopeIdent].Metrics().AppendEmpty()
		m.CopyTo(metric)
	}

	return metrics
}
