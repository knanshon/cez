# Go HTMX Todo Application

This is a simple Todo application built with Go and HTMX. The application allows you to view and add Todo items using a web interface.

## Launch Instructions

1. Ensure that [Go is installed](https://golang.org/dl/).
2. From the project root, run:
    ```sh
    go run main.go
    ```
3. The server will start on port `9090`.

## Usage Instructions

1. Open your web browser and navigate to [http://localhost:9090](http://localhost:9090).
2. The page displays a list of todos.
3. To add a new todo, fill in the "Message" field and hit the **Create** button. The new todo will be appended to the list using HTMX.

## Project Structure

- [main.go](/Users/barrynorthern/personal/cez/main.go) &ndash; Contains the HTTP server implementation and handler functions.
- [index.html](/Users/barrynorthern/personal/cez/index.html) &ndash; The HTML template that displays the todo list and form.
- [LICENSE](/Users/barrynorthern/personal/cez/LICENSE) &ndash; License details.

Enjoy using the application!