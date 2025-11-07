# To install the postgres
`sudo apt install postgresql postgresql-contrib`


# For .env file check .env.example


# To run the project
`go run server.go`


# Implemented below mentioned routes for the project

To read the key value
Method: GET
URL: `localhost:8080/read?key=16`

To create the key value
Method: POST
URL: `localhost:8080/create`
Body:
    {
        "key": 18,
        "value": "Val 18"
    }

To update the key value
Method: PUT
URL `localhost:8080/update`
Body:
    {
        "key": 3,
        "value": "Val 3 Updated"
    }

To delete the key
Method: DELETE
Body:
    {
        "key": 9
    }