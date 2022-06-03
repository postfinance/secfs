// Package io provides basic interfaces to I/O primitives.
package io

// Loader is the interface that wraps the Load method,
type Loader interface {
	Load(Secreter) error
}

// Storer is the interface that wraps the Store method.
type Storer interface {
	Store(Secreter) error
}

// Deleter is the interface that wraps the Delete method.
type Deleter interface {
	Delete(Secreter) error
}

// LoadStoreDeleter is the interface that groups the basic Load, Store and Delete methods.
type LoadStoreDeleter interface {
	Loader
	Storer
	Deleter
}

// Secreter is the interface that abstracts the Kubernetes secret
type Secreter interface {
	Name() string
	Namespace() string
	Data() map[string][]byte
	SetData(map[string][]byte)
}
