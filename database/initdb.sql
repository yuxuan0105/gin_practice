-- create a user 'webapp'
-- create a database 'webdb' 
-- import tables with commend psql -U postgres webdb -a -f initdb.sql

--for webapp
CREATE SCHEMA IF NOT EXISTS AUTHORIZATION webapp;
CREATE TABLE IF NOT EXISTS webapp.account (
    user_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    nickname VARCHAR UNIQUE NOT NULL,
    created_on TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
ALTER TABLE webapp.account OWNER TO webapp;

