package kv

// Wrapper swaps a modified environment onto the process, runs your functions,
// and then puts the original environment back when it's done.
type Wrapper struct {
	original Environment
	modified Environment
}

// WrapCalls drops modified onto the process environment, runs fns in order, and
// then restores the original environment. It bails at the first error.
func WrapCalls(modified Environment, fns ...func() error) error {
	return NewWrapper(modified).CallWithRestore(fns...)
}

// NewWrapper takes a snapshot of the current environment and then lays modified on top.
func NewWrapper(modified Environment) Wrapper {
	return Wrapper{
		original: New(),
		modified: modified.Set(),
	}
}

// CallWithRestore runs fns and then puts the original environment back afterward.
func (w Wrapper) CallWithRestore(fns ...func() error) error {
	defer w.Restore()
	return w.Call(fns...)
}

// Call runs fns in order, stopping at and handing back the first error. Any nil
// functions just get skipped.
func (w Wrapper) Call(fns ...func() error) error {
	for _, fn := range fns {
		if fn == nil {
			continue
		}
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

// Restore clears out the modified environment and lays the original back down.
func (w Wrapper) Restore() {
	w.modified.Unset()
	w.original.Set()
}
