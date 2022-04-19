package item_test

import (
	"testing"

	"github.com/marcsauter/secretfs/internal/item"
	"github.com/marcsauter/secretfs/internal/secret"
	"github.com/stretchr/testify/require"
)

func TestNewItem(t *testing.T) {
	s, err := secret.New("/default/testsecret")
	require.NoError(t, err)
	require.NotNil(t, s)

	it := item.New(s, "key1", []byte("value1"))
	require.NotNil(t, it)
}

func TestAferoFileInterface(t *testing.T) {
}