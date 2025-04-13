CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE employees
(
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email           VARCHAR(50),
    password_user   VARCHAR(255),
    role_user       VARCHAR(50)
);

CREATE TABLE pickup_points
(
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    registration_date    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    city                VARCHAR(50)
);


CREATE TABLE receptions
(
    id                      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    date_time               TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status_reception        VARCHAR(50),
    pickup_point_id UUID    REFERENCES pickup_point (id) ON DELETE CASCADE
);

CREATE TABLE products
(
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    date_time           TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    type_product        VARCHAR(50),
    reception_id UUID   REFERENCES reception (id) ON DELETE CASCADE
);