CREATE ROLE api_group INHERIT NOLOGIN;

CREATE DATABASE users_db;
CREATE DATABASE products_db;
CREATE DATABASE orders_db;

-- Users DB Migration
\c users_db;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL ON TABLES TO api_group ;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL ON SEQUENCES TO api_group ;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL ON FUNCTIONS TO api_group ;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL ON TYPES TO api_group ;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL ON SCHEMAS TO api_group ;

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email VARCHAR(255),
  phone_number VARCHAR(255),
  full_name VARCHAR(255),
  gender int,
  role  VARCHAR(255),
  password VARCHAR(255),
  created_at timestamp,
  updated_at timestamp
);

CREATE TABLE users_log (
  id SERIAL PRIMARY KEY,
  user_id int,
  email VARCHAR(255),
  phone_number VARCHAR(255),
  full_name VARCHAR(255),
  gender int,
  role  VARCHAR(255),
  password VARCHAR(255),
  created_at timestamp,
  updated_at timestamp
);

CREATE ROLE users_admin WITH ENCRYPTED PASSWORD 'password123' LOGIN;
GRANT api_group to users_admin;

-- Products DB Migration
\c products_db;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL ON TABLES TO api_group ;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL ON SEQUENCES TO api_group ;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL ON FUNCTIONS TO api_group ;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL ON TYPES TO api_group ;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL ON SCHEMAS TO api_group ;

CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255),
  price float,
  qty int,
  status int,
  created_at timestamp,
  updated_at timestamp
);

CREATE TABLE products_log (
  id SERIAL PRIMARY KEY,
  product_id int,
  user_id int,
  name VARCHAR(255),
  price float,
  qty int,
  status int,
  event VARCHAR(255),
  created_at timestamp,
  updated_at timestamp
);

INSERT INTO products ("name", "price", "qty", "status", "created_at", "updated_at") VALUES
('PRODUCT001', 15000, 15, 1, NOW(), NOW()),
('PRODUCT002', 18000, 18, 1, NOW(), NOW()),
('PRODUCT003', 7000, 2, 2, NOW(), NOW()),
('PRODUCT004', 5000, 0, 2, NOW(), NOW()),
('PRODUCT005', 35000, 10, 1, NOW(), NOW()),
('PRODUCT006', 30000, 18, 1, NOW(), NOW());


CREATE ROLE products_admin WITH ENCRYPTED PASSWORD 'password123' LOGIN;
GRANT api_group to products_admin;

-- Orders DB Migration
\c orders_db;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL ON TABLES TO api_group;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL ON SEQUENCES TO api_group;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL ON FUNCTIONS TO api_group;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL ON TYPES TO api_group;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL ON SCHEMAS TO api_group;

CREATE TABLE orders (
  id SERIAL PRIMARY KEY,
  user_id int,
  product_id int, 
  product_name VARCHAR(255),
  price float,
  qty int,
  total_price float,
  status int,
  created_at timestamp,
  updated_at timestamp
);

CREATE TABLE orders_log (
  id SERIAL PRIMARY KEY,
  user_id int,
  order_id int,
  product_id int, 
  product_name VARCHAR(255),
  price float,
  qty int,
  total_price float,
  status int,
  event VARCHAR(255),
  admin_id int,
  created_at timestamp,
  updated_at timestamp
);

CREATE ROLE orders_admin WITH ENCRYPTED PASSWORD 'password123' LOGIN;
GRANT api_group to orders_admin;
