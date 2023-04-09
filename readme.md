

<h3 align="center">Customer Rest API <br>
<h5 align="center" >Golang Gin MVC<h5>
<br>
</h4>
<p align="left">
<h2>
  Content <br></h2>
  • Key Features <br>
  • Installing Using Github<br>
  • Installing Using Docker<br>
  • End Point<br>
  • Technologi that i use<br>
  • Contact me<br>
</p>

## 📱 Features

* Auth
* Customers
* Products
* Orders


## ⚙️ Installing and Runing from Github

installing and running the app from github repository <br>
To clone and run this application, you'll need [Git](https://git-scm.com) and [Golang](https://go.dev/dl/) installed on your computer. From your command line:

```bash
# Clone this repository
$ git clone https://github.com/reski-id/costumer.git

# Go into the repository
$ cd costumer

# Install dependencies
$ go get

# Run the app
$ go run main.go

# if you have problem while running you can use bash cmd and type this..
$ source .env #then type 
$ go run main.go #again
```

> **Note**
> Make sure you allready create database mysql `costumerdb` for this app.more info in local `.env` and `utils/database.go` file.


## ⚙️ Installing and Runing with Docker
if you are using docker or aws/google cloud server you can run this application by creating a container. <br>

```bash
# Pull this latest app from dockerhub 
$ docker pull programmerreski/costumer

# if you have mysql container you can skip this
$ docker pull mysql

$ docker run --name mysqlku -p 3306:3306 -d -e MYSQL_ROOT_PASSWORD="yourmysqlpassword" mysql 

# create app container
$ docker run --name costumer -p 80:8000 -d --link mysqlku -e SECRET="secr3t" -e SERVERPORT=8000 -e Name="costumer" -e Address=mysqlku -e Port=3306 -e Username="root" -e Password="yourmysqlpassword" programmerreski/costumer

# Run the app
$ docker logs costumer
```

## 📜 End Point  

Auth
| Methode       | End Point      | used for            
| ------------- | -------------  | -----------                  
| `POST`        | /api/v1/register            | Register
| `POST`        | /api/v1/login         | Login

Customers
| Methode       | End Point      | used for            
| ------------- | -------------  | -----------                  
| `GET`         | /api/v1/customers             | Get all customers      
| `GET`         | /api/v1/customers/:id          | Get One customers      
| `GET`         | /api/v1//customers/search       | Searching a customers      
| `POST`        | /api/v1/customers              | Insert customers 
| `PUT`         | /api/v1/customers/:id         | Update data customers
| `DELETE`      | /api/v1/customers/:id         | Delete customers  


Products
| Methode       | End Point      | used for            
| ------------- | -------------  | -----------                  
| `GET`         | /api/v1/products             | Get all products      
| `GET`         | /api/v1/products/:id          | Get One products      
| `GET`         | /api/v1//products/search          | Searching products      
| `POST`        | /api/v1/products              | Insert products 
| `PUT`         | /api/v1/products/:id         | Update data products
| `DELETE`      | /api/v1/products/:id         | Delete products  

Orders
| Methode       | End Point      | used for            
| ------------- | -------------  | -----------                  
| `GET`         | /api/v1/orders             | Get all orders      
| `GET`         | /api/v1/myorder             | Get only my order   
| `PUT`         | /api/v1/myorder/:id                 | Update my order   
| `GET`         | /api/v1/orders/:id          | Get One orders      
| `GET`         | /api/v1//orders/search          | Search an orders      
| `POST`        | /api/v1/orders              | Create orders 
| `POST`        | /api/v1/multiorder              | Create multi order / batch order
| `PUT`         | /api/v1/orders/:id         | Update data orders
| `DELETE`      | /api/v1/orders/:id         | Delete orders  

## 📜 Swagger Open Api
after you running the app you can access swagger open api in this 
 http://localhost:8080/swagger/index.html

## 📜 Postman 
you can find postman testing in  `/screenshoot/` folder

## 🛠️ Technology

This software uses the following Tech:

- [Golang](https://go.dev/dl/)
- [Gin Framework](https://gin-gonic.com/)
- [Gorm](https://gorm.io/index.html)
- [OpenAPI Swaggo](https://github.com/swaggo/gin-swagger)
- [UUID](github.com/google/uuid)
- [mysql](https://www.mysql.com/)
- [Linux](https://www.linux.com/)
- [Docker](https://www.docker.com/)
- [Dockerhub](https://hub.docker.com/u/programmerreski)
- [Git Repository](https://github.com/reski-id)
- [Trunk Base Development](https://trunkbaseddevelopment.com/)


## 📱 Contact me
feel free to contact me ... 
- Email programmer.reski@gmail.com 
- [Linkedin](https://www.linkedin.com/in/reski-id)
- [Github](https://github.com/reski-id)
- Whatsapp <a href="https://wa.me/+6281261478432?text=Hello">Send WhatsApp Message</a>
