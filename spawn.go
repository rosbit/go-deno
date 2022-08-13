package deno

import (
	"github.com/rosbit/go-expect"
)

func spawn(envs map[string]string, denoExePath string, args ...string) (e *expect.Expect, err error) {
	return expect.SpawnPTYWithEnvs(envs, denoExePath, args...)
}
