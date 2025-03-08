package config

import (
  "context"
  "testing"

  "github.com/stretchr/testify/require"
)

func TestOkServer(t *testing.T) {
  cfg, err := New()
  require.NoError(t, err)
  require.NotNil(t, cfg)
  ctx := WrapServerContext(context.Background(), &cfg.Server)
  res := FromServerContext(ctx)
  require.Equal(t, &cfg.Server, res)
}

func TestFailServer(t *testing.T) {
  cfg, err := New()
  require.NoError(t, err)
  require.NotNil(t, cfg)
  ctx := WrapServerContext(context.Background(), cfg.Server)
  res := FromServerContext(ctx)
  require.Nil(t, res)
}

func TestOkCookie(t *testing.T) {
  cfg, err := New()
  require.NoError(t, err)
  require.NotNil(t, cfg)
  ctx := WrapCookieContext(context.Background(), &cfg.Cookie)
  res := FromCookieContext(ctx)
  require.Equal(t, &cfg.Cookie, res)
}

func TestFailCookie(t *testing.T) {
  cfg, err := New()
  require.NoError(t, err)
  require.NotNil(t, cfg)
  ctx := WrapCookieContext(context.Background(), cfg.Cookie)
  res := FromCookieContext(ctx)
  require.Nil(t, res)
}
