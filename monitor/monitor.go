package monitor

type Monitor struct {
	linksToProcessChan chan string
	pendingCountChan   chan int
}

func NewMonitor(pendingCountChan chan int, linksToProcessChan chan string) *Monitor {
	return &Monitor{
		pendingCountChan:   pendingCountChan,
		linksToProcessChan: linksToProcessChan,
	}
}

func (m *Monitor) Start() {
	count := 0

	for c := range m.pendingCountChan {
		count += c
		if count == 0 {
			close(m.linksToProcessChan)
			close(m.pendingCountChan)
		}
	}
}
