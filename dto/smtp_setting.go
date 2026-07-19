package dto

type SMTPSetting struct {
	From         string
	FromName     string
	SMTPHost     string
	SMTPPort     int
	SMTPUserName string
	SMTPPassword string
}

type SMTPMessage struct {
	To      string
	Subject string
	Html    string
	Plain   string
	CC      []string
}
