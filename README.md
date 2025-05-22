 # QR Code Generator

 Go + Vue.js application for generating customizable QR codes with optional SVG logo overlays.

 ## Development

 ### Frontend

 ```bash
 cd frontend/frontend-app
 npm install
 npm run dev
 ```

 Starts the frontend dev server at http://localhost:5173 with hot reloading and API proxy to http://localhost:8080.

 ### Backend

 ```bash
 go run main.go
 ```

 Starts the backend server at http://localhost:8080.

 ## Docker

 Build the Docker image:

 ```bash
 docker build -t qr-generator .
 ```

 Run the container:

 ```bash
 docker run -p 8080:8080 qr-generator
 ```

 Access the application at http://localhost:8080.