package run

/*
RequestTarget describe the model in the request files.
*/
type RequestTarget struct {
	ID     string
	Method string
	URL    string

	Headers map[string]string

	Body map[string]interface{}

	AWS *struct {
		Profile string
		Region  string
	}
}

/*
RunnerSpec is the specification of runner files.
It's only used to load files from the system.
*/
type RunnerSpec struct {
	Requests []*RequestTarget
}
