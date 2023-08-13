# Url Shortener
A service for shortening URLs using golang and mongo db 

## API Reference

#### Shorten a Long URL
- This endpoint needs long url and user id in the request body, user id is taken as input to generate a unique short url for each user
- This implementation generates the same short url for a given combination of long url and userId

```bash
POST /shorten
```

Sample curl request
```bash
curl --location 'http://localhost:8080/shorten' \
--header 'Content-Type: application/json' \
--data '{
    "long_url" : "https://www.google.com",
    "user_id" : "asdasd"
}'
```
Request body
```json
{ 
    "long_url": "<put long url here>" ,
    "user_id" : "<put user id here>"
}
```

201 Response in case where short url gets created:
```json
{
    "shortUrl": "<short url>"
}
```

403 Response in case where short url already exists:
```json
{
    "Error": "Short url already exists: <short url>"
}
```

#### Browse a Short URL
When the user opens a short url in browser this endpoint redirects the request to the corresponding long url

```bash
GET /:shortUrlCode
```

| Parameter      | Type     | Description                       |
|:---------------| :------- | :-------------------------------- |
| `shortUrlCode` | `string` | **Required**. short url code of a given long url|


#### Delete a Short URL
- User can delete a short url 
- Upon deletion the short url access logs are also deleted

```bash
DELETE /:shortUrlCode
```

| Parameter   | Type     | Description                       |
|:------------| :------- | :-------------------------------- |
| `shortUrlCode` | `string` | **Required**. short url code of a given long url|

#### Access Details of a Short URL
- Access counts are calculated based on the logs stored in the database
- Timestamp is updated in the database everytime a request for redirect comes
- When a request to shorten an existing short url is made the access count for that short url is returned in the response
- Access count has 3 variables that represent 3 things:
    - All time count
    - Count in last 24 hours
    - Count in last 7 days

```bash
GET access-details/:shortUrlCode
```

| Parameter   | Type     | Description                       |
|:------------| :------- | :-------------------------------- |
| `shortUrlCode` | `string` | **Required**. short url code of a given long url|

Response
```json
{
    "Access details": "AllTimeAccessCount : <some value>, pastTwentyFourHoursAccessCount : <some value>, pastWeekAccessCount : <some value>"
}
```


## Installation
- Please make use of Safari browser or any other browser apart from Google Chrome since chrome caches some of the requests and does not invoke the application which leads to incorrect access counts
- Make sure you have mongo db & golang installed on your local machine
- One can make use of GoLand IDE provided by intellij in order to run this project
- MAC users can follow the below installation steps, in case of other OS one can follow similar steps 

Mongo installation:
```
brew tap mongodb/brew
brew install mongodb-community
brew services start mongodb-community
```

Golang installation:
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

## Running Tests Locally
Test are located under respective packages
