package domain

type Notifyer interface {
  Notify(notification Notification) error
}
