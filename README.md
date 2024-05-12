# ibox

<img src="./docs/ibox.png" height="200px" align="right" width="200px">

**ibox** is a tool designed to simplify the deployment and testing of applications by running them inside Docker containers. It abstracts the complexity of Docker commands into a simple interface, allowing developers to launch applications "out of the box".
- **Simple Command Line Interface**: Run applications with simple commands as if they were installed locally.
- **Port Mapping Support**: Out of the box, all application ports are automatically mapped to the host.
- **Lightweight and Fast**: Minimal overhead and speedy execution of containers.
- **Automatic Working Directory Mounting**: Automatically mounts the working directory as if you were running the program locally.
- **Automatic Platform Detection**: Automatically detects your OS and processor architecture to download the most suitable container image.
- **Flexible Software Versioning**: Use any version of the software by specifying the image tag after a colon. For example: `node:14 test.js`, `npm:14 install`, `yarn:14 add [package]`.

## Getting Started

### Prerequisites
Ensure you have Docker installed on your machine. **ibox** interfaces directly with Docker, so it's required for operation.

### Usage

To display the available programs, enter

```sh
ibox -h
```
Output:

![ibox help](./docs/help-example.png)

To run an application using ibox:

```sh
ibox node test.js
```

```sh
ibox npm install
```

```sh
ibox cast call 0x6b175474e89094c44da98b954eedeac495271d0f 'totalSupply()(uint256)' --rpc-url https://eth-mainnet.alchemyapi.io/v2/Lc7oIGYeL_QvInzI0Wiu_pOZZDEKBrdf
```
