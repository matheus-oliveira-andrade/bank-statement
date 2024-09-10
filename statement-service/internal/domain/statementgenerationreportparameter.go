package domain

type StatementGenerationReportParameter struct {
	Document     string
	CustomerName string
	Movements    []MovementReportParameter
}
