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

func TestSingleSourceFileWriterCopy(t *testing.T) {
    stat, err := os.Stat(this)
    if err != nil {
        t.Errorf("os.Stat failed: %s", err.Error())
        return
    }

    destOpen, err := os.Create(dest)
    if err != nil {
        t.Errorf("os.Open failed: %s", err.Error())
        return
    }
    defer destOpen.Close()

    n, err := Copy(destOpen, this)
    defer os.Remove(dest)
    if err != nil {
        t.Errorf("Copy failed: %s", err.Error())
        return
    }

    if n != stat.Size() {
        t.Errorf("Copy failed: unmatching file sizes")
        return
    }
}

func TestSingleSourceFileStringCopy(t *testing.T) {
    stat, err := os.Stat(this)
    if err != nil {
        t.Errorf("os.Stat failed: %s", err.Error())
        return
    }

    n, err := Copy(dest, this)
    defer os.Remove(dest)
    if err != nil {
        t.Errorf("Copy failed: %s", err.Error())
        return
    }

    if n != stat.Size() {
        t.Errorf("Copy failed: unmatching file sizes")
        return
    }
}

func TestMultiSourceStringCopy(t *testing.T) {
    stat, err := os.Stat(this)
    if err != nil {
        t.Errorf("os.Stat failed: %s", err.Error())
        return
    }

    file, err = os.Create(src)
    if err != nil {
        os.Remove(src)
        t.Errorf("Test file creation failed: %s", err.Error())
        return
    }

    written, err := file.Write(data)
    if err != nil {
        file.Close()
        os.Remove(src)
        t.Errorf("Test file write failed: %s", err.Error())
        return
    }

    if err := file.Close(); err != nil {
        os.Remove(src)
        t.Errorf("Test file close failed: %s", err.Error())
        return
    }

    n, err := Copy(dest, this, src)
    if err != nil {
        os.Remove(src)
        t.Errorf("Copy failed: %s", err.Error())
        return
    }

    if n != int64(written) + stat.Size() {
        os.Remove(src)
        os.Remove(dest)
        t.Errorf("Copy failed: unmatching file sizes")
        return
    }

    if err := os.Remove(dest); err != nil {
        os.Remove(src)
        t.Errorf("Destination file remove failed: %s", err.Error())
        return
    }

    if err := os.Remove(src); err != nil {
        t.Errorf("Source file remove failed: %s", err.Error())
        return
    }
}
