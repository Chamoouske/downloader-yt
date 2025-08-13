package termux

import (
  "os/exec"
  "fmt"
  "downloader/internal/domain"
)

type TermuxNotifyer struct {
}

func NewTermuxNotifyer() *TermuxNotifyer {
  return &TermuxNotifyer{}
}

func (tn *TermuxNotifyer) Notify(notification domain.Notification) error {
  cmd := exec.Command("termux-notification", "--title", notification.Title, "--content", notification.Message)

  if err := cmd.Run(); err != nil {
    return fmt.Errorf("erro notifing user: %w", err)
  }

  return nil
}
