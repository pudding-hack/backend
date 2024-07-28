
<h1 align="center">
  <br>
  Sightry
  <br>
</h1>

<h4 align="center">Inventory for all</h4>
Sightry is a mobile application that focuses on helping blind people manage stock, such as checking stock, adding, and reducing stock based on image rekognition supported by AWS recognition. 
<br>
<br>


<p align="center">
  <a href="#key-features">Key Features</a> •
  <a href="#how-to-use">How To Use</a> •
  <a href="#author">Author</a>
</p>


## Key Features

* Inventory Management
* Stock Management
* Image Recognition
* Authentication
  
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

#### Notes : 
On the docker-compose.yaml file you would also need to add these following environment variable on the Inventory Service to be able to use the image rekognition service from AWS.

```bash
AWS_ACCESS_KEY_ID=<aws access key ID>
AWS_SECRET_ACCESS_KEY=<aws access key>
AWS_REGION=<aws region>
```

But, we hosted a demo api for testing purposes https://puddinghack.varomnrg.me/api


## Author

- [Timothy Aurelio Cannavaro](https://github.com/varomnrg) - Backend
- [Nanda Wijaya Putra](https://github.com/nanwp) - Backend
- [Muhammad Alif Vidi](https://github.com/MuhammadAlifVidi) - Mobile
- [Nadira Belinda]() - Hipster
