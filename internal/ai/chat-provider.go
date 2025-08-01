package ai

type ChatResponse struct {
	Answer     string
	Successful bool
}

type ChatProvider interface {
	Ask(request string) (ChatResponse, error)
}
