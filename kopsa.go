package kopsa

import (
    "fmt"
    "os"
    "path/filepath"
    "io"
)

const mega = 1000000
var bufferSize = 10 * mega

func SetBufferSize(size int) {
    bufferSize = size
}

func Copy(dst string, srcs ...string) (int64, error) {
    var (
        err         error
        totalBytes  int64
        source      *os.File // This file must be closed inside this function
        destination *os.File
    )

    defer source.Close()

    dst, err = filepath.Abs(dst)
    if err != nil {
        return totalBytes, err
    }

    destination, err = os.Create(dst + ".tmp")
    if err != nil {
        return totalBytes, err
    }
    defer destination.Close()

    for _, src := range srcs {
        src, err = filepath.Abs(src)
        if err != nil {
            return totalBytes, err
        }

        sfi, err := os.Stat(src)
        if err != nil {
            return totalBytes, err
        }

        if !(sfi.Mode().IsRegular()) {
            return totalBytes, fmt.Errorf("non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
        }

        f, err := os.Open(src)
        if err != nil {
            return totalBytes, err
        }

        buffer := make([]byte, bufferSize)
        n, err := io.CopyBuffer(destination, f, buffer)
        totalBytes = totalBytes + n
        if err != nil {
            f.Close()
            return totalBytes, err
        }

        if err = f.Close(); err != nil {
            return totalBytes, err
        }
    }

    if destination.Sync() != nil {
        destination.Close()
        os.Remove(dst + ".tmp")
        return totalBytes, err
    }

    err = destination.Close()
    if err != nil {
        os.Remove(dst + ".tmp")
        return totalBytes, err
    }

    err = os.Rename(dst + ".tmp", dst)
    if err != nil {
        os.Remove(dst + ".tmp")
        return totalBytes, err
    }

    return totalBytes, err
}
