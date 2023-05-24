# api-go-expert

## Health Check
    http://localhost:8000/ping
    
## Swagger
    - http://localhost:8000/docs/index.html

## Using
    - create a user
        {
            "name": "Roger",
            "email": "email@email",
            "password": "123456"
        }
        - POST : http://localhost:8000/users

    - get a token
        {
            "email": "email@email",
            "password": "123456"
        }
        - POST : http://localhost:8000/users/generate_token

## Docker Prod (test)
    - docker build -t mrpsousa/api-go-expert:latest -f Dockerfile.prod .
    - docker run --rm -p 8000:8000 mrpsousa/api-go-expert:latest

