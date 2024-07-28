
<h1 align="center">
  <br>
  Sightry
  <br>
</h1>

<h4 align="center">Inventory for all</h4>

<p align="center">
  <a href="#key-features">Key Features</a> •
  <a href="#how-to-use">How To Use</a> •
  <a href="#author">Author</a>
</p>


## Key Features

* Check Inventory
* Stock Management

## How To Use

To clone and run this application, you'll need [Git](https://git-scm.com), [Docker](https://www.docker.com/), [Golang](https://go.dev/), [PostgreSQL](https://www.postgresql.org/), installed on your computer, From your command line:

Clone this repository and go inside the folder
```bash
$ git clone https://github.com/pudding-hack/backend
```

Run the following command to migrate the database (you would need your own local/hosted database)

```bash
$ migrate -path database/migration/ -database "postgresql://<username>:<password>@<hostname>:5432/sightry?sslmode=disable" up
```


Setup your docker compose file and replace the environment variables with your own then run the following command
```bash
$ docker-compose up -d --build
```

When docker is running, you can access the application on the port you set up in the docker-compose file, the default is http://localhost:8081

## Author

- [Timothy Aurelio Cannavaro](https://github.com/varomnrg) - Backend
- [Nanda Wijaya Putra](https://github.com/nanwp) - Backend
- [Muhammad Alif Vidi](https://github.com/MuhammadAlifVidi) - Mobile
- [Nadira Belinda]() - Hipster