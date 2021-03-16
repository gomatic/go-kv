# kv

[![CI](https://github.com/gomatic/kv/actions/workflows/ci.yml/badge.svg)](https://github.com/gomatic/kv/actions/workflows/ci.yml)

Package `kv` provides small, dependency-light helpers for reading and manipulating key/value environment data: an `Environment` map, lookups with fallbacks, loading from JSON/YAML, and scoped set/restore of the process environment.

## Install

```bash
go get github.com/gomatic/kv
```

Requires Go 1.26+.

## Usage

### Lookups with fallback

```go
port := kv.GetOr("PORT", "8080")
url := kv.FirstOr("http://localhost", "PRIMARY_URL", "FALLBACK_URL")
value, ok := kv.Lookup("OPTIONAL")
```

The package-level helpers read the process environment; the same methods exist on an `Environment` map for standalone use:

```go
env := kv.New()                 // snapshot of os.Environ()
env = kv.Parse(os.Environ())    // or parse any "KEY=VALUE" slice
got := env.GetOr("PORT", "8080")
```

### Typed values

Environment values are strings; `LookupAs` / `GetOrAs` convert them with any `func(string) (T, error)`. Standard-library parsers already have that signature, so they pass directly — the type is inferred:

```go
port := kv.GetOrAs("PORT", strconv.Atoi, 8080)            // int
debug := kv.GetOrAs("DEBUG", strconv.ParseBool, false)    // bool
timeout := kv.GetOrAs("TIMEOUT", time.ParseDuration, time.Second)

// Distinguish "unset" from "set but invalid":
value, ok, err := kv.LookupAs("PORT", strconv.Atoi)
```

### Load from JSON or YAML

```go
env, err := kv.LoadFromJSONFile("config.json")
env, err = kv.LoadFromYAMLFile("config.yaml")
```

### Load a value from a file

```go
// If SECRET is unset and SECRET_FILE names a readable file, load it into SECRET.
err := kv.SetFromFile("SECRET", "SECRET_FILE")
```

### Scoped set and restore

```go
// Apply an environment for the duration of the calls, then restore it.
err := kv.WrapCalls(kv.Environment{"DEBUG": "1"}, func() error { return run() })

// Or set a single variable and restore it later.
restore := kv.SetWithRestore("DEBUG", "1", false)
defer restore()
```

## API

| Symbol | Description |
|--------|-------------|
| `Get` / `Lookup` / `GetTrimmed` | Read a process environment variable |
| `GetOr` / `First` / `FirstOr` / `LookupFirst` | Read with fallbacks across one or more keys |
| `LookupAs[T]` / `GetOrAs[T]` | Read a variable converted to `T` via a `func(string) (T, error)` |
| `SetWithRestore` | Set a variable and get a function that restores its prior state |
| `SetFromFile` | Populate a variable from a file named by another variable |
| `Environment` | A `map[string]string` with the same lookup/load methods |
| `New` / `Parse` | Build an `Environment` from the process environment or a `KEY=VALUE` slice |
| `LoadFromJSON*` / `LoadFromYAML*` | Load an `Environment` from JSON/YAML readers or files |
| `Wrapper` / `WrapCalls` | Run functions between environment set and restore |
| `Name` | A typed environment-variable name with `Value` / `ValueOr` / `Lookup` |
| `Names` | A name→value map with process-environment fallback and `${VAR}` expansion |

## License

See [LICENSE](LICENSE).
