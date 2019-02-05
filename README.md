# go-shop-item


## Description
This is an example of implementation of Shop Item in Go (Golang) projects.



### How To Run This Project
> Make Sure you have run the shop.sql in your mysql

```bash
#move to directory
cd $GOPATH/src/github.com

# Clone into YOUR $GOPATH/src
git clone https://github.com/hendriksiallagan/go-shop-item.git

#move to project
cd go-shop-item

# Install Dependencies
dep ensure

# Test the code
make test

# Run Project
go run main.go

```
Or With `go get`
> Make Sure you have run the shop.sql in your mysql

```bash
# GET WITH GO GET
go get github.com/hendriksiallagan/go-shop-item.git

# Go to directory

cd $GOPATH/src/github.com/go-shop-item

# Install Dependencies
dep ensure

# Test the code
make test

# Run Project
go run main.go
```

Or with `docker-compose`

```bash
#move to directory
cd $GOPATH/src/github.com

# Clone into YOUR $GOPATH/src
git clone https://github.com/hendriksiallagan/go-shop-item.git

#move to project
cd go-shop-item

# Build the docker image first
make docker

# Run the application
make run

# check if the containers are running
docker ps

# Execute the call
curl localhost:9090/items

# Stop
make stop
```


### Notes:
In this project, you can access some links :

- [GET] /items -> to display a list data of items
- [GET] /calculate -> to display how much price subtotal, tax subtotal & grand total
- [POST] /items -> to insert data item : name, price, taxcode to database

