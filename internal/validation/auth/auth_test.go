package auth

import (
  "testing"

  errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
  "github.com/stretchr/testify/require"
)

func TestOKAuth(t *testing.T) {
  tests := []struct {
    password string
  }{
    {`.(FdeZO7`},
    {`=Ix7U!Kvk=8P`},
    {`g5~(Lh/Y<4Yz.PJXu`},
    {`V_0m>04–w@Q{x%mR`},
  }
  for _, tt := range tests {
    t.Run(tt.password, func(t *testing.T) {
      t.Parallel()
      err := IsValidPassword(tt.password)
      require.NoError(t, err)
    })
  }
}

func TestShortPassword(t *testing.T) {
  tests := []struct {
    password string
  }{
    {`O\Rx2`},
    {`=qX1*`},
    {`Y%5i`},
    {`gQ!0`},
  }
  for _, tt := range tests {
    t.Run(tt.password, func(t *testing.T) {
      t.Parallel()
      err := IsValidPassword(tt.password)
      require.Error(t, err)
      require.Equal(t, err.Error(), errs.ErrPasswordTooShort)
    })
  }
}

func TestLongPassword(t *testing.T) {
  tests := []struct {
    password string
  }{
    {`B}Run:yarlpeO\=RVFMM5.[vG]`},
    {`?lfJ=T#mb6EGoI5W\Yqwp59,YF{}<{St60`},
    {`bw=fb\QM&+qpLt19}[#q[TQiO~–:#{;V*iPsvbi},<`},
    {`=+;h$)7\Qwt2/fP(c6{1F^sIybJcf,e*;q2ujrZVA{PH2–sd]j`},
  }
  for _, tt := range tests {
    t.Run(tt.password, func(t *testing.T) {
      t.Parallel()
      err := IsValidPassword(tt.password)
      require.Error(t, err)
      require.Equal(t, err.Error(), errs.ErrPasswordTooLong)
    })
  }
}

func TestEmptyPassword(t *testing.T) {
  tests := []struct {
    password string
  }{
    {`          `},
    {`           `},
    {`            `},
    {`             `},
  }
  for _, tt := range tests {
    t.Run(tt.password, func(t *testing.T) {
      t.Parallel()
      err := IsValidPassword(tt.password)
      require.Error(t, err)
      require.Equal(t, err.Error(), errs.ErrEmptyPassword)
    })
  }
}

func TestValidEmail(t *testing.T) {
  tests := []struct {
    email string
  }{
    {"test@example.com"},
    {"user.tag+tag@example.com"},
    {"user@sub.example.com"},
    {"validemail123@mail.ru"},
  }
  for _, tt := range tests {
    t.Run(tt.email, func(t *testing.T) {
      t.Parallel()
      err := IsValidEmail(tt.email)
      require.NoError(t, err)
    })
  }
}

func TestInvalidEmail(t *testing.T) {
  tests := []struct {
    email string
  }{
    {"invalid-email"},
    {"invalid@.com"},
    {"invalid@com"},
    {"invalid@"},
    {"invalid@com."},
    {"invalid@com.."},
    {"invalid@com "},
    {"invalid@com("},
    {"invalid@com,com"},
    {"invalid@com;com"},
    {"invalid@com:com"},
  }
  for _, tt := range tests {
    t.Run(tt.email, func(t *testing.T) {
      t.Parallel()
      err := IsValidEmail(tt.email)
      require.Error(t, err)
      require.Equal(t, err.Error(), errs.ErrInvalidEmail)
    })
  }
}

func TestEmptyEmail(t *testing.T) {
  tests := []struct {
    email string
  }{
    {""},
    {" "},
    {"  "},
    {"   "},
  }
  for _, tt := range tests {
    t.Run(tt.email, func(t *testing.T) {
      t.Parallel()
      err := IsValidEmail(tt.email)
      require.Error(t, err)
      require.Equal(t, err.Error(), errs.ErrInvalidEmail)
    })
  }
}

func TestEmailWithoutAtSymbol(t *testing.T) {
  tests := []struct {
    email string
  }{
    {"invalidemail.com"},
    {"invalidemail"},
    {"invalid.email.com"},
  }
  for _, tt := range tests {
    t.Run(tt.email, func(t *testing.T) {
      t.Parallel()
      err := IsValidEmail(tt.email)
      require.Error(t, err)
      require.Equal(t, err.Error(), errs.ErrInvalidEmail)
    })
  }
}

func TestEmailWithMultipleAtSymbols(t *testing.T) {
  tests := []struct {
    email string
  }{
    {"invalid@@example.com"},
    {"invalid@example@com"},
    {"invalid@example@com@"},
  }
  for _, tt := range tests {
    t.Run(tt.email, func(t *testing.T) {
      t.Parallel()
      err := IsValidEmail(tt.email)
      require.Error(t, err)
      require.Equal(t, err.Error(), errs.ErrInvalidEmail)
    })
  }
}
