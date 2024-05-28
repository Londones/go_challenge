# How to use the logger

## General informations
- Our logger use loggrus to write logs
- Our logger use fmt and os to write logs in a file named logs.log
- logs.log is located at "./"

1. Import the package:
```go
 import (
    "go-challenge/internal/utils"
)
```

2. Call the function: 
```go
    utils.Logger(logLevel, field1, field2, msg string)
```

3. Call the function with a string concatenation:
```go
    utils.Logger("debug", "TEST", "test", fmt.Sprintf("your message %v", v))
```

## Log levels
- DebugLevel => "debug"
- InfoLevel => "info"
- WarnLevel => "warn"
- ErrorLevel => "error"
- FatalLevel => "fatal"
- PanicLevel => "panic"

