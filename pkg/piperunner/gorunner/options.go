package gorunner

type RunnerOption func(*Runner)

// WithFunctions user defined functions
func WithFunctions(functions *GoFunction) RunnerOption {
	return func(r *Runner) {
		r.functions = functions
	}
}
