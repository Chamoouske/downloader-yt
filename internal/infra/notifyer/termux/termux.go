package termux

import (
	"downloader/internal/domain"
	"fmt"
	"os/exec"
)

// Commander interface to abstract exec.Command operations
type Commander interface {
	Command(name string, arg ...string) Cmd
}

// Cmd interface to abstract exec.Cmd operations
type Cmd interface {
	Run() error
}

// DefaultCommander implements Commander using os/exec
type DefaultCommander struct{}

func (c *DefaultCommander) Command(name string, arg ...string) Cmd {
	return &DefaultCmd{exec.Command(name, arg...)}
}

// DefaultCmd implements Cmd using os/exec.Cmd
type DefaultCmd struct {
	*exec.Cmd
}

func (d *DefaultCmd) Run() error {
	return d.Cmd.Run()
}

// TermuxNotifyer struct with injected dependencies
type TermuxNotifyer struct {
	commander Commander
}

func NewTermuxNotifyer(commander Commander) *TermuxNotifyer {
	return &TermuxNotifyer{commander: commander}
}

func (tn *TermuxNotifyer) Notify(notification domain.Notification) error {
	cmd := tn.commander.Command("termux-notification", "--title", notification.Title, "--content", notification.Message)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("erro notifing user: %w", err)
	}

	return nil
}
