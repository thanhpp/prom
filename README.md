# prom
- Project Management. HUST_20202_DS
- [DOCUMENT](https://docs.google.com/document/d/1P3y1ooL-uSGomqsTag6MspTsKyCxqANajjhT1v88H7c/edit?usp=sharing)
# Project folders

1. Golang
   1. cmd: for specified apps
   2. pkg: shared packages  
2. web: web resources
3. docs: documentation
4. build: build docker images, docker-compose
5. api: api document (move to embedded swagger) 

## DB desgin
### 1. Overall DB
![DBDesign](docs/dbdiagram/prom_db.png)

### 2. Actual DB

![ActualDB](docs/dbdiagram/prom_db_microservice.png)

## System design

![SystemDesign](docs/drawio/prom_design.png)

## DEVELOPMENT

1. Postman Collection: https://www.getpostman.com/collections/654a270d5382abd97404

2. SwaggerUI: http://127.0.0.1:12345/docs/index.html

3. Reorder cards
   1. Reorder in 1 column
      - Specify card id(1)
      - index of card(2) that you want to set card(1) above of it.
  
   2. Reorder card between column
      - Specify card id(1)
      - Specify column id that you want to move card(1) to.
      - index of card(2) that you want to set card(1) above of it.

   3. Example
      - ![ReorderCardExample](docs/drawio/ReorderCardExp.png)

4. Reorder columns

- ![ReorderColumnExample](docs/drawio/ReorderColumnExp.jpg)
