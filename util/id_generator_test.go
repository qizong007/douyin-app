package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateId(t *testing.T) {
	InitIdGenerator()
	id1 := GenerateId()
	id2 := GenerateId()
	assert.Equal(t, int64(0), sequenceIdStepMask&id1)
	assert.Equal(t, int64(1), sequenceIdStepMask&id2)
	assert.Equal(t, true, (id1>>sequenceIdStep) == (id2>>sequenceIdStep))
}
