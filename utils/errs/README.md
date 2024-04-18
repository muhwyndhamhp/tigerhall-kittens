# Error Handling
Each uniquely identifiable error will have to use the `ServiceError` struct, which is a struct that implements the `error` interface. 

Later down the line when returning from the resolvers, it also should return by using the `errs.RespError()` function
instead. 

This is to ensure that error with unique code can be identified and handled properly by the client. `ServiceError` struct will have `ErrorCode` field, which is a unique code for the error. And if returned using `errs.RespError()`, it will also fill the `errors.extensions.code` field in the response.

Example Response of not using `ServiceError`:
```json
{
  "errors": [
    {
      "message": "ErrUserByCtxNotFound: user not found in context",
      "path": [
        "createSighting"
      ],
    }
  ],
  "data": null
}
```

Example Response of using `ServiceError`:
```json
{
  "errors": [
    {
      "message": "user not found in context",
      "path": [
        "createSighting"
      ],
      "extensions": {
        "code": "ErrUserByCtxNotFound" // Unique Error Code
      }
    }
  ],
  "data": null
}
```
