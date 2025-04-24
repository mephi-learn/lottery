package main

import (
	"context"
	"homework/pkg/errors"
	"homework/pkg/log/testlog"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewButler(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	log := testlog.New(t)
	ctx := context.Background()

	butler := NewManager(&ctx, log)
	require.NotNil(butler.ctx)
	require.NotNil(butler.log)
	require.NotNil(butler.quit)

	require.NotEmpty(butler.build.Name)
	require.NotEmpty(butler.build.Version)

	require.Equal(ctx, *butler.ctx)
}

func TestButlerRunSuccess(t *testing.T) {
	t.Parallel()

	stub := func() error {
		return nil
	}

	ctx := context.Background()
	butler := testButler(t, &ctx)

	go butler.run(stub)

	// Убеждаемся что мы не повисли на выходе
	<-butler.quit
}

func TestButlerRunError(t *testing.T) {
	t.Parallel()

	stub := func() error {
		return errors.New("sorry bro")
	}

	ctx := context.Background()
	butler := testButler(t, &ctx)

	go butler.run(stub)

	// Убеждаемся что мы не повисли на выходе
	<-butler.quit
}

type testShutdown struct {
	err error
}

func (s testShutdown) Shutdown(context.Context) error {
	return s.err
}

func TestButlerStopSuccess(t *testing.T) {
	t.Parallel()

	var svc testShutdown
	svc.err = nil

	ctx := context.Background()
	butler := testButler(t, &ctx)

	butler.stop(svc)
}

func TestButlerStopError(t *testing.T) {
	t.Parallel()

	var svc testShutdown
	svc.err = errors.New("i am the error")

	ctx := context.Background()
	butler := testButler(t, &ctx)

	butler.stop(svc)
}

func testButler(t *testing.T, ctx *context.Context) Manager {
	t.Helper()

	log := testlog.New(t)

	return NewManager(ctx, log)
}
