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
```

## Setting Up The Project
### Step 1: Clone the Repository
```
git clone https://github.com/Kevin-T-Cook/TrackPointGPS.git
cd TrackPointGPS
```

### Step 2: Set Up Backend
##### In the server, install dependencies for the project:
```
go mod tidy
```

##### Next you want to copy the example env file and rename it:
```
cp .env.example .env
```

##### Open the .env and update the following with your database credentials:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_postgres_user
DB_PASSWORD=your_postgres_password
DB_NAME=your_database_name
API_KEY=your_api_key_here  -- This was sent via email
```

##### Then you will run the database initialization script to create tables:
```
psql -U <your_postgres_user> -d <your_database_name> -f init_db.sql
```

##### Here is an example:
```
psql -U kevincook -d onestepdb -f init_db.sql
```

##### Lastly, you can start the backend server with this command:
```
go run main.go
```

### Step 3: Set Up Frontend
##### Install dependencies for the project:
```
npm install
```

##### Run the development server:
```
npm run serve
```

## Lints and fixes files
```
npm run lint
```