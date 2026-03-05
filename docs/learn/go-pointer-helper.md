# Generic Pointer Helper (TLDR: Value -> \*Value)

**Input**: Any value (e.g., `"active"`, `123`, `true`)
**Output**: A pointer to that value (e.g., `*string`, `*int`, `*bool`)

---

## What is it?

The `Ptr[T any](v T) *T` function is a generic Go helper that takes a value of any type and returns a pointer to it. It uses Go's escape analysis to safely allocate the value on the heap and return its address.

## Why is it important?

In Go, you cannot directly take the address of a literal value.

- Invalid: `&"string"`
- Invalid: `&10`
- Invalid: `&true`

This becomes painful when working with structs that use pointers for optional fields (like in database models, configuration structs, or API clients). Without this helper, you must declare a temporary variable for every single literal you want to pass as a pointer.

## Examples

### The Problem (Without Helper)

```go
type User struct {
    Name *string
    Age  *int
}

func main() {
    // verbose and clunky
    name := "Alice"
    age := 30
    u := User{
        Name: &name,
        Age:  &age,
    }
}
```

### The Solution (With Helper)

```go
import "your/project/utils"

type User struct {
    Name *string
    Age  *int
}

func main() {
    // clean and inline
    u := User{
        Name: utils.Ptr("Alice"),
        Age:  utils.Ptr(30),
    }
}
```
