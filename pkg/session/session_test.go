package session

import (
	"strconv"
	"testing"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/stretchr/testify/require"
)

func TestOK(t *testing.T) {
	tests := []struct {
		length int
	}{
		{32},
		{64},
		{128},
		{256},
	}
	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.length), func(t *testing.T) {
			t.Parallel()
			sessionID, err := GenerateSessionID(tt.length)
			require.NoError(t, err)
			require.NotEmpty(t, sessionID)
		})
	}
}

func TestLong(t *testing.T) {
	tests := []struct {
		length int
	}{
		{3231423},
		{32413123412},
		{513},
		{1000},
	}
	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.length), func(t *testing.T) {
			t.Parallel()
			_, err := GenerateSessionID(tt.length)
			require.Equal(t, err.Error(), errs.ErrMsgLengthTooLong)
		})
	}
}

func TestShort(t *testing.T) {
	tests := []struct {
		length int
	}{
		{0},
		{1},
		{5},
		{7},
	}
	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.length), func(t *testing.T) {
			t.Parallel()
			_, err := GenerateSessionID(tt.length)
			require.Equal(t, err.Error(), errs.ErrMsgLengthTooShort)
		})
	}
}

func TestNegative(t *testing.T) {
	tests := []struct {
		length int
	}{
		{-1},
		{-8},
		{-384673198276},
		{-124112},
	}
	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.length), func(t *testing.T) {
			t.Parallel()
			_, err := GenerateSessionID(tt.length)
			require.Equal(t, err.Error(), errs.ErrMsgNegativeSessionIDLength)
		})
	}
}
