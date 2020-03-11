--CREATE USER webapp WITH PASSWORD '23698741';
--CREATE USER test WITH PASSWORD 'test';

--CREATE DATABASE webdb;
--\c webdb

--for test
CREATE SCHEMA IF NOT EXISTS AUTHORIZATION test;
CREATE TABLE IF NOT EXISTS test.account (
    user_id serial PRIMARY KEY,
    email VARCHAR (355) UNIQUE NOT NULL,
    password VARCHAR (50) NOT NULL,
    nickname VARCHAR (50) UNIQUE NOT NULL,
    created_on TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
ALTER TABLE test.account OWNER TO test;

--for webapp
CREATE SCHEMA IF NOT EXISTS AUTHORIZATION webapp;
CREATE TABLE IF NOT EXISTS webapp.account (
    user_id serial PRIMARY KEY,
    email VARCHAR (355) UNIQUE NOT NULL,
    password VARCHAR (50) NOT NULL,
    nickname VARCHAR (50) UNIQUE NOT NULL,
    created_on TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
ALTER TABLE webapp.account OWNER TO webapp;

