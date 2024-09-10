package templatecompiler

import (
	"bytes"
	"html/template"
	"log/slog"

	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/domain"
)

type TemplateCompileInterface interface {
	Compile(parameters *domain.StatementGenerationReportParameter) (string, error)
}

func NewTemplateCompile() TemplateCompileInterface {
	return &TemplateCompile{}
}

type TemplateCompile struct {
}

func (*TemplateCompile) Compile(parameters *domain.StatementGenerationReportParameter) (string, error) {
	tmpl, err := template.ParseFiles("./templates/statement.html")
	if err != nil {
		slog.Error("Error loading template", "error", err)
		return "", err
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, parameters)
	if err != nil {
		slog.Error("Error executing template", "error", err)
		return "", err
	}

	return buffer.String(), nil
}
