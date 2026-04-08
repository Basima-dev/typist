package main 

import(
	"testing"
	"github.com/stretchr/testify/assert"
) 

func TestWrapIntoLines(t *testing.T){
	lines, offests := wrapiIntoLines("hello world foo", 10)
	assert.Len(t, lines, 2)
}

