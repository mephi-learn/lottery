package log

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppendFileWriterSingle(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	var opt options
	var err error

	opt.writers, err = appendFileWriter(opt.writers, os.Stdout, DestConfig{})
	require.NoError(err)
	require.Len(opt.writers, 1)

	w := opt.writers[0]
	require.IsType(fileWriter{}, w)
	require.Equal(os.Stdout, w.(fileWriter).File)
}

func TestAppendFileWriterDupe(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	var opt options
	var err error

	opt.writers, err = appendFileWriter(opt.writers, os.Stdout, DestConfig{})
	require.NoError(err)
	require.Len(opt.writers, 1)

	opt.writers, err = appendFileWriter(opt.writers, os.Stdout, DestConfig{})
	require.Error(err)
}

func TestAppendPathWriter(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	var opt options
	var err error

	opt.writers, err = appendPathWriter(opt.writers, os.DevNull, DestConfig{})
	require.NoError(err)
	require.Len(opt.writers, 1)
}

func TestFileWriterConfig(t *testing.T) {
	t.Parallel()

	config := DestConfig{Format: FormatJSON}

	fw := fileWriter{DestConfig: config}
	require.Equal(t, config, fw.Config())
}

func TestFileWriterEquals(t *testing.T) {
	t.Parallel()

	fw := fileWriter{File: os.Stdout}

	eq := fw.equals(os.Stdout)
	require.True(t, eq)
}
