package main

import (
	"golang.org/x/tools/go/analysis"

	domainlint "yoshiyoshifujii/go-ddd-sample/internal/lint/domain"
	usecaselint "yoshiyoshifujii/go-ddd-sample/internal/lint/usecase"
)

func New(conf any) ([]*analysis.Analyzer, error) {
	_ = conf
	return []*analysis.Analyzer{domainlint.Analyzer, usecaselint.Analyzer}, nil
}
