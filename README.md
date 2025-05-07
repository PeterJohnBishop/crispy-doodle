# crispy-doodle

# 
docker build -t peterjbishop/crispy-doodle:latest .
docker push peterjbishop/crispy-doodle:latest
docker pull peterjbishop/crispy-doodle:latest
docker-compose down
docker-compose build --no-cache
docker-compose up 