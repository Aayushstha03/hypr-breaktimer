package state

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

type Lock struct{ f *os.File }

func Acquire(lockPath string) (*Lock, bool, error) {
	if err := os.MkdirAll(filepath.Dir(lockPath), 0o755); err != nil {
		return nil, false, err
	}
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, false, err
	}
	if err := unix.Flock(int(f.Fd()), unix.LOCK_EX|unix.LOCK_NB); err != nil {
		_ = f.Close()
		return nil, false, nil
	}
	return &Lock{f: f}, true, nil
}

func (l *Lock) Release() error {
	if l == nil || l.f == nil {
		return nil
	}
	defer func() { l.f = nil }()
	if err := unix.Flock(int(l.f.Fd()), unix.LOCK_UN); err != nil {
		_ = l.f.Close()
		return err
	}
	return l.f.Close()
}
