# DirWatcher

This project monitors a specified directory for changes and counts occurrences of a magic string in files.

## Prerequisites

- [Go](https://golang.org/) 
- [PostgreSQL](https://www.postgresql.org/) (you can use any other database)


## Setup

- Clone the repository `git clone https://github.com/yourusername/background-task-monitor.git`

- Navigate to the project directory

- Install dependencies
```bash
$ go get github.com/gin-gonic/gin
$ go get github.com/fsnotify/fsnotify
```

## Database Setup
- Create a PostgreSQL database and note the connection URL

- Replace "yourusername", "yourpassword", "yourdbname" with your actual database credentials

## Running the Application
```bash 
go run .
```

## Monitor the Results

* Check the configuration
  - Send a GET request to http://localhost:8080/config

* Check the task result
  - Send a GET request to http://localhost:8080/task/result

## Stopping the Application

- Press Ctrl + C in the terminal where the API server is running


# API Endpoints

| Route               | Method | Description                           | Request Body (JSON)                        | Sample Response (JSON)                        |
|---------------------|--------|---------------------------------------|--------------------------------------------|----------------------------------------------|
| `/task/start`       | POST   | Start the background task             | None                                       | `{"message": "Background task started"}`     |
| `/config`           | GET    | Get the current configuration         | None                                       | `{"directory": "/path/to/directory", "time_interval": "5s", "magic_string": "magic", "task_in_progress": false}` |
| `/config`           | POST   | Update the configuration              | `{"directory": "/new/path", "time_interval": "10s", "magic_string": "updated", "task_in_progress": true}` | `{"message": "Configuration updated successfully"}` |
| `/task/result`      | GET    | Get the result of the background task | None                                       | `{"start_time": "2024-01-30T01:01:33Z", "end_time": "2024-01-30T01:01:38Z", "runtime": "5s", "files_added": ["file1.txt"], "files_deleted": ["file2.txt"], "magic_string_cnt": 10, "status": "success"}` |

## Sample Usage

1. **Start Background Task:**

   - **Route:** `/task/start`
   - **Method:** `POST`
   - **Sample Response:**
     ```json
     {"message": "Background task started"}
     ```

2. **Get Current Configuration:**

   - **Route:** `/config`
   - **Method:** `GET`
   - **Sample Response:**
     ```json
     {"directory": "/path/to/directory", "time_interval": "5s", "magic_string": "magic", "task_in_progress": false}
     ```

3. **Update Configuration:**

   - **Route:** `/config`
   - **Method:** `POST`
   - **Request Body:**
     ```json
     {"directory": "/new/path", "time_interval": "10s", "magic_string": "updated", "task_in_progress": true}
     ```
   - **Sample Response:**
     ```json
     {"message": "Configuration updated successfully"}
     ```

4. **Get Task Result:**

   - **Route:** `/task/result`
   - **Method:** `GET`
   - **Sample Response:**
     ```json
     {"start_time": "2024-01-30T01:01:33Z", "end_time": "2024-01-30T01:01:38Z", "runtime": "5s", "files_added": ["file1.txt"], "files_deleted": ["file2.txt"], "magic_string_cnt": 10, "status": "success"}
     ```

