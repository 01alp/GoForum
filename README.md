# Forum

## Description

This project consists in creating a web forum that allows :

    - communication between users.
    - associating categories to posts.
    - liking and disliking posts and comments.
    - filtering posts.

## Author
01alp
## Usage
  
To run the web-site on local machine:

- Download the repository

- Run with a command `./run.sh`

- Open [http://localhost:8080/](http://localhost:8080/) in browser

- Register or login with a test user(test@gmail.com, 1234)

- Stop with a command `./stop.sh`

To work with a Docker manually:

- Build image with `sudo docker build -t application .`

- Run image with `docker run -p 8080:8080 application:latest`

## Implementation details

- SQL3

- Adaptive web-design
