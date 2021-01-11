# PersonalDiary

* Make sure you have MySQL and Redis with ReJSON installed on your local machine.
* You can download MySQL from [here](https://dev.mysql.com/downloads/installer/) and Redis from [here](https://redis.io/download).
* If you're using Windows, you need to have Docker installed on your local machine in order to run Redis With ReJSON. Paste the following on the command line to install Redis with ReJSON.
```
docker run -p 6379:6379 --name redis-redisjson redislabs/rejson:latest
```

### Installing
* To clone this repository, you need have to install [GIT](https://git-scm.com) on your local machine.
* Paste the following on the command line:
```
$ git clone https://github.com/evrintobing17/PersonalDiary
```

### Deployment
* After you have cloned the repository, navigate to ***.env*** file and enter your MySQL database configurations.
```
# Mysql
API_SECRET=evrin17 # For JSON Web tokens, can be anything
DB_HOST=127.0.0.1
DB_DRIVER=mysql 
DB_USER=DB_USER_NAME
DB_PASSWORD=DB_USER_PASSWORD
DB_NAME=YOUR_DB
DB_PORT=3306

# Mysql Test
API_SECRET_TEST=evrin17 # For JSON Web tokens, can be anything
DB_HOST_TEST=127.0.0.1
DB_DRIVER_TEST=mysql 
DB_USER_TEST=DB_USER_NAME
DB_PASSWORD_TEST=DB_USER_PASSWORD
DB_NAME_TEST=YOUR_DB_TEST
DB_PORT_TEST=3306
```
* Run the ***diary_db.sql*** file in your database workbench.
* Navigate to the folder where ***main.go*** resides and enter the following command to run the program:
```
go run main.go
```
## Running the tests
* Each repository folder has a test file. We will create a separate database for testing purposes. 
* Run the ***diary_db_test.sql*** file in your database workbench.
* In order to run a particular test, run the following command:
```
go test --run TestName
```