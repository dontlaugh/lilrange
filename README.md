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
to a the `Within` method. 

```go
// get a lilrange.Range
range, _ := lilrange.Parse("0130-0230")
// get the current time with Go's time package
now := time.Now()
// Are we within the Range?
if range.Within(now) {
    fmt.Println("We are inside the range")
}
```





