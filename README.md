# Url Shortener
An internal service for shortening URLs using golang and mongo db that anybody can use

## API Reference

#### Shorten a Long URL

This endpoint needs long url and user id in the request body, user id is required in order to generate a unique short url for each user

```bash
POST /shorten
```

```json
{ 
    "long_url": "<put long url here>" ,
    "user_id" : "<put user id here>"
}
```

Response in case where short url does not exist:
```json
{
    "shortUrl": "<short url>"
}
```

Response in case where short url exists:
```json
{
    "info": "Short url already exists: <short url> | totalAccessCount : <some value>, pastTwentyFourHoursAccessCount : <some value>, pastWeekAccessCount : <some value>"
}
```

#### Browse a Short URL
When the user opens a short url in browser this endpoint redirects the request to the corresponding short url

```bash
GET /:code
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `code`      | `string` | **Required**. short url code of a given long url|





## Installation
- Please make use of Safari browser or any other browser apart from Google Chrome since chrome caches some of the requests and does not invoke the application which leads to incorrect access counts
- Make sure you have mongo db & golang installed on your local machine
- One can make use of GoLand IDE provided by intellij in order to run this project
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


## Run Locally

Clone the project

```bash
  git clone https://github.com/kvkadakia/UrlShortener.git
```

Go to the project directory

```bash
  cd UrlShortener
```


Start the server

```bash
  go run main.go
```

