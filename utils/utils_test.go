package utils

import (
	"testing"
)

func TestCompile(t *testing.T) {
	build := NewBuildResult()
	build.pathOutput = "/tmp/build/0171e548-e67d-4219-808a-a53ec0785213"
	build.pathCompile = build.pathOutput + "0171e548-e67d-4219-808a-a53ec0785213"
	CompileCode(build)
	if build.Errors != nil {
		t.Errorf("Could not compile")
	}
}
