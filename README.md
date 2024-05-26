# cubx

<img src="./docs/cubx.png" height="300px" align="right" width="300px">

**Cubx** is a versatile tool that simplifies running console programs inside Docker containers. It is suitable for anyone who wants to use applications without the need to install them on their device, providing ease and convenience in operation.
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
> [!NOTE]
> You can then specify a `PATH` variable in your environment to make `cubx` available globally

This will allow you to start using the `cubx` tool.


## Usage
Cubx has a set of preset programs, using the command you can check the list of available programs.
To display the available programs, enter:

```sh
cubx -h
```
Output:

![cubx help](./docs/help.png)

## Basic Usage

Cubx's main goal is to provide the easiest possible way for applications inside containers. Everything looks like you are running them locally.

```sh
cubx node --eval 'console.log(`node version: ${process.version}`)'
```
The command analog without the Cubx:

```sh
node --eval 'console.log(`node version: ${process.version}`)'
```

Output:

```sh
node version: v22.1.0
```

The output may be different only if you don't have Node.js installed or if you don't have the latest version of it.

As you can see, the way to use it is almost no different from running a regular program in a terminal. Only by running programs with cubx, you run programs in an isolated environment and get a lot of convenience and means to secure your data.

- the application runs in a separate container
- only the current volume is mounted, your data from other directories is safe
- If you need to restrict access to folders or files in the working directory, you can do that too.
- you can disconnect the application from the local network or from the entire Internet.

## Features
### Version control
In addition to security, cubx provides a user-friendly interface for working with applications. For example, working with versions. nodejs has nvm (but not all programs have an analog). Cubx provides a version control mechanism that works as simply as possible.
Just add ":version" to your program at the terminal prompt.

```sh
cubx node:14 --eval 'console.log(`node version: ${process.version}`)'
```

Output:

```sh
node version: v14.21.3
```

Any other program works the same way, as long as it has a Docker-registry (all possible Docker-registries are supported).

### Interactive version selection

Sometimes we don't know what specific versions are out there right now and we don't want to search the internet for the exact spelling of the version. With the `--select` flush we can activate the interactive version selection interface.

```sh
cubx --select node --eval 'console.log(`node version: ${process.version}`)'
```


Select:

![cubx select](./docs/select.png)

Result:

![cubx select result](./docs/select-result.png)

### File Exclusion

It is a fairly common task to restrict access (visibility) of certain files and folders. You can store locally an .env file with keys that you are afraid of losing access to.

As you know, when installing dependencies in Node.js, you can compromise your data due to unscrupulous third-party code.

> This situation can happen not only in Node.js, but in any platform that executes code downloaded from the Internet.

Let's demonstrate how this works with an example:

`test-env.js`
```js
const fs = require('fs').promises;
const path = require('path');

async function readEnvFile() {
    try {
        const envFilePath = path.resolve(__dirname, '.env');
        const data = await fs.readFile(envFilePath, 'utf8');
        console.log(data);
    } catch (err) {
        console.error('Error reading .env file:', err);
    }
}

readEnvFile();

```
`.env`
```env
SOME_PRIVATE_KEY=123
```
Let's create a file and write the code in JS. We will read the local file with cofiguration and output it to the console. As if it were a malicious script.

Next, let's call the script:

```sh
cubx node test-env.js 
```

Output

```sh
SOME_PRIVATE_KEY=123
```

Our script read the configuration without problems and output everything to the console.

Let's now exclude the file and try again.

```sh
cubx --ignore-path .env node test-env.js 
```

Output

```sh

```

As a result we get empty output, because inside the container this file will be empty when the program is called.

```sh
cubx npm install
```

```sh
cubx cast call 0x6b175474e89094c44da98b954eedeac495271d0f 'totalSupply()(uint256)' --rpc-url https://eth-mainnet.alchemyapi.io/v2/Lc7oIGYeL_QvInzI0Wiu_pOZZDEKBrdf
```
