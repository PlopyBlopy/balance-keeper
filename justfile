set dotenv-path := "./docker/.env"
set windows-shell := ["powershell.exe", "-NoProfile", "-c"]

# create containers applies migrations on database and remove migration container
docker-up:
    docker-compose -f ./docker/docker-compose.yml up --abort-on-container-exit balance-keeper-migrate
    docker-compose -f ./docker/docker-compose.yml rm -f balance-keeper-migrate

# down services and remove volumes
docker-down-volumes:
    docker-compose -f ./docker/docker-compose.yml down --volumes

# for migrations, use container balance-keeper-migrate
migrate cmd *args:
    docker-compose -f ./docker/docker-compose.yml run --rm balance-keeper-migrate \
        -path /migration  \
        -database {{ env_var('POSTGRES_DB_URL') }} \
        {{ cmd }} {{ args }}
