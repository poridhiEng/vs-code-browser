tag := v1.3
repo := poridhi
image := codeserver-reverse-proxy

build:
	@ echo "Building Docker image for repository ${repo}, image ${image}, tag ${tag}..."
	@ docker build --platform linux/amd64 -t ${repo}/${image}:${tag} -f Dockerfile .

run:
	@ echo "Starting Docker container for repository ${repo}, image ${image}..."
	@ docker run --name ${image} -p 8080:8080 ${repo}/${image}:${tag}

clean:
	@ echo "Stopping and removing Docker container for repository ${repo}, image ${image}..."
	@ docker stop ${image}
	@ docker rm ${image}

push:
	@ echo "Pushing Docker image for repository ${repo}, image ${image}, tag ${tag} to the registry..."
	@ docker push ${repo}/${image}:${tag}
	@ echo "Image ${repo}/${image}:${tag} successfully pushed to the registry."