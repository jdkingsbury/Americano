package msgtypes

type ErrMsg struct {
  err error
}

func (e ErrMsg) Error() string {
  return e.err.Error()
}

func NewErrMsg(err error) ErrMsg {
  return ErrMsg{err: err}
}
