CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS employees
(
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email           VARCHAR(50),
    password_user   VARCHAR(255),
    role_user       VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS pickup_points
(
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    registration_date    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    city                VARCHAR(50)
);


CREATE TABLE IF NOT EXISTS receptions
(
    id                      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    date_time               TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status_reception        VARCHAR(50),
    pickup_point_id UUID    REFERENCES pickup_points (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS products
(
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    date_time           TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    type_product        VARCHAR(50),
    reception_id UUID   REFERENCES receptions (id) ON DELETE CASCADE
);