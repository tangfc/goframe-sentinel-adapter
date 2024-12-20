### What This PR Does / Why We Need It
This PR adds a goframe framework adapter for Sentinel. With this adapter, developers can easily integrate Sentinel's rate limiting and circuit breaker functionalities into their goframe projects, thereby enhancing system stability and reliability.

### How It Was Done
1. Implemented the `SentinelMiddleware` function to integrate Sentinel into the goframe framework.
2. Added the `withResourceExtractor` option to allow developers to customize resource extraction logic.
3. Added the `WithBlockFallback` option to allow developers to customize fallback logic when blocked.
4. Wrote corresponding test cases to ensure the correctness and stability of the adapter.

### How to Verify It
1. go test -run ^TestSentinelMiddlewareDefault -v 
2. go test -run ^TestSentineIMiddlewareExtractor -v 
3. go test -run ^TestSentinelMiddlewareFallback -v
