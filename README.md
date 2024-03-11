### VS Code Editor on Browser

Resources: [Code Server Installation Guide](https://coder.com/docs/code-server/latest/install#docker)

- Step 1: Start the vscode server using `Dockerfile.base`

```bash

docker build -t demo1:v1 -f Dockerfile.base .

docker run --rm -p 8080:8080 demo1:v1

```

Visit localhost:8080

- Step 2: Login into the system from browser. Install necessary configuration for your coding environment.

- Step 3: Copy the `/root/.config` and `/root/.local` diretories to host

```bash

docker cp <container_id>:/path/in/container /path/on/host

# For example,

docker cp 4c432caca6f9:/root/.config ./config
docker cp 4c432caca6f9:/root/.local ./local

```

Now, we will create an image that will have these configs predefined as soon as the container starts

- Step 4: Now, Build another image using docker file `Dockerfile.python` 

```bash

docker build -t demo2:v1 -f Dockerfile.python .

docker run --rm -p 8080:8080 demo2:v1

```

- Step 6: Go to url `http://localhost:3000/?folder=/app`