

Useful commands

```bash
docker-compose up -d # Run docker
docker-compose down # Stop docker
docker ps # Get docker list
docker logs go-gotracker_sqldb_1 # Get logs
watch docker logs go-gotracker_sqldb_1 # Watch logs

# Mysql
docker exec -it <container_id> mysql -u root -p # SQL interface

# Postgres
docker exec -it go-gotracker_sqldb_1 sh # Open shell
createdb -U postgres go-gotracker # Create database
\l # List databases
\c <database> # Change to db
\dt # List tables
```
