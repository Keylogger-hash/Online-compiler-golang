package utils

import (
	"testing"
)

func TestCompile(t *testing.T) {
	build := NewBuildResult()
	build.pathOutput = "/tmp/build/bf4e2f33-d26c-4f6d-aec8-8cef2ca67dd3/"
	build.pathCompile = build.pathOutput + "main_bf4e2f33-d26c-4f6d-aec8-8cef2ca67dd3"
	CompileCode(build)
	if build.Errors != nil {
		t.Errorf("Could not compile")
	}
}
