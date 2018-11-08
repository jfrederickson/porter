package exec

import (
	"fmt"
	"io"
	"os/exec"
)

// Exec is the logic behind the exec mixin
type Exec struct {
	Out io.Writer
	Err io.Writer
}

type Instruction struct {
	Name       string            `yaml:"name"`
	Command    string            `yaml:"command"`
	Arguments  []string          `yaml:"arguments"`
	Parameters map[string]string `yaml:"parameters"`
}

func (e *Exec) Execute(instruction Instruction) error {
	cmd := exec.Command(instruction.Command, instruction.Arguments...)
	stderr, err := cmd.StderrPipe()
	if stderr == nil || err != nil {
		return fmt.Errorf("couldnt get stderr pipe")
	}
	stdout, err := cmd.StdoutPipe()
	if stdout == nil || err != nil {
		return fmt.Errorf("couldnt get stdout pipe")
	}
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start...%s", err)
	}
	_, err = io.Copy(e.Err, stderr)
	if err != nil {
		return fmt.Errorf("error copying stderr")
	}
	_, err = io.Copy(e.Out, stdout)
	if err != nil {
		return fmt.Errorf("error copying stdout")
	}
	return nil
}