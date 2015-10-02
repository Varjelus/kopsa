package kopsa

import (
    "os"
    "testing"
)

var (
    this    = os.Args[0] // Use this code as a source
    src     = "test-source"
    data    = []byte("some other test data")
    dest    = "destination"
    file    *os.File
    err     error
)

func TestSingleSourceCopy(t *testing.T) {
    stat, err := os.Stat(this)
    if err != nil {
        t.Errorf("os.Stat failed: %s", err.Error())
    }

    n, err := Copy(dest, this)
    defer os.Remove(dest)
    if err != nil {
        t.Errorf("Copy failed: %s", err.Error())
    }

    if n != stat.Size() {
        t.Errorf("Copy failed: unmatching file sizes")
    }
}

func TestMultiSourceCopy(t *testing.T) {
    stat, err := os.Stat(this)
    if err != nil {
        t.Errorf("os.Stat failed: %s", err.Error())
    }

    file, err = os.Create(src)
    if err != nil {
        os.Remove(src)
        t.Errorf("Test file creation failed: %s", err.Error())
    }

    written, err := file.Write(data)
    if err != nil {
        file.Close()
        os.Remove(src)
        t.Errorf("Test file write failed: %s", err.Error())
    }

    if err := file.Close(); err != nil {
        os.Remove(src)
        t.Errorf("Test file close failed: %s", err.Error())
    }

    n, err := Copy(dest, this, src)
    if err != nil {
        os.Remove(src)
        t.Errorf("Copy failed: %s", err.Error())
    }

    if n != int64(written) + stat.Size() {
        os.Remove(src)
        os.Remove(dest)
        t.Errorf("Copy failed: unmatching file sizes")
    }

    if err := os.Remove(dest); err != nil {
        os.Remove(src)
        t.Errorf("Destination file remove failed: %s", err.Error())
    }

    if err := os.Remove(src); err != nil {
        t.Errorf("Source file remove failed: %s", err.Error())
    }
}
