package metrics

import "errors"

type CmdType int

const (
	INC CmdType = iota
	DEC
	SET
	CLEAR
)

//type metricType interface {
//	int | float32 | float64 | string
//}

type MetricsProto struct {
	Name  string
	Cmd   CmdType
	Value int
}

type Metrics struct {
	Mapping map[string]int
}

func (m *Metrics) Process(mp MetricsProto) {
	if v, ok := m.Mapping[mp.Name]; ok {
		switch mp.Cmd {
		case INC:
			v += 1
		case DEC:
			v -= 1
		case SET:
			v = mp.Value
		case CLEAR:
			v = 0
		}
		m.Mapping[mp.Name] = v
	}
}

func (m *Metrics) Inc(name string) {
	m.Process(MetricsProto{Name: name, Cmd: INC})
}

func (m *Metrics) Dec(name string) {
	m.Process(MetricsProto{Name: name, Cmd: DEC})
}

func (m *Metrics) Set(name string, val int) {
	m.Process(MetricsProto{Name: name, Cmd: SET, Value: val})
}

func (m *Metrics) Clear(name string) {
	m.Process(MetricsProto{Name: name, Cmd: CLEAR})
}

func (m *Metrics) GetValue(name string) (int, error) {
	if v, ok := m.Mapping[name]; ok {
		return v, nil
	} else {
		return 0, errors.New("metric name not found: " + name)
	}
}

func (m *Metrics) AddMetric(name ...string) {
	for _, n := range name {
		m.Mapping[n] = 0
	}
}

func New() *Metrics {
	return &Metrics{Mapping: make(map[string]int)}
}
