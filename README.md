# cubx

<img src="./docs/cubx.png" height="300px" align="right" width="300px">

**Cubx** is a versatile tool that simplifies running console programs inside Docker containers. It is suitable for anyone who wants to use applications without the need to install them on their device, providing ease and convenience in operation..
- **Simple Command Line Interface**: Run applications with simple commands as if they were installed locally.
- **Port Mapping Support**: All application ports are automatically mapped to the host.
- **Lightweight and Fast**: Minimal overhead and speedy execution of containers.
- **Automatic Working Directory Mounting**: Automatically mounts the working directory as if you were running the program locally.
- **Automatic Platform Detection**: Automatically detects your OS and processor architecture to download the most suitable container image.
- **Flexible Software Versioning**: Use any version of the software by specifying the image tag after a colon. For example: `cubx node:14 test.js`, `cubx npm:14 install`, `cubx yarn:14 add [package]`.

## Getting Started

### Prerequisites
Ensure you have Docker installed on your machine. **cubx** interfaces directly with Docker, so it's required for operation.

#### MacOS

If you are using macOS, you need to enable the net=host feature in Docker Engine. This feature allows port mapping to the host. While this functionality is enabled by default on Linux, it is still in beta on macOS and requires manual activation.

For more information, visit the following link: https://docs.docker.com/network/drivers/host/#docker-desktop

### Installation
You need to clone the project repository:

```bash
git clone git@github.com:eddort/cubx.git
```

Then, navigate to the project directory:

```bash
cd cubx
```

Run the following command to build the project:

```bash
go build
```

Now, you can use the command:

```bash
./cubx
```
> You can then specify a `PATH` variable in your environment to make `cubx` available globally

This will allow you to start using the `cubx` tool.


### Usage

To display the available programs, enter

```sh
cubx -h
```
Output:

![cubx help](./docs/help-example.png)

To run an application using cubx:

```sh
cubx node --eval 'console.log("Hello from Cubx\n",`node version: ${process.version}`)'
```

Output:

```sh
Hello from Cubx
 node version: v22.1.0
```

We can also specify a specific version of the program with ":"

```sh
cubx node:14 --eval 'console.log("Hello from Cubx\n",`node version: ${process.version}`)'
```

Output:

```sh
Hello from Cubx
 node version: v14.21.3
```

```sh
cubx npm install
```

```sh
cubx cast call 0x6b175474e89094c44da98b954eedeac495271d0f 'totalSupply()(uint256)' --rpc-url https://eth-mainnet.alchemyapi.io/v2/Lc7oIGYeL_QvInzI0Wiu_pOZZDEKBrdf
```
