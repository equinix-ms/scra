package audit

type Auditor interface {
	AuditOnce() error
}
