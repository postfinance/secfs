package secfs_test

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"sort"
	"syscall"
	"testing"

	"github.com/marcsauter/secfs"
	"github.com/marcsauter/secfs/internal/backend"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

// TODO: add afero.File tests

func TestCreateOpenClose(t *testing.T) {
	namespace := "default"
	secret := "testsecret1"
	key := "testfile1"

	filename := path.Join(namespace, secret, key)
	secretname := path.Join(namespace, secret)

	cs := backend.NewFakeClientset()

	// prepare
	sfs := secfs.New(cs)
	err := sfs.Mkdir(secretname, os.FileMode(0))
	require.NoError(t, err)

	b := backend.New(cs)

	t.Run("FileCreate on secret", func(t *testing.T) {
		f, err := secfs.FileCreate(b, secretname)
		require.ErrorIs(t, err, syscall.EISDIR)
		require.Nil(t, f)
	})

	t.Run("FileCreate on file", func(t *testing.T) {
		f, err := secfs.FileCreate(b, filename)
		require.NoError(t, err)
		require.NotNil(t, f)

		// interface backend.Metadata
		require.Equal(t, namespace, f.Namespace())
		require.Equal(t, secret, f.Secret())
		require.Equal(t, key, f.Key())

		// interface os.FileInfo
		require.Equal(t, key, f.Name())
		require.Equal(t, int64(1), f.Size())
		require.Equal(t, fs.FileMode(0), f.Mode())
		require.False(t, f.ModTime().IsZero())
		require.False(t, f.IsDir())
		require.Equal(t, f, f.Sys())

		require.NoError(t, f.Close())
		require.ErrorIs(t, f.Close(), afero.ErrFileClosed)
	})

	t.Run("Open file", func(t *testing.T) {
		f, err := secfs.Open(b, filename)
		require.NoError(t, err)
		require.NotNil(t, f)

		// interface backend.Metadata
		require.Equal(t, namespace, f.Namespace())
		require.Equal(t, secret, f.Secret())
		require.Equal(t, key, f.Key())

		// interface os.FileInfo
		require.Equal(t, key, f.Name())
		require.Equal(t, int64(1), f.Size())
		require.Equal(t, fs.FileMode(0), f.Mode())
		require.False(t, f.ModTime().IsZero())
		require.False(t, f.IsDir())
		require.Equal(t, f, f.Sys())

		require.NoError(t, f.Close())
		require.ErrorIs(t, f.Close(), afero.ErrFileClosed)
	})

	t.Run("Open secret", func(t *testing.T) {
		f, err := secfs.Open(b, secretname)
		require.NoError(t, err)
		require.NotNil(t, f)

		// interface backend.Metadata
		require.Equal(t, namespace, f.Namespace())
		require.Equal(t, secret, f.Secret())
		require.Empty(t, f.Key())

		// interface os.FileInfo
		require.Equal(t, secret, f.Name())
		require.Equal(t, int64(1), f.Size())
		require.Equal(t, fs.ModeDir, f.Mode())
		require.False(t, f.ModTime().IsZero())
		require.True(t, f.IsDir())
		require.Equal(t, f, f.Sys())

		require.ErrorIs(t, f.Close(), syscall.EISDIR)
	})
}

func TestReadSeekWriteSyncTruncateSecret(t *testing.T) {
	namespace := "default"
	secret := "testsecret2"

	secretname := path.Join(namespace, secret)

	cs := backend.NewFakeClientset()

	// prepare
	sfs := secfs.New(cs)
	err := sfs.Mkdir(secretname, os.FileMode(0))
	require.NoError(t, err)

	b := backend.New(cs)

	t.Run("Read", func(t *testing.T) {
		f, err := secfs.Open(b, secretname)
		require.NoError(t, err)
		require.NotNil(t, f)

		n, err := f.Read([]byte{})
		require.Zero(t, n)
		require.ErrorIs(t, err, syscall.EISDIR)

		n, err = f.ReadAt([]byte{}, 10)
		require.Zero(t, n)
		require.ErrorIs(t, err, syscall.EISDIR)
	})

	t.Run("Seek", func(t *testing.T) {
		f, err := secfs.Open(b, secretname)
		require.NoError(t, err)
		require.NotNil(t, f)

		n, err := f.Seek(10, io.SeekStart)
		require.Zero(t, n)
		require.ErrorIs(t, err, syscall.EISDIR)
	})

	t.Run("Write", func(t *testing.T) {
		f, err := secfs.Open(b, secretname)
		require.NoError(t, err)
		require.NotNil(t, f)

		n, err := f.Write([]byte{})
		require.Zero(t, n)
		require.ErrorIs(t, err, syscall.EISDIR)

		n, err = f.WriteAt([]byte{}, 10)
		require.Zero(t, n)
		require.ErrorIs(t, err, syscall.EISDIR)
	})

	t.Run("Sync", func(t *testing.T) {
		f, err := secfs.Open(b, secretname)
		require.NoError(t, err)
		require.NotNil(t, f)

		err = f.Sync()
		require.ErrorIs(t, err, syscall.EISDIR)
	})

	t.Run("Truncate", func(t *testing.T) {
		f, err := secfs.Open(b, secretname)
		require.NoError(t, err)
		require.NotNil(t, f)

		err = f.Truncate(10)
		require.ErrorIs(t, err, syscall.EISDIR)
	})

	t.Run("WriteString", func(t *testing.T) {
		f, err := secfs.Open(b, secretname)
		require.NoError(t, err)
		require.NotNil(t, f)

		n, err := f.WriteString("")
		require.Zero(t, n)
		require.ErrorIs(t, err, syscall.EISDIR)
	})
}

func TestReadSeekWriteSyncTruncateFile(t *testing.T) {
	namespace := "default"
	secret := "testsecret2"
	key := "testfile2"

	filename := path.Join(namespace, secret, key)
	secretname := path.Join(namespace, secret)

	cs := backend.NewFakeClientset()

	// prepare
	sfs := secfs.New(cs)
	err := sfs.Mkdir(secretname, os.FileMode(0))
	require.NoError(t, err)

	b := backend.New(cs)

	f, err := secfs.FileCreate(b, filename)
	require.NoError(t, err)
	require.NotNil(t, f)

	t.Run("Read", func(t *testing.T) {
		f, err := secfs.Open(b, filename)
		require.NoError(t, err)
		require.NotNil(t, f)

		n, err := f.Read([]byte{})
		require.Zero(t, n)
		require.ErrorIs(t, err, io.EOF)

		n, err = f.ReadAt([]byte{}, 10)
		require.Zero(t, n)
		require.ErrorIs(t, err, io.EOF)
	})

	t.Run("Seek", func(t *testing.T) {
		f, err := secfs.Open(b, filename)
		require.NoError(t, err)
		require.NotNil(t, f)

		n, err := f.Seek(10, io.SeekStart)
		require.Equal(t, int64(10), n)
		require.NoError(t, err)
	})

	t.Run("Write", func(t *testing.T) {
		f, err := secfs.Open(b, filename)
		require.NoError(t, err)
		require.NotNil(t, f)

		n, err := f.Write([]byte{})
		require.Zero(t, n)
		require.NoError(t, err)

		n, err = f.WriteAt([]byte{}, 10)
		require.Zero(t, n)
		require.NoError(t, err)
	})

	t.Run("Sync", func(t *testing.T) {
		f, err := secfs.Open(b, filename)
		require.NoError(t, err)
		require.NotNil(t, f)

		err = f.Sync()
		require.NoError(t, err)
	})

	t.Run("Truncate", func(t *testing.T) {
		f, err := secfs.Open(b, filename)
		require.NoError(t, err)
		require.NotNil(t, f)

		err = f.Truncate(10)
		require.NoError(t, err)
	})

	t.Run("WriteString", func(t *testing.T) {
		f, err := secfs.Open(b, filename)
		require.NoError(t, err)
		require.NotNil(t, f)

		n, err := f.WriteString("")
		require.Zero(t, n)
		require.NoError(t, err)
	})
}

func TestReaddir(t *testing.T) {
	namespace := "default"
	secret := "testsecret3"
	key := "testfile3"

	secretname := path.Join(namespace, secret)

	cs := backend.NewFakeClientset()

	// prepare
	sfs := secfs.New(cs)
	err := sfs.Mkdir(secretname, os.FileMode(0))
	require.NoError(t, err)

	b := backend.New(cs)

	count := 10

	for i := 0; i < count; i++ {
		filename := path.Join(namespace, secret, fmt.Sprintf("%s%d", key, i))
		f, err := secfs.FileCreate(b, filename)
		require.NoError(t, err)
		require.NotNil(t, f)
	}

	t.Run("Readdir not dir", func(t *testing.T) {
		f, err := secfs.Open(b, path.Join(namespace, secret, key+"0"))
		require.NoError(t, err)
		require.NotNil(t, f)

		fi, err := f.Readdir(0)
		require.ErrorIs(t, err, syscall.ENOTDIR)
		require.Len(t, fi, 0)
	})

	t.Run("Readdir", func(t *testing.T) {
		f, err := secfs.Open(b, secretname)
		require.NoError(t, err)
		require.NotNil(t, f)

		fi, err := f.Readdir(count)
		require.NoError(t, err)
		require.Len(t, fi, count)

		n := make([]string, count)
		for i := range fi {
			n[i] = fi[i].Name()
		}

		sort.Strings(n)

		for i := 0; i < count; i++ {
			require.Equal(t, n[i], fmt.Sprintf("%s%d", key, i))
		}

	})

	t.Run("Readnames not dir", func(t *testing.T) {
		f, err := secfs.Open(b, path.Join(namespace, secret, key+"0"))
		require.NoError(t, err)
		require.NotNil(t, f)

		n, err := f.Readdirnames(0)
		require.ErrorIs(t, err, syscall.ENOTDIR)
		require.Len(t, n, 0)
	})

	t.Run("Readnames", func(t *testing.T) {
		f, err := secfs.Open(b, secretname)
		require.NoError(t, err)
		require.NotNil(t, f)

		n, err := f.Readdirnames(count)
		require.NoError(t, err)
		require.Len(t, n, count)

		sort.Strings(n)

		for i := 0; i < count; i++ {
			require.Equal(t, n[i], fmt.Sprintf("%s%d", key, i))
		}
	})
}
