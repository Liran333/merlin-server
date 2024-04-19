package primitive

const (
	statusPassed     = "passed"
	statusProcessing = "processing"
	statusFailed     = "failed"
)

type CertificateStatus interface {
	CertificateStatus() string
}

func NewProcessingStatus() CertificateStatus {
	return certificateStatus(statusProcessing)
}

func NewPassedStatus() CertificateStatus {
	return certificateStatus(statusPassed)
}

func CreateCertificateStatus(v string) CertificateStatus {
	return certificateStatus(v)
}

type certificateStatus string

func (c certificateStatus) CertificateStatus() string {
	return string(c)
}
