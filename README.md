codeserver-reverse-proxy:


# Run at local environment [ Using Docker ]

1. Pull the docker image to host.
```bash
docker pull poridhi/codeserver-reverse-proxy:v1.3
```

2. Run the container.
```bash
docker run --name codeserver-reverse-proxy -p 8080:8080 poridhi/codeserver-reverse-proxy:v1.3
```

3. Open a new terminal tab and get the container ip

```bash
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' codeserver-reverse-proxy
```

4. Open another new terminal tab and run

```bash
docker exec -it network-tools-container sh
```
5. Now curl to get the logs 
```bash
curl http://172.17.0.2:3000/ns1/?folder=app
```