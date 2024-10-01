package msgtypes

type ErrMsg struct {
	Err error
}

func (e ErrMsg) Error() string {
	return e.Err.Error()
}

func NewErrMsg(err error) ErrMsg {
	return ErrMsg{Err: err}
}

type NotificationMsg struct {
	Notification string
}

func NewNotificationMsg(notification string) NotificationMsg {
	return NotificationMsg{Notification: notification}
}
