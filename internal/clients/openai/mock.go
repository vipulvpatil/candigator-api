package openai

type MockClientSuccess struct {
	Text string
}

func (m *MockClientSuccess) CallCompletionApi(prompt string) (string, error) {
	return m.Text, nil
}
