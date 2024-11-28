# OneStepGPS Project called TrackPoint GPS

## Project Goals
```
The aim of the project was to showcase my skills with Vue and Go and demonstrate my ability to pick up new technologies fast. 

TrackPointGPS is an application that provides real-time GPS tracking functionality for devices. This application is designed to manage, track, and display device locations on an interactive map while offering user-friendly preference management features.
```

## Test User Login Credentials
```
Username: testuser
Password: password123
```

## Prerequisites
##### Before running the project, ensure you have the following installed:
```
Node.js (version 14.x or later)
npm (comes with Node.js)
Go (version 1.18 or later)
PostgreSQL (version 12 or later)
 - Ensure your PostgreSQL installation allows local connections and has a properly configured superuser account (e.g., `postgres`).
```

## Setting Up The Project
### Step 1: Create a new project and Clone the Repository
```
git clone https://github.com/Kevin-T-Cook/TrackPointGPS.git
cd TrackPointGPS
```

### Step 2: Install and Configure PostgreSQL

##### For Mac Users:
###### Install PostgreSQL
```
brew install postgresql
```

###### Start PostgreSQL
```
brew services start postgresql
```

##### For Windows Users:
###### Install PostgreSQL
```
sudo apt update
sudo apt install postgresql postgresql-contrib
```

###### Start PostgreSQL
```
sudo systemctl start postgresql
```

##### Additional Steps for Windows:
```
Download PostgreSQL from https://www.postgresql.org/download/.
Follow the installation wizard and set up a superuser during the installation process.
```

##### Log into PostgreSQL
```
psql -U <your_user_system_name> -d postgres
```

##### Here is an example
```
psql -U kevincook -d postgres
```

##### If you are unsure what your user system name is run this command to open PostgreSQL:
```
psql -U $(whoami) -d postgres
```
##### You can check your system user name by running the command:
```
\du
```

##### Configure your pg_hba.conf file to use trust authentication. Run the following command in the PostgreSQL shell to locate your pg_hba.conf:
```
SHOW hba_file;
```

##### Copy the path of the your pg_hba.conf file before running this command to exit the PostgreSQL (psql) shell:
```
\q
```

##### Run the following command to open the pg_hba.conf file:
```
code <pg_hba_file_path>
```

##### Make sure your pg_hba.conf looks like the following:
```
# TYPE  DATABASE        USER            ADDRESS                 METHOD

# "local" is for Unix domain socket connections only
local   all             all                                     trust
# IPv4 local connections:
host    all             all             127.0.0.1/32            trust
# IPv6 local connections:
host    all             all             ::1/128                 trust
# Allow replication connections from localhost, by a user with the
# replication privilege.
local   replication     all                                     trust
host    replication     all             127.0.0.1/32            trust
host    replication     all             ::1/128                 trust

```

##### After modifying your pg_hba.conf, you will want to restart PostgreSQL so that the changes you made are reflected.
##### For Mac Users:
```
brew services restart postgresql
```

##### Log back into PostgreSQL
```
psql -U <your_user_system_name> -d postgres
```

##### Here is an example
```
psql -U kevincook -d postgres
```

### Step 3: Create a Database for the project
##### Run the following commands in the PostgreSQL (psql) shell
```
CREATE DATABASE your_database_name OWNER your_project_user;
```

##### Here is an example
```
CREATE DATABASE testdb OWNER kevincook;
```

##### Confirm database was created by running this command:
```
\l
```

##### Switch to newly created database
```
\c your_database_name
```

##### Here is an example
```
\c testdb
```

##### Exit PostgreSQL (psql) shell
```
\q
```

##### Move to your server directory then run the following command to run the database initialization script to create tables:
```
cd server
psql -U <your_postgres_user> -d <your_database_name> -f init_db.sql
```

##### Here is an example:
```
cd server
psql -U kevincook -d testdb -f init_db.sql
```

### Step 4: Set Up Backend
##### Install dependencies for the project:
```
go mod tidy
```

##### Next you want to copy the example env file and rename it:
```
cp .env.example .env
```

##### Move back to the TrackPointGPS directory and open the code in your text editor:
```
cd ..
code .
```

##### Open the .env and update the following with your database credentials:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_postgres_user
DB_NAME=your_database_name
API_KEY=your_api_key_here  -- This was sent via email
```

##### Here is an example
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=kevincook
DB_NAME=testdb
API_KEY=your_api_key_here  -- This was sent via email
```

##### Lastly, open your server terminal and start the backend server with this command:
```
go run main.go
```

### Step 5: Set Up Frontend
##### Open your client terminal and install dependencies for the project:
```
npm install
```

##### Run the development server:
```
npm run serve
```

## Additional Notes
### Lints and fixes files
```
npm run lint
```