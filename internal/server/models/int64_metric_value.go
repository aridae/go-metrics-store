package models

import (
	"fmt"
)

type int64MetricValue struct {
	Val int64
}

func NewInt64MetricValue(val int64) MetricValue {
	return int64MetricValue{Val: val}
}

func (mv int64MetricValue) String() string {
	return fmt.Sprintf("%d", mv.Val)
}

func (mv int64MetricValue) Inc(v MetricValue) (MetricValue, error) {
	int64Val, ok := v.(int64MetricValue)
	if !ok {
		return nil, fmt.Errorf("expected int64 metric value, got %T", v)
	}

	newVal := mv.Val + int64Val.Val

	return int64MetricValue{Val: newVal}, nil
}

func (mv int64MetricValue) UnsafeCastFloat() float64 {
	return float64(mv.Val)
}

func (mv int64MetricValue) UnsafeCastInt() int64 {
	return mv.Val
}
