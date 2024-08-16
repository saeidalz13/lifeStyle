# DOCKER
create_network:
	docker network create lifestyle_network
	
run_container:
	docker run --network lifestyle_network --name lifestyle_auth -p 8765:5432 --env-file "./config/.docker.env" -d postgres:16.2 

create_database:
	docker exec lifestyle_auth psql -U lifestyleauthuser -c "create database lifestyle_auth;" 