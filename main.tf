terraform {
    required_providers {
      docker = {
          source = "kreuzwerker/docker"
          version = "~> 2.13.0"
      }
    }
}

provider "docker" {}

resource "docker_network" "golang-todos" {
  name = "golang-todos"
}

resource "docker_image" "mongo" {
    name = "stevenweatherspoon/todos-mongodb:latest"
    keep_locally = false
}

resource "docker_container" "mongo" {
    image = docker_image.mongo.latest
    name = "mongodb"
    ports {
        internal = 27017
        external = 27017
    }
    networks_advanced {
      name = "golang-todos"
    }
}

resource "docker_image" "golang-todos" {
  name         = "stevenweatherspoon/golang-todos:latest"
  keep_locally = false
}

resource "docker_container" "golang-todos" {
    name = "golang-todos"
    image = docker_image.golang-todos.latest
    ports {
        internal = 8080
        external = 8080
    }
    networks_advanced {
      name = "golang-todos"
    }
}