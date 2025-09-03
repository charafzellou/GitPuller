# GitPuller - a Github Organization Repository Cloner

This is a Go application that uses the Github API to clone all the repositories in a given Github organization, and pull updates for all the cloned repositories.

## Getting Started:

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites:

Before you can use this application, you need to have the following:

- Go 1.21 or higher
- Docker 23.0 or higher
- A Github account
    - An active PAT [(Personal Access Token)](https://github.com/settings/tokens) with the `repo` scope

> **⚠️** Note that **the PAT is optional** if the repositories you want to pull are public. If not, you will need to create one and add it in the `.env` file.

### Setting up:

To get started, clone this repository to your local machine:

```bash
git clone https://github.com/charafzellou/gitpuller
```

Change into the project directory:

```bash
cd gitpuller
```

Create an environment file `.env` with your Github organization name and personal access token:

```bash
ORG_NAME=<your-github-organization-name>
PERSONAL_TOKEN=<your-github-access-token>
```

### Usage as a binary:

To clone all the repositories in the organization using the Go binary directly, run the following command:

```bash
make build
make run
```

The program will detect the repositories, and apply a `git pull --all` on them.

### Usage as a Docker container:

Build the Docker image with the following command:

```bash
make build . -t gitpuller
```

Run the Docker container from the image we created:

```bash
make run gitpuller
```

This will start the application and clone all the repositories in the specified organization to a folder called `repositories` in the current directory.

## Authors:

- **Charaf ZELLOU**
    - [Github](https://github.com/charafzellou)
    - [LinkedIn](https://linkedin.com/charafzellou)

## Collaborators:

- **Your name could be here, make a useful PR 😊**

## Built With:

- [Go](https://golang.org/) - The programming language used
- [Docker](https://www.docker.com/) - Containerization platform

## License:

This project is licensed under the AGPL-3.0 License - see the [LICENSE](./LICENSE.md) file for details.