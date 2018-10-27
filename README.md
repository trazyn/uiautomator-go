
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
client := ug.New("10.10.0.101:7912")
info := client.GetDeviceInfo()
fmt.Printf(info)

```

[https://github.com/openatx/uiautomator2#basic-api-usages](https://github.com/openatx/uiautomator2#basic-api-usages)
