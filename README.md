# Docker Mailserver Aliases

## Overview

This project is an addon for the **[Docker Mailserver](https://github.com/docker-mailserver/docker-mailserver)** (DMS) project. It provides a simple web interface to manage email aliases. This addon is packaged as a Docker container that hosts a REST API written in Go and a frontend built with Svelte, Tailwind CSS, and DaisyUI.

## Features

- List existing mail aliases and the email address they redirect to.
- Add new aliases.
- Delete existing aliases.

## Technologies

- [Gin Web Framework](https://gin-gonic.com/)
- [Docker Engine SDK](https://pkg.go.dev/github.com/docker/docker/client)
- [Svelte](https://svelte.dev/)
- [Tailwind CSS](https://tailwindcss.com/)
- [DaisyUI](https://daisyui.com/)

## Installation

### Prerequisites

To run the Docker container, you will need:

- A working and configured [Docker Mailserver](https://github.com/docker-mailserver/docker-mailserver) instance.
- (Optional) A reverse proxy, e.g., [Caddy](https://caddyserver.com/), to serve the frontend over HTTPS.

### Configuration

You can change the port of the web server with the following environment variable:

```bash
export GIN_ADDR=":8080"
```

### Docker Compose

Add the `mailserver-aliases` container to your `docker-compose.yaml` file:

```yaml
services:
  mailserver-aliases:
    image: chscheid/docker-mailserver-aliases:1.0.1
    restart: unless-stopped
    ports:
      - "8080:8080"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
```

> **Note**: There is no built-in authentication for the frontend, so the web interface will be publicly available under port 8080. You may want to secure it with a reverse proxy and authentication.

Here is an example to serve the frontend with Caddy and Basic Authentication:

`docker-compose.yaml`:

```yaml
services:
  caddy:
    image: caddy:2.8
    restart: unless-stopped
    environment:
      - FRONTEND_BASIC_AUTH_USER="username"
      - FRONTEND_BASIC_AUTH_PASSWORD="HASHED_PASSWORD"
    ports:
      - "80:80"
      - "443:443"
      - "443:443/udp"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile

  mailserver-aliases:
    image: chscheid/docker-mailserver-aliases:1.0.1
    restart: unless-stopped
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
```

`Caddyfile`:

```
aliases.yourdomain.com {
  basic_auth {
    {$FRONTEND_BASIC_AUTH_USER} {$FRONTEND_BASIC_AUTH_PASSWORD}
  }
  reverse_proxy mailserver-aliases:8080
}
```

Replace `username` and `HASHED_PASSWORD` with your values. For more information on configuring Caddy and hashing the password, see the [Caddy documentation](https://caddyserver.com/docs/caddyfile/directives/basic_auth).

## Development

To develop and contribute to this project, you can run both the backend and frontend locally. You can mock the API for frontend development using [Mockoon](https://mockoon.com/).

### REST API

The backend REST API is written in Go and provides endpoints to manage aliases. It interacts with the Docker Mailserver instance through the Docker Engine API.

#### Prerequisites

- Go (1.23 or higher)

#### Run Development Server

To start the development server for the backend:

```bash
go run main.go
```

#### Swagger Documentation

The Swagger documentation is generated with [Swag](https://github.com/swaggo/swag).

Generate the documentation with:

```bash
swag init
```

You can view the documentation at:

[http://localhost:8080/docs/index.html](http://localhost:8080/docs/index.html)

### Frontend

The frontend is built with Svelte, Tailwind CSS, and daisyUI. It communicates with the backend REST API to manage email aliases.

#### Prerequisites

- Node.js (20 or higher)

#### Run Development Server

To run the frontend locally, follow these steps:

1. Navigate to the `frontend` directory:

   ```bash
   cd ./frontend
   ```

2. Install the necessary dependencies:

   ```bash
   npm install
   ```

3. Use Mockoon to mock the REST API if the Docker Mailserver is not running:

   ```bash
   npm run mockoon
   ```

4. Start the frontend development server:

   ```bash
   npm run dev
   ```

This will start the frontend on [http://localhost:5173/](http://localhost:5173/).

## Contributing

Contributions are welcome! If you'd like to contribute to the project, please follow these steps:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes and commit them (`git commit -am 'Add new feature'`).
4. Push the branch (`git push origin feature-branch`).
5. Open a pull request.

Please make sure to update tests as appropriate.

## License

This project is open-source and available under the [MIT License](LICENSE).