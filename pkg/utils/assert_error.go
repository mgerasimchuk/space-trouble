package utils

import "github.com/stretchr/testify/assert"

func EqualError(t assert.TestingT, expected error, actual error, msgAndArgs ...interface{}) {
	if expected != nil && actual != nil {
		assert.EqualError(t, expected, actual.Error(), msgAndArgs)
	} else {
		assert.Equal(t, expected, actual, msgAndArgs)
	}
}
