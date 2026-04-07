TARGET = cat-database
IMAGE_NAME = cat-database
CONTAINER_NAME = cat-database
VOLUME_NAME = cat-database

build-image:
	podman build -t $(IMAGE_NAME) --format docker .

clean:
	podman stop $(CONTAINER_NAME) || true
	podman rm -f $(CONTAINER_NAME) || true
	podman rmi -f $(IMAGE_NAME) || true

run-container:
	podman run -d -v $(VOLUME_NAME):/app/data --network=host --name $(CONTAINER_NAME) --replace $(IMAGE_NAME)

.PHONY: build-image clean run-container
