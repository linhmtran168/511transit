# 511 Transit Data Search

## Overview

Real time transit data for SF Bay area based on Transit Data Feed API provided on [https://511.org/](https://511.org/)
Demo: [https://511transit.fly.dev/](https://511transit.fly.dev/)

## Run locally

- Requirements
  - Install Docker
  - Get an API Key from [511 Request Token](https://511.org/open-data/token)

- Run the following command in the root folder:

  ```bash
  docker build -t 511transit .
  docker run -p 8080:8080 -e API_KEY=${requested_api_key} 511transit
  ```

## Deployment

- Since this project have no external dependencies, it can be deployed with a single Docker image (built using provided Dockerfile) to any platform that support running container
- Just need to set the `API_KEY` environment variable to the API key you got from [511 Request Token](https://511.org/open-data/token)
- The project comes with `.fly.toml` for default settings when deploy to [Fly.io](https://fly.io). Highly recommended this platform for simple website and workload

## Development

### Requirements

- Go 1.18 or higher
- Node.js 16 or higher
- Protocol Buffer Compiler (3.12 or higher)
  - Only if you need to re-generate GTFS protobuf files
- After install the above requirements, run the following command in the root folder to install necessary Go and npm dependencies:

  ```bash
  make install-deps
  ```

### Run development server

- Create `.env.local` file at the root of the project with the following content:

  ```conf
  API_KEY=<511 API key>
  SERVER_ADDR=":5050" # Only needed if you want to change the backend default listening port of 8080
  ENV=local
  ```

- Backend (Hot reload function provided by [Air](https://github.com/cosmtrek/air))

  ```bash
  make back-dev
  ```

- Frontend (HMR enabled development server provided by [Vite](https://vitejs.dev/))

  ```bash
  make front-dev
  ```

## General Architecture

### Backend

- The backend is built with [Chi](https://github.com/go-chi/chi) and provide a HTTP server that listen on port `8080` by default.
- `/` route is the root route that served Frontend compiled file from `./web/dist` folder.
  - This route is available only in `production` mode. In `development` mode, the static files are served by Vite dev server
- `/ws` route provide a WebSocket server that response to request messages and return data from 511 Transit API in accordance with the request type. The format of the request and response messages are as below:
  - `{ "requestType": "operators"}`
    - Response
      - `{"responseType":"operators","data":[{"id":"<id>","name":"<name>"}, ..]}`
  - `{ "requestType": "tripUpdates", "operatorIds": "<operatorIds>"}`
    - Response:
      - `{"responseType":"tripUpdates","data": { "operatorId": "<id>", "tripUpdates": [...]}`

- The Backend requests 511 Transit API to get the data. Because 511 Transit API has a rate limit and to provide better response time, the data from the API is cached in memory.
  - Operators data: 1 hours
  - Trip updates data: 1 minutes

### Frontend

- A SPA app built with Vite and Vue 3.x, communicate with the backend using WebSocket
  - Automatically reconnect after connection lost
- In local env, the frontend is served by Vite dev server, default run on port `:5173`
- In production, the frontend codes is compiled to `web/dist` folder and served by the static server provided by the backend
  - Run `make front-build` to execute this process

### Folder structure

```text
- cmd
  - transit.go # Main entrypoint for backend code
- internal # All backend logic codes
  - data
    - memory # Hold code related to caching and query data
  - handler # HTTP handler
  - mocks # Mockery generated mocks
  - models # # Hold domain objects, interfaces
- web # All frontend codes live here
  - dist # Compiled codes for production
  - public # Public assets
  - src # SPA code
    - components # Vue components
    - composables # Vue composables
    - store # Pinia store
    - index.html # Main index page
    - main.ts # Entrypoint for frontend
```

### Test

- Currently only unit test for caching part of the Backend is provided. Run them using the following command:

  ```bash
  make back-dev
  ```

## To do

- [ ] Add more tests
  - [ ] Frontend
  - [ ]Integration tests for backend
- [ ] `Ping pong` check between the frontend and backend to check the status of websocket connection and keep it alive
- [ ] Add option to use better in-memory DB to hold the cache data (like Redis) instead of using Go lang objects right now
