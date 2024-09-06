This is a ***fullstack***ish application that allows users to create, read, update, and delete todos. Frontend is built with HTMX, Bootstrap, Go Template, CSS and JS; it is not good, but it works. Backend is built with Go.

### Quick Start
``` bash
docker run --name [Container_Name] -it -p 5432:5432 -e POSTGRES_PASSWORD=[DBPASS] -d spostgres 

make
```

Container name is up to you. 

DBPASS is the password you want to use for the database. You should also set the DBPASS variable in the .env file.

SECRET variable in the .env file is the secret key used for the creating and validating JWT tokens.