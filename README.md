# kopsa
Drop-in file copying module

# Usage
Get it: `go get github.com/Varjelus/kopsa`

Import it: `import github.com/Varjelus/kopsa`


## Use it

`bytesCopied, err := kopsa.Copy("destination.txt", "source.txt")`, handle error

or `... kopsa.Copy("destination.txt", "source1.txt", "source2.txt", "source3.txt") ...`
