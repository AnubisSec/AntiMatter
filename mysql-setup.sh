# since i don't know how to use dockerfiles yet, this is the best I can do


# Pull the latest mysql image
docker pull mysql/mysql-server:latest

# Run the image naming it "anti_sql" making it expose ports 3306 and 22 with ports 3307 and 3308 respectively
docker run --name=anti_sql -p 3307:3306 -p 3308:22 -e MYSQL_ROOT_PASSWORD=Passw0rd! -e MYSQL_DATABASE=Anti --restart on-failure -d mysql/mysql-server:latest

