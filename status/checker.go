package status

type Checker interface {
	CheckStatus() Status
}
