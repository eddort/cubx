# ibox

<img src="./ibox.png" height="200px" align="right" width="200px">

## Introduction
ibox is a tool designed to simplify the deployment and testing of applications by running them inside Docker containers. It abstracts the complexity of Docker commands into a simple interface, allowing developers to launch applications "out of the box".
- **Simple Command Line Interface**: Run your applications with simple commands.
- **Port Mapping Support**: Easily map ports from the container to the host.
- **Volume Management**: Quickly define and manage volumes for persistent data.
- **Lightweight and Fast**: Minimal overhead and speedy execution of containers.

## Getting Started

### Prerequisites
Ensure you have Docker installed on your machine. ibox interfaces directly with Docker, so it's required for operation.

### Usage
To run an application using ibox:

```sh
ibox node test.js
```

```sh
ibox npm install
```