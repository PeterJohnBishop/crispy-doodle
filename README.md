# crispy-doodle

## backend

- Gin/Http API with JWT authentication
- Postgres Database with Users table, Messages table
- OpenAI API
- AWS/S3 API for media storage

# run in docker
docker build -t peterjbishop/crispy-doodle:latest .
docker push peterjbishop/crispy-doodle:latest
docker pull peterjbishop/crispy-doodle:latest
docker-compose down
docker-compose build --no-cache
docker-compose up 

## open postgres container
docker exec -it postgres psql -U postgres -d postgres

