package chart

// SizedLineChart represents a chart with a fixed size and a list of float64 values.
type SizedLineChart struct {
	Values []float64
	size   int
}

// NewSizedLineChart creates a new SizedLineChart with the specified size.
func NewSizedLineChart(size int) SizedLineChart {
	if size < 1 {
		size = 1
	}

	return SizedLineChart{
		Values: make([]float64, size),
		size:   size,
	}
}

// Add appends a new value to the chart. If the chart is at full capacity, it removes the oldest value.
func (s *SizedLineChart) Add(value float64) []float64 {
	if len(s.Values) >= s.size {
		s.Values = s.Values[1:]
	}
	s.Values = append(s.Values, value)
	return s.Values
}
