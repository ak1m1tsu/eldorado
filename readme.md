# Eldorado

Eldorado is a simple todo application implemented using microservices.

## Technologies

- RabbitMQ (communication between auth service/statistics service and email sending service)
- gRPC (communication between auth service and todo service)
- Redis (cache, session storage)
- PostgreSQL (database)
- Cron (schedule statistics service)
- Docker (deploy services, message broker, etc.)
- SMTP (send emails)

## Design Diagram

![](./assets/Eldorado%20Design.svg)

## Services

### Auth Service

The authorization service runs on the gRPC framework, its main functionality is to create new users, issue JWTs to access the functionality of the main service, and refresh the access token.

All tokens are stored in Redis. When a new user is created, the auth service communicates with the email sending service to send a welcome message to the new user's email.

### Email sending Service

The email sending service is created useing RabbitMQ. It's main functionality is to listen to the queue in which email messages arrive and send them to users using SMTP.

### Statistics Service

The statistics service collect statistic about users uncompleted tasks and send it to email sending service. It collects data on a schedule, with cron task syntax, which can be edited in the config file.

### Todo Service

ToDo service provides CRUD operations over tasks for authorized users, possibility of new user registration and authorization for anonymous users. For authorization and registration data the service interacts with authorization service using gRPC. All data is cached in Redis, and the cache is cleared when the data is changed.

## Hot to run?

Create `.env` file:

```shell
# API
API_CONFIG_PATH=/etc/api.local.yaml
API_PORT=8080
# AUTH SERVICE
AUTH_SERVICE_CONFIG_PATH=/etc/auth.local.yaml
AUTH_SERVICE_PORT=8081
RSA_PRIVATE_KEY="LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb2dJQkFBS0NBUUVBaUYyMGlTUVFtRlVnZGZlM3RsQmVVcXFwRDRUTDdhRFAwS0g1ZzJob1o3MEFjcDhVCm9rdTZZN2MrNklBVk91alJ6TEJjQlpKeHBpNGl3bkhGaFd3c1ZEV2ZBbU55TVZtYnc4L2JKa0p3dnJXVVNOaDEKYXJVdUpaK3JPbWhNU1l5UXVxZ3hNOUx5OVVvelFuU2MwU1V1Zkl6Z2k2RitvQ0FyOGZiRFVscGdXUDJFM0xVdwp5WFhnWWUrdWE5MW1UUzdCMTdiRlZWN2N6Q3ZuMmZxbDkzQjFCeTFsN1htVEtwemxHcVNaL1JYMlhtWGtnN1BkCmw3ZFljdG5mbWN0TW9sbVQ1RDF5djEwOFZOVVFsUGNuMUxTTUJMS0w5RVNTdTBIU0hYVWlKRU1CWDN2bkFoS1IKQjNLTHI3cVl3c3crTXp3MTdGeUgyYmJhcHZwczJMRUpoREFBQ1FJREFRQUJBb0lCQUF4dngxeE9qcmpsNHB6LwpwNkNYK2RJK1FFYnJESkl5ZldHQXREblkxdFRITnZnOUsrdVZUbjF6bytnZWJsRStGSXcvZFZVSXd4YXQxSHU1ClZwTlJoMFZ5MG5xc1NTalpERXl5YzBFdEJBMVFrQ0tJbzBURkcxMVJENU8zR1dZSHpOZEpLWnVaWEpFa1lFSnEKVWpiODFoMkQxNkxFYXNEOXppUTJKaUFubit2ZHFKa3RidXp2a2NpRDhJbG9PRkJkQW9rS0o3ak05enVVY2VjMgpvTVplQ1V4UGFBcWtSSUhvckR2K3Bic09TRFdLK3A5NmlkMzc1dGx3aVRFdE9ZTDMxRTMxKzJnWmhWUVJpT1FJCkJzQXV0SXFrT0lpRWRNVDlDSENBKzE5V3EvL0lMb2gzQStjU25ObWNGdVZLVzQ1aGd3KzQrbHZIUUV0K21NY0oKeG42ajlmVUNnWUVBM3M5UXNCNTJNT0gyMHdZTW5yWmtJdDBvYklaaFJwcEVYTTdmZWU2Z3kxSy9PcXJFUXY1eQpqQUtOMHRiWmhjZHcyZk9DUjNhWVN3TnZ6TzJmQ21iZ1RMVXhKcS94UUZ1bzcybGM3Nms5LzFjTTQzSGZadS9vCmxwZVMrOGhWUm9sU2lmYTFmc1c5S29obkZOQXk5Mktqay9ZaTdTU3d6QWFMUU9vSi92M29waGNDZ1lFQW5LM3IKMS92NGdJTEVtZVQ5b3EvS2REcytiMHZ4Q1N3aTE5cU9XRHh1TWJpRERpMG1yaGZRSTk2ZGN5dTN5QVBUUFprWQplUVpjcG4rQWt5MXJKZFBFQVplMjRHdzlSOEk1N0QvK1V3OEEvWnZCQlFrRFZVZEVSUmdYQWljK3FtbTB3NjZuCmMraEJjSkF0d29IYmFmMTY0NXlKR0V3d3dhQ3lKYnpHRWhaNGZ0OENnWUJNMjdqYXQvZFBUM05FUWU3eWhMb3cKS0pmOEw3SzRseGlENXp6Qm8rWkZuT3FvUXlYbWNqMXpQS0pObTUxM2YxL1hYeCtPcFVOTmhRYjQ2dk1VdEg0bApGOCszcDdPUjNzeDhvWnpVYzA0V21hR3hoNk5ucjlSTXRrYVdvZm1BbG5ncGJUZ1lYZit0LzFXSG9YWWpUaHhkCk91ZTdaQSswb1lGQnlEbmdneGZReXdLQmdFQWUrWWRDQ3BobGJCcGpXZDNydlpwRjZLNVowUUprK3JtR0szMDIKOWc0SktqRnlEd051b1hNY2x5bGNPZkYwaDA0TlNyTzFBOVBzR0cyalI5ZUtUQXd1ejl3VzBCbC9CbitHVFFvbQorTXZSbzNQeEZWa2dPbk5nZ1lJVEY1VmNmMnNhSGxQVU9IdmR0YXlGd29zay90Y0o2QjEyaVBtbFQyTWNWNTl0CmJORFhBb0dBZitWQzFZOGpoazVaY1NsNFVibXQrWWVPRjh6VXh3eGN5bDAvR1E4N2hKME5oMkN3RFJ6UFc2b0gKdDRjbGVNZU02RHhTanhkMHlHOE5kalVBd3NqQTFWK2ZLRU91cmlkNVNIdzBzQ1RVUXYzUDN0OWNtSDhqaXZRdQpXTEF5SGZzcW9BdUJQNTIvVmZoNXZWSVlqNnIvN2Q3emVJbFVvQkpjY09hQTduaDZhUW89Ci0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0t"
RSA_PUBLIC_KEY="LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUFpRjIwaVNRUW1GVWdkZmUzdGxCZQpVcXFwRDRUTDdhRFAwS0g1ZzJob1o3MEFjcDhVb2t1Nlk3Yys2SUFWT3VqUnpMQmNCWkp4cGk0aXduSEZoV3dzClZEV2ZBbU55TVZtYnc4L2JKa0p3dnJXVVNOaDFhclV1Slorck9taE1TWXlRdXFneE05THk5VW96UW5TYzBTVXUKZkl6Z2k2RitvQ0FyOGZiRFVscGdXUDJFM0xVd3lYWGdZZSt1YTkxbVRTN0IxN2JGVlY3Y3pDdm4yZnFsOTNCMQpCeTFsN1htVEtwemxHcVNaL1JYMlhtWGtnN1BkbDdkWWN0bmZtY3RNb2xtVDVEMXl2MTA4Vk5VUWxQY24xTFNNCkJMS0w5RVNTdTBIU0hYVWlKRU1CWDN2bkFoS1JCM0tMcjdxWXdzdytNencxN0Z5SDJiYmFwdnBzMkxFSmhEQUEKQ1FJREFRQUIKLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0t"
# EMAIL SERVICE
EMAIL_SENDER_CONFIG_PATH=/etc/esender.local.yaml
EMAIL_SENDER_PORT=8082
# STATISTIC SERVICE
STATISTIC_CONFIG_PATH=/etc/statistic.local.yaml
STATISTIC_SERVICE_PORT=8083
# POSTGRES
PG_USER=postgres
PG_DB=eldorado
PG_PASS=postgres
# REDIS
REDIS_PASS=guest
# SMTP
SMTP_EMAIL=your_smtp_email
SMTP_PASS=your_smtp_password
```

Next, run `make up` command

# ToDo API endpoints

**Error response**

```json
{
    "error": "message",
}
```

## Check health status

```shell
curl http://localhost:8080/health
```

**Response**

```json
{
    "message": "alive",
}
```

## Register new user

```shell
curl -X POST --data '{"email":"admin@example.com","username":"admin","password":"test1234"}' http://localhost:8080/api/auth
```

**Response**

```json
{
    "message": "ok",
}
```

## Get access and refresh tokens

```shell
curl -X POST --data '{"email":"admin@example.com","password":"test1234"}' http://localhost:8080/api/auth/token
```

**Response**
```json
{
    "access_token": "string",
    "refresh_token": "string",
}
```

## Refresh access token

```shell
curl -X POST --cookie "refresh_token=string" http://localhost:8080/api/refresh
```

**Response**

```json
{
    "access_token": "string",
}
```


## Tasks CRUD

Each of these requests requires the header authorization bearer with access token

```shell
curl -H "Authorization: Bearer <token>"
```

### Get tasks

```shell
curl http://localhost:8080/tasks
```

**Response**

```json
{
  "tasks": [
    {
      "id": "a4501171-30f5-4fd3-88a2-3d4089fb7c63",
      "title": "first task",
      "description": "this is my first task, haha!",
      "created_at": "2023-09-25T11:40:35Z",
      "is_completed": false
    }
  ]
}
```

### Create new task

```shell
curl -X POST --data '{"title":"hello","description":"go to home"}' http://localhost:8080/api/tasks
```

**Response**

```json
{
  "task": {
    "id": "8673ce18-6bcc-4c02-9c9a-997c3784f84b",
    "title": "hello",
    "description": "go to home",
    "is_completed": false,
    "created_on": "2023-10-01T04:44:58Z"
  }
}
```

### Update task

```shell
curl -X PUT --data '{"title":"go back", "description":"welcome", "is_completed":true}' http://localhost:8080/api/tasks/8673ce18-6bcc-4c02-9c9a-997c3784f84b
```

**Response**

```json
{
  "task": {
    "id": "8673ce18-6bcc-4c02-9c9a-997c3784f84b",
    "title": "go back",
    "description": "welcome",
    "created_at": "0001-01-01T00:00:00Z",
    "is_completed": true
  }
}
```

### Delete task

```shell
curl -X DELETE http://localhost:8080/api/tasks/8673ce18-6bcc-4c02-9c9a-997c3784f84b
```

**Response**

```json
{
    "message": "ok"
}
```