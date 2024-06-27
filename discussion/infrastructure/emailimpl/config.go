package emailimpl

type Config struct {
	RootUrl             string   `json:"root_url" required:"true"`
	ReportTitle         string   `json:"report_title" required:"true"`
	ReportEmailReceiver []string `json:"report_email_receiver" required:"true"`
}
