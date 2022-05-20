package mock

type ManagerMock struct {
	aliases map[string]string
}

func NewManagerMock() *ManagerMock {
	return &ManagerMock{make(map[string]string)}
}

func (m *ManagerMock) set(name, value string) {
	m.aliases[name] = value
}

func (m *ManagerMock) Load() (map[string]string, error) {
	return m.aliases, nil
}

func (m *ManagerMock) Save(aliases map[string]string) error {
	m.aliases = aliases
	return nil
}
