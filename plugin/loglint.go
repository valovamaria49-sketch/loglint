//go:build plugin
// +build plugin

package main

import (
	myanalysis "github.com/valovamaria49-sketch/loglint/analysis"
	"golang.org/x/tools/go/analysis"
)

func New(conf any) ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{myanalysis.Analyzer}, nil
}
