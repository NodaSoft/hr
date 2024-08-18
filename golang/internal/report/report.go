package report

type Report struct {
	succeeded map[int]string
	errors    map[int]string
}

func NewReport(succeeded, errors map[int]string) Report {
	return Report{
		succeeded: succeeded,
		errors:    errors,
	}
}
