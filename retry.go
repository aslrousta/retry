package retry

// Option defines a retry option.
type Option interface {
	option()
}

// MaxTries is the option for setting max number of tries.
type MaxTries int

// If is the option for setting a condition for errors which imply a retry.
type If func(error) bool

func (MaxTries) option() {}
func (If) option()       {}

// Func is the type of a retryable function.
type Func func() error

// Retry attempts to perform f with the conditions specified by the given
// options. It returns immediately if f returns nil at any attempt, or returns
// the error returned by f if no further attempts can be done.
//
// If no MaxTries option is specified or it is a negative number, 3 attempts
// will be performed by default. A zero or one results no retry.
//
// If no If option is specified, any non-nil error implies a retry.
//
// It panics if f is nil.
func Retry(f Func, opts ...Option) error {
	if f == nil {
		panic("f is nil")
	}

	maxTries := 3
	implies := func(error) bool { return true }

	for _, opt := range opts {
		switch v := opt.(type) {
		case MaxTries:
			if v > 1 {
				maxTries = int(v)
			} else if v >= 0 {
				maxTries = 1
			}
		case If:
			if v != nil {
				implies = v
			}
		}
	}

	var err error

	for try := 0; try < maxTries; try++ {
		if err = f(); err == nil || !implies(err) {
			break
		}
	}

	return err
}
