# cmutex

`cmutex` provides a mutex which accepts a native `context` item as a means of guaranteeing bounded wait.

## Example usage

```golang
var m cmutex.Mutex
ctx, cf := context.WithTimeout(context.TODO(), 1 * time.Second)

if m.Lock(ctx) != nil {
  return err
}

defer m.Unlock()

// computations, etc..
```

N.B. This package uses a `channel` under the hood and is therefore a little slower than the native `sync.Mutex` provided in the `stdlib`.