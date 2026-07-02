package md2pdf_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	fakeconverter "github.com/damianoneill/md2pdf/adapter/fakeconverter"
	fakerenderer "github.com/damianoneill/md2pdf/adapter/fakerenderer"
	"github.com/damianoneill/md2pdf/usecase/md2pdf"
)

func TestConvertUseCase_Convert(t *testing.T) {
	src := []byte("hello")
	htmlOut := []byte("<p>hello</p>")

	tests := []struct {
		name       string
		convOut    []byte
		convErr    error
		rendOut    []byte
		rendErr    error
		wantOut    string
		wantConvIn []byte
		wantRendIn []byte
		wantErr    bool
	}{
		{
			name:       "success: src reaches converter, converter output reaches renderer",
			convOut:    htmlOut,
			rendOut:    []byte("%PDF-1.4"),
			wantOut:    "%PDF-1.4",
			wantConvIn: src,
			wantRendIn: htmlOut,
		},
		{
			name:    "converter error: propagated, renderer not called",
			convErr: errors.New("bad markdown"),
			wantErr: true,
		},
		{
			name:    "renderer error: propagated",
			convOut: htmlOut,
			rendErr: errors.New("render failed"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conv := &fakeconverter.Converter{Output: tt.convOut, Err: tt.convErr}
			rend := &fakerenderer.Renderer{Output: tt.rendOut, Err: tt.rendErr}
			svc := md2pdf.NewConvertUseCase(conv, rend)

			got, err := svc.Convert(context.Background(), src)
			if (err != nil) != tt.wantErr {
				t.Fatalf("wantErr=%v, got err=%v", tt.wantErr, err)
			}
			if tt.wantErr {
				return
			}
			if string(got) != tt.wantOut {
				t.Fatalf("output: want %q, got %q", tt.wantOut, got)
			}
			if !bytes.Equal(conv.Input, tt.wantConvIn) {
				t.Fatalf("converter input: want %q, got %q", tt.wantConvIn, conv.Input)
			}
			if !bytes.Equal(rend.Input, tt.wantRendIn) {
				t.Fatalf("renderer input: want %q, got %q", tt.wantRendIn, rend.Input)
			}
		})
	}
}
