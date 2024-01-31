// Copyright 2023 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package command

import (
	"io"
	"strconv"
)

type CmdOptionFunc func(c *Command)

// WithAction set the action of the Git command, e.g. "set-url" in `git remote set-url`.
func WithAction(action string) CmdOptionFunc {
	return func(c *Command) {
		c.Action = action
	}
}

// WithFlag set optional flags to pass before positional arguments.
func WithFlag(flags ...string) CmdOptionFunc {
	return func(c *Command) {
		c.Flags = append(c.Flags, flags...)
	}
}

// WithArg add arguments that shall be passed after all flags.
func WithArg(args ...string) CmdOptionFunc {
	return func(c *Command) {
		c.Args = append(c.Args, args...)
	}
}

// WithPostSepArg set arguments that shall be passed as positional arguments after the `--`.
func WithPostSepArg(args ...string) CmdOptionFunc {
	return func(c *Command) {
		c.PostSepArgs = append(c.PostSepArgs, args...)
	}
}

// WithEnv sets environment variable using key value pair
// for example: WithEnv("GIT_TRACE", "true").
func WithEnv(keyValPairs ...string) CmdOptionFunc {
	return func(c *Command) {
		for i := 0; i < len(keyValPairs); i += 2 {
			k, v := keyValPairs[i], keyValPairs[i+1]
			c.Envs[k] = v
		}
	}
}

// WithCommitter sets given committer to the command.
func WithCommitter(name, email string) CmdOptionFunc {
	return func(c *Command) {
		c.Envs[GitCommitterName] = name
		c.Envs[GitCommitterEmail] = email
	}
}

// WithAuthor sets given author to the command.
func WithAuthor(name, email string) CmdOptionFunc {
	return func(c *Command) {
		c.Envs[GitAuthorName] = name
		c.Envs[GitAuthorEmail] = email
	}
}

// WithConfig function sets key and value for config command.
func WithConfig(key, value string) CmdOptionFunc {
	return func(c *Command) {
		c.Envs["GIT_CONFIG_KEY_"+strconv.Itoa(c.configEnvCounter)] = key
		c.Envs["GIT_CONFIG_VALUE_"+strconv.Itoa(c.configEnvCounter)] = value
		c.configEnvCounter++
		c.Envs["GIT_CONFIG_COUNT"] = strconv.Itoa(c.configEnvCounter)
	}
}

// RunOption contains option for running a command.
type RunOption struct {
	// Dir is location of repo.
	Dir string
	// Stdin is the input to the command.
	Stdin io.Reader
	// Stdout is the outputs from the command.
	Stdout io.Writer
	// Stderr is the error output from the command.
	Stderr io.Writer
}

type RunOptionFunc func(option *RunOption)

// WithDir set directory RunOption.Dir, this is repository dir
// where git command should be running.
func WithDir(dir string) RunOptionFunc {
	return func(option *RunOption) {
		option.Dir = dir
	}
}

// WithStdin set RunOption.Stdin reader.
func WithStdin(stdin io.Reader) RunOptionFunc {
	return func(option *RunOption) {
		option.Stdin = stdin
	}
}

// WithStdout set RunOption.Stdout writer.
func WithStdout(stdout io.Writer) RunOptionFunc {
	return func(option *RunOption) {
		option.Stdout = stdout
	}
}

// WithStderr set RunOption.Stderr writer.
func WithStderr(stderr io.Writer) RunOptionFunc {
	return func(option *RunOption) {
		option.Stderr = stderr
	}
}
