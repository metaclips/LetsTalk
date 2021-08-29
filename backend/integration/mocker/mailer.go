package mocker

type Mocker struct {
	Mailer func(email, uuid string) error
}

func (m *Mocker) MailTo(email, uuid string) error {
	return m.Mailer(email, uuid)
}
