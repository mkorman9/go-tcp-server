package main

import "fmt"

const (
	UnitBps = iota
	UnitKbps
	UnitMbps
	UnitGbps
	UnitTbps
)

type throughput struct {
	value float64
	unit  int
}

func (t throughput) String() string {
	switch t.unit {
	case UnitBps:
		return fmt.Sprintf("%f bps", t.value)
	case UnitKbps:
		return fmt.Sprintf("%f kbps", t.value)
	case UnitMbps:
		return fmt.Sprintf("%f mbps", t.value)
	case UnitGbps:
		return fmt.Sprintf("%f gbps", t.value)
	case UnitTbps:
		return fmt.Sprintf("%f tbps", t.value)
	default:
		return ""
	}
}

func formatThroughput(value float64) throughput {
	unit := UnitBps

	for unit < UnitTbps {
		v := value / 1024
		if v >= 1 {
			value = v
			unit++
		} else {
			break
		}
	}

	return throughput{
		value: value * 8,
		unit:  unit,
	}
}
