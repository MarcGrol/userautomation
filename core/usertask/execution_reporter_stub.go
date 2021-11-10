package usertask

import "context"

type ExecutionReporterStub struct {
	Reports []string
}

func NewExecutionReporterStub() *ExecutionReporterStub {
	return &ExecutionReporterStub{
		Reports: []string{},
	}
}

func (s *ExecutionReporterStub) ReportExecution(ctx context.Context, message string) {
	s.Reports = append(s.Reports, message)
}
