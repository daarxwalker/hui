package headlessui

import (
	"context"
	"math"

	"github.com/a-h/templ"
)

type MeterProps struct {
	Base
	Value     float64
	Min       float64
	Max       float64
	Low       float64
	High      float64
	Optimum   string
	Label     string
	ValueText string
}

func percentOf(value, low, high float64) string {
	if high <= low {
		return "0%"
	}
	clamped := min(max(value, low), high)
	return formatNumber(math.Round((clamped-low)/(high-low)*100000)/1000) + "%"
}

func labelledBy(ctx context.Context, own templ.Attributes, label string) string {
	if label != "" {
		own["aria-label"] = label
	}
	f, inField := fieldFrom(ctx)
	if !inField {
		return ""
	}
	own["aria-labelledby"] = f.ID + "-label"
	if by := f.describedBy(); by != "" {
		own["aria-describedby"] = by
	}
	return f.ID
}

func meterRange(p MeterProps) (low, high float64) {
	low, high = p.Min, p.Max
	if low == 0 && high == 0 {
		high = 100
	}
	return low, high
}

func meterLevel(p MeterProps, value, low, high float64) string {
	lowBound, highBound := low, high
	if p.Low != 0 {
		lowBound = min(max(p.Low, low), high)
	}
	if p.High != 0 {
		highBound = min(max(p.High, lowBound), high)
	}
	switch {
	case value < lowBound:
		return "low"
	case value > highBound:
		return "high"
	default:
		return "medium"
	}
}

func meterState(optimum, level string) string {
	switch optimum {
	case "low", "high":
		far := map[string]string{"low": "high", "high": "low"}[optimum]
		switch level {
		case optimum:
			return "optimum"
		case far:
			return "critical"
		default:
			return "suboptimal"
		}
	case "medium":
		if level == "medium" {
			return "optimum"
		}
		return "suboptimal"
	default:
		return ""
	}
}

func meterAttrs(ctx context.Context, p MeterProps) templ.Attributes {
	low, high := meterRange(p)
	value := min(max(p.Value, low), high)

	own := make(templ.Attributes, 12)
	own["role"] = "meter"
	own["aria-valuemin"] = formatNumber(low)
	own["aria-valuemax"] = formatNumber(high)
	own["aria-valuenow"] = formatNumber(value)
	if id := labelledBy(ctx, own, p.Label); id != "" {
		p.ID = id
	}
	if p.ValueText != "" {
		own["aria-valuetext"] = p.ValueText
	}
	own["style"] = "--hui-meter: " + percentOf(value, low, high)

	level := meterLevel(p, value, low, high)
	own["data-level"] = level
	if state := meterState(p.Optimum, level); state != "" {
		own["data-state"] = state
	}
	return p.merge(own)
}
