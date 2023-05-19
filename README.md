# api-go-expert

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

## Docker actions
    - docker build -t mrpsousa/api-go-expert:latest -f Dockerfile.prod .