package usertask

import "context"

type ExecutionReporterStub struct {
	Reports []UserTaskExecutionReport
}

func NewExecutionReporterStub() *ExecutionReporterStub {
	return &ExecutionReporterStub{
		Reports: []UserTaskExecutionReport{},
	}
}

func (s *ExecutionReporterStub) ReportExecution(ctx context.Context, report UserTaskExecutionReport) error {
	s.Reports = append(s.Reports, report)
	return nil
}
