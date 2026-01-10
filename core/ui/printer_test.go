// Package ui provides console output helpers for elmos.
package ui

import (
	"os"
	"reflect"
	"testing"
)

func TestNewPrinter(t *testing.T) {
	tests := []struct {
		name string
		want *Printer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPrinter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPrinter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrinter_Success(t *testing.T) {
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		p    *Printer
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.Success(tt.args.format, tt.args.args...)
		})
	}
}

func TestPrinter_Error(t *testing.T) {
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		p    *Printer
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.Error(tt.args.format, tt.args.args...)
		})
	}
}

func TestPrinter_Warn(t *testing.T) {
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		p    *Printer
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.Warn(tt.args.format, tt.args.args...)
		})
	}
}

func TestPrinter_Info(t *testing.T) {
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		p    *Printer
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.Info(tt.args.format, tt.args.args...)
		})
	}
}

func TestPrinter_Step(t *testing.T) {
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		p    *Printer
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.Step(tt.args.format, tt.args.args...)
		})
	}
}

func TestPrinter_Print(t *testing.T) {
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		p    *Printer
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.Print(tt.args.format, tt.args.args...)
		})
	}
}

func TestPrinter_Writer(t *testing.T) {
	tests := []struct {
		name string
		p    *Printer
		want *os.File
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Writer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Printer.Writer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrintSuccess(t *testing.T) {
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintSuccess(tt.args.format, tt.args.args...)
		})
	}
}

func TestPrintError(t *testing.T) {
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintError(tt.args.format, tt.args.args...)
		})
	}
}

func TestPrintWarn(t *testing.T) {
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintWarn(tt.args.format, tt.args.args...)
		})
	}
}

func TestPrintInfo(t *testing.T) {
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintInfo(tt.args.format, tt.args.args...)
		})
	}
}

func TestPrintStep(t *testing.T) {
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintStep(tt.args.format, tt.args.args...)
		})
	}
}
