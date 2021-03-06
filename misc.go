// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package curio

import (
	"io"
	"os"
)

type (
	stater interface {
		Stat() (os.FileInfo, error)
	}
	lener interface {
		Len() int
	}
	sizer interface {
		Size() int64
	}
)

// Len returns the length of a Seeker.
// If s has a Len, Size, or Stat method, one of those will be used. Otherwise,
// Seek will be used to determine the length, before restoring the cursor to its
// previous position.
func Len(s io.Seeker) (int64, error) {
	switch s := s.(type) {
	case sizer:
		return s.Size(), nil
	case lener:
		return int64(s.Len()), nil
	case stater:
		info, err := s.Stat()
		if err != nil {
			return 0, err
		}
		return info.Size(), nil
	}
	i, err := s.Seek(0, os.SEEK_CUR)
	if err != nil {
		return 0, err
	}
	j, err := s.Seek(0, os.SEEK_END)
	if err != nil {
		return j, err
	}
	_, err = s.Seek(i, os.SEEK_SET)
	return j, err
}
