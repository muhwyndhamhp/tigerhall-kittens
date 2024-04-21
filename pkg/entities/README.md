# Authentication Flow
## User Authentication
We're using JWT for the authentication flow. The flow is as follows:
1. User will send a request to the server with their credentials.
2. The server will validate the credentials, and if it's correct, the server will generate a JWT token.
3. The server will send the JWT token back to the user.
4. The user will store the JWT token in their local storage.
5. The user will send the JWT token in the header (using `Authentication` header) of the request to the server.
6. The server will validate the JWT token, and if it's valid, the server will process the request.

## How the JWT Token is Generated
The JWT token is generated using the following steps:
1. The server will generate a JWT token using the following payload:
```json
{
  "id": 123,
  "name", "user_name",
  "email": "user_email",
  "exp": 123, // Expiry Time
}
```
2. The server will sign the JWT token using the `JWT_SECRET` environment variable.

## How the JWT Token is Validated
The JWT token is validated using the following steps:
1. The server will validate the JWT token using the `JWT_SECRET` environment variable.
2. The server will check the expiry time of the JWT token. If the expiry time is in the past, parser will return an error.
3. The server will check the `id` with real user id in the database. If the user id is not found, parser will return an error.
4. If all checks are passed, the server will append User Entity to the Request Context.
5. Resolvers can access the User Entity from the Request Context, and validate the user's permission if needed.

## Password Hashing
We're using `bcrypt` for hashing the password. The flow is as follows:
1. User will send a request to the server with their credentials.
2. The server will hash the password using `bcrypt` and store it in the database.
3. The server will compare the hashed password with the password stored in the database.

## Token Invalidation and Refresh Token
As token do have expiration date (24 hours), we need to make sure that user have the best user experience possible without required to login every 24 hours. We can achieve this by using Refresh Token. As we're using basic auth via JWT, we don't really have refresh token itself. So the implementation is a bit simplistic, but no less secure. 

There are no risk if token being reused far after it's expiration date, as the JWT itself has expiration mechanism and won't parse if the token is expired. The only risk is if the token is being used after it's invalidated. But this is a risk that we can take, as the token is invalidated after it's refreshed.

The flow to handle those scenarios is as follows:
- When user login, we will generate a new JWT token with 24 hours expiration date.
- Ideally when user do another request, client should decide whether the token is almost expired or not. If it's almost expired, client should request a new token by sending the old token to the server.
- When token is refreshed, the old token will be recorded in the database as invalidated token. So if the old token is used, server will return an error.
- If the token is invalidated, user should login again to get a new token.

