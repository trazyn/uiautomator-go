
This is a simple golang wrapper for working uiautomator.

## Setup


To install this library, simple:



```bash
go get github/trazyn/uiautomator-go
```


Import the package:



```go
import ug "github/trazyn/uiautomator-go"
```


## Quick start

First, let yours mobile and PC join the same network.

```go
ua := ug.New(&ug.Config{
    Host: "10.10.20.78",
    Port: 7912,
})

ua.Unlock()

// Show toast
toast := ua.NewToast()
toast.Show("hallo world", 10)
```

[https://github.com/openatx/uiautomator2#basic-api-usages](https://github.com/openatx/uiautomator2#basic-api-usages)
