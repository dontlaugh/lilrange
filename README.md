# lilrange

An opinionated string DSL for creating UTC time ranges.

## Usage

Get a `lilrange.Range` from a string.

```go
range, err := lilrange.Parse("0000-0400")
```

The returned range's `End` time is guaranteed to be in the future. The `Start`
time could be in the past, or the future.

To determine if your current system time is inside the range, pass a `time.Time`
to the `Within` method. 

```go
// get a lilrange.Range
r, _ := lilrange.Parse("0130-0230")

// get the current time with Go's time package
now := time.Now()

// Are we within the Range?
if r.Within(now) {
    fmt.Println("We are inside the range")
}
```

## Test Program

A test program can be compiled from the *lilrange* sub directory, or use `go get -u`
if `$GOPATH/bin` is on your PATH.

```
go get -u github.com/dontlaugh/lilrange/lilrange
```

Give it a lilrange string to test the behavior.

```
lilrange 0415-0445
```

