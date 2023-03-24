## Elotus API

import to Postman: 

```json
{
	"info": {
		"_postman_id": "c609784b-8296-44d1-bf22-e266dbdba023",
		"name": "elotus",
		"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
	},
	"item": [
		{
			"name": "Ping",
			"request": {
				"method": "GET",
				"header": []
			},
			"response": []
		},
		{
			"name": "Login",
			"request": {
				"method": "POST",
				"header": [],
				"body": {
					"mode": "raw",
					"raw": "{\r\n    \"username\": \"abc\",\r\n    \"password\": \"123\"\r\n}",
					"options": {
						"raw": {
							"language": "json"
						}
					}
				},
				"url": {
					"raw": "localhost:8080/login",
					"host": [
						"localhost"
					],
					"port": "8080",
					"path": [
						"login"
					]
				}
			},
			"response": []
		},
		{
			"name": "Register",
			"request": {
				"method": "POST",
				"header": []
			},
			"response": []
		},
		{
			"name": "Upload",
			"request": {
				"auth": {
					"type": "bearer"
				},
				"method": "POST",
				"header": [],
				"body": {
					"mode": "formdata",
					"formdata": []
				},
				"url": {
					"raw": "localhost:8080/upload",
					"host": [
						"localhost"
					],
					"port": "8080",
					"path": [
						"upload"
					]
				}
			},
			"response": []
		}
	]
}

```

JWT will exprive in 5 minutes

Register -> Login
Use access token from Login, attach to Bear Authoration in Upload request

I will update unit tests later. My laptop has some issues so I didn't have time. Sorry :(((

I prepared Dockerfile.

Run 
```
docker-compose up --build
```