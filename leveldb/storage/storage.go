// Copyright (c) 2012, Suryandaru Triandana <syndtr@gmail.com>
// All rights reserved.
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package storage provides storage abstraction for LevelDB.
package storage

import (
	"errors"
	"io"

	"github.com/jrallison/goleveldb/leveldb/util"
)

type FileType uint32

const (
	TypeManifest FileType = 1 << iota
	TypeJournal
	TypeTable

	TypeAll = TypeManifest | TypeJournal | TypeTable
)

func (t FileType) String() string {
	switch t {
	case TypeManifest:
		return "manifest"
	case TypeJournal:
		return "journal"
	case TypeTable:
		return "table"
	}
	return "<unknown>"
}

var (
	ErrInvalidFile = errors.New("leveldb/storage: invalid file for argument")
	ErrLocked      = errors.New("leveldb/storage: already locked")
	ErrClosed      = errors.New("leveldb/storage: closed")
)

// Syncer is the interface that wraps basic Sync method.
type Syncer interface {
	// Sync commits the current contents of the file to stable storage.
	Sync() error
}

// Reader is the interface that groups the basic Read, Seek, ReadAt and Close
// methods.
type Reader interface {
	io.ReadSeeker
	io.ReaderAt
	io.Closer
}

// Writer is the interface that groups the basic Write, Sync and Close
// methods.
type Writer interface {
	io.WriteCloser
	Syncer
}

// File is the file.
type File interface {
	// Open opens the file for read. Returns os.ErrNotExist error
	// if the file does not exist.
	// Open returns error if the underlying storage is closed.
	Open() (r Reader, err error)

	// Create creates the file for writting. Truncate the file if
	// already exist.
	// Returns error if the underlying storage is closed.
	Create() (w Writer, err error)

	// Type returns the file type
	Type() FileType

	// Num returns the file number.
	Num() uint64

	// Remove removes the file.
	// Returns error if the underlying storage is closed.
	Remove() error
}

// Storage is the storage.
type Storage interface {
	// Lock locks the storage. Any subsequent attempt to call Lock will fail
	// until the last lock released.
	// After use the caller should call the Release method.
	Lock() (l util.Releaser, err error)

	// Log logs a string. This is used for logging. An implementation
	// may write to a file, stdout or simply do nothing.
	Log(str string)

	// GetFile returns a file for the given number and type. GetFile will never
	// returns nil, even if the underlying storage is closed.
	GetFile(num uint64, t FileType) File

	// GetFiles returns a slice of files that match the given file types.
	// The file types may be OR'ed together.
	GetFiles(t FileType) ([]File, error)

	// GetManifest returns a manifest file. Returns os.ErrNotExist if manifest
	// file does not exist.
	GetManifest() (File, error)

	// SetManifest sets the given file as manifest file. The given file should
	// be a manifest file type or error will be returned.
	SetManifest(f File) error

	// Close closes the storage. It is valid to call Close multiple times.
	// Other methods should not be called after the storage has been closed.
	Close() error
}
