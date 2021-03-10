Before i start, i'm so sorry if this readme fill crispy (read: garing) because i actually never make the readme so i don't know how to configure the font

--- What is this project about? ---
This is mini project about 3 microservices that can:
- Users Services:
  - create users with customer or admin roles
  - update users profile
  - update users password
- Products Services:
  - list all products
  - products detail
  - add and update products by admin roles
  - update products to inactive will find if there any orders still in pending, if there is still in pending, the product cannot be set to inactive instead you can make the qty 0 first
- Orders Services:
  - Create orders
  - Update orders by users
  - Cancel orders by users
  - List orders by users orders
  - List orders all users by admin roles
  - Approve and reject orders
  - note: all update, cancel, and reject orders will update the quantity products on products services

This project using clean architecture with microservices approach with monorepo structure
there is also migration script sql query when you run the docker-compose

need the postman with all endpoint and example response? here it is: https://www.getpostman.com/collections/0a75523cfba669054828

you can run all this project with docker-compose, just need to run `docker-composes up --build -d` (if you doesn't want to run background just remove the `-d`) the docker compose also will run the infrastructure like postgres and redis with default port on your container, so watch out, is your port is available?

what is the minus? i still not implement the unit test, why? because i'm still not familiar with golang testing (before when i still create unit test was on node js based project), i will continue to learn and maybe update this project, especially the unit test and maybe i will create the service diagram

that's all, if you wanna ask more just email me or make issue on this repo

Thank you, have a nice day