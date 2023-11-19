# how to handle error in a service
## error for user
we should mapping each error to a specific error code, this code should be user friendly and meaningful for user, the translate can be happened at user interface, e.g at controller
## error for developer
we just chain each error from buttom to top, and log this error at top level, for example
```go
func outerFunc() error {
    if err := innerFunc(); err != nil {
        return fmt.Errorf("failed to handle inner: %w", err)
    }
    // other logic
}
```

this practice was inspired by [this post](https://www.sobyte.net/post/2023-05/go-error/)

# handle errors among micro-services
## error trace in one api call
most tracing system like zipkin or skywalking can be used to achive a calling chain trace

## error trace in one logic cal
need to aggregate some of the trace logs into one logicid, inspired by [this blog](https://tech.meituan.com/2022/07/21/visualized-log-tracing.html)