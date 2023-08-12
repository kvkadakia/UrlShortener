
## API Reference

#### Shorten a Long URL

This endpoint needs long url and user id in the request body, user id is required in order to generate a unique short url for each user

Endpoint:
```http
  POST /shorten
```

Request body:
```json
{ 
    "long_url": "<put long url here>" ,
    "user_id" : "<put user id here>"
}
```

#### Browse a Short URL
When the user opens a short url in browser this endpoint redirects the request to the corresponding short url

```http
  GET /:code
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `code`      | `string` | **Required**. short url code of a given long url|



## Installation

- Make sure you have mongo db & golang installed on your local machine
- One can make use of GoLand IDE provided by intellij in order to run this project.
- MAC users can follow the below installation steps, in case of other OS one can follow similar steps
- Mongo installation:
```
brew tap mongodb/brew
brew install mongodb-community
brew services start mongodb-community
```

-  Golang installation:
 ```
 brew install golang
 ```