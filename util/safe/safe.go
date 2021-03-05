package safe

import (
	"errors"
	"fmt"
)

// Do solves panic automatically, converts it to an error and returns it.
func Do(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = fmt.Errorf("%#v", x)
			}
		}
	}()
	err = fn()
	return err
}

// Func returns the method to be executed safely.
func Func(fn func() error) func() error {
	return func() error {
		return Do(fn)
	}
} 
