package kopsa

import (
    "os"
    "io"
)

const mega = 1000000
var bufferSize = 10 * mega

func copyReaderWriter(dst io.Writer, src io.Reader) (int64, error) {
    buffer := make([]byte, bufferSize)
    return io.CopyBuffer(dst, src, buffer)
}

// If err == nil, file returned is OPEN
func appendFiles(dst string, srcs []string) (int64, error) {
    var (
        totalBytes  int64
        err         error
        destination *os.File
    )

    destination, err = os.Create(dst)
    if err != nil {
        return totalBytes, err
    }
    defer destination.Close()

    for _, src := range srcs {
        f, err := os.Open(src)
        if err != nil {
            return totalBytes, err
        }

        n, err := copyReaderWriter(destination, f)
        totalBytes = totalBytes + n
        if err != nil {
            f.Close()
            return totalBytes, err
        }

        if err = f.Close(); err != nil {
            return totalBytes, err
        }
    }

    return totalBytes, nil
}

// Does not close the source file
func copyFile(dst string, source *os.File) (int64, error) {
    var (
        totalBytes  int64
        destination *os.File
        err          error
    )

    destination, err = os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0660)
    if err != nil {
        destination.Close()
        os.Remove(dst)
        return totalBytes, err
    }

    totalBytes, err = copyReaderWriter(destination, source)
    if err != nil {
        destination.Close()
        os.Remove(dst)
        return totalBytes, err
    }

    if destination.Sync() != nil {
        destination.Close()
        os.Remove(dst)
        return totalBytes, err
    }

    err = destination.Close()
    if err != nil {
        os.Remove(dst)
        return totalBytes, err
    }

    return totalBytes, nil
}

func SetBufferSize(size int) {
    bufferSize = size
}

func Copy(dst string, srcs ...string) (int64, error) {
    var (
        err         error
        totalBytes  int64
        source      *os.File // This file must be closed inside this function
    )

    defer source.Close()

    if len(srcs) > 1 {
        totalBytes, err = appendFiles(dst, srcs)
        if err != nil {
            return totalBytes, err
        }
    } else {
        source, err = os.Open(srcs[0])
        if err != nil {
            return totalBytes, err
        }
        defer source.Close()

        n, err := copyFile(dst, source)
        totalBytes = totalBytes + n
        if err != nil {
            return totalBytes, err
        }

        err = source.Close()
        if err != nil {
            return totalBytes, err
        }
    }

    return totalBytes, err
}
