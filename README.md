# Backend authentication

This project was made for a test assignment - part of the authentication service.

Issues and refreshes `SHA512` **Access** and **Refresh** JWT tokens.
These tokens are generated together with cohesion, so they cannot be used apart from each other.

## Build
Use the command 
```code:shell
$ go run ./main/
```
Copy this request localhost:8080/signin/?user-id=<your-user-id> 
and paste it into postman.
copy the body and send it to the body to check the refresh operation
localhost:8080/insert