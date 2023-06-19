package worker

// WorkRequest contains the type that the workers use
type WorkRequest struct {
	Line      string
	Delimiter string
	NumFields int
}
