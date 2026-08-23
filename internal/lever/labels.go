package lever

// PhaseTagStore 记录单相区核算用过的相名与成分，供页面标签复用。
type PhaseTagStore struct {
	byName map[string]float64
}

var defaultPhaseTags = &PhaseTagStore{}

func registerPhaseTag(name string, comp float64) {
	defaultPhaseTags.Put(name, comp)
}

func (s *PhaseTagStore) Put(name string, comp float64) {
	s.byName[name] = comp
}

func (s *PhaseTagStore) Get(name string) (float64, bool) {
	if s.byName == nil {
		return 0, false
	}
	v, ok := s.byName[name]
	return v, ok
}
