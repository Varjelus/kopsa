package kopsa

import (
    "path/filepath"
    "fmt"
    "io"
    "os"
)

const mega = 1000000

// Set default buffer size to ~10MB
var bufferSize = 10 * mega

// SetBufferSize sets the buffer size used in copying
func SetBufferSize(size int) {
    bufferSize = size
}

// Copy takes a destination path as the first argument and any number of 
// additional string arguments to be used as copy source paths. 
func Copy(dst string, srcs ...string) (int64, error) {
    // We want these to be declared at all times
    var (
        err          error
        totalBytes   int64
        medium       *os.File
    )

    // Get the absolute path of the destination file
    dst, err = filepath.Abs(dst)
    if err != nil {
        return totalBytes, err
    }

    // Create the intermediary file or truncate it if it 
    // already exists. It takes the name of the destination 
    // file + ".tmp".
    medium, err = os.Create(dst + ".tmp")
    if err != nil {
        return totalBytes, err
    }
    defer medium.Close()            // In case it's still open/existing
    defer os.Remove(dst + ".tmp")   // when this function returns.

    // Iterate over sources
    for _, src := range srcs {

        // Get the absolute path of the source file.
        src, err = filepath.Abs(src)
        if err != nil {
            return totalBytes, err
        }

        // Get file info of the source file.
        sfi, err := os.Stat(src)
        if err != nil {
            return totalBytes, err
        }

        // If the source file is not a regular file 
        // (it's a directory, device...), return error.
        if !(sfi.Mode().IsRegular()) {
            return totalBytes, fmt.Errorf("non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
        }

        // Open the source file
        f, err := os.Open(src)
        if err != nil {
            return totalBytes, err
        }

        // Create a buffer of bufferSize
        buffer := make([]byte, bufferSize)

        // Read the source file through the buffer 
        // to the intermediary file.
        n, err := io.CopyBuffer(medium, f, buffer)

        // Append the bytes to the total byte count
        totalBytes = totalBytes + n
        if err != nil {
            f.Close()
            return totalBytes, err
        }

        // Close the source file
        if err = f.Close(); err != nil {
            return totalBytes, err
        }
    }

    // Flush in-memory data to disk
    if medium.Sync() != nil {
        return totalBytes, err
    }

    // Close the intermediary file
    err = medium.Close()
    if err != nil {
        return totalBytes, err
    }

    // Rename (aka move) the intermediary file, 
    // in effect removing the ".tmp" extension 
    // so the file path satisfies the requested 
    // destination path.
    err = os.Rename(dst + ".tmp", dst)
    if err != nil {
        return totalBytes, err
    }

    // Return the total bytes copied and a nil error
    return totalBytes, nil
}
