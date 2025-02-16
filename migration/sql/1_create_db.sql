CREATE DATABASE merch_store
    WITH
    OWNER = merch_service_user
    ENCODING = 'UTF8'
    LOCALE_PROVIDER = 'libc'
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;

CREATE SCHEMA merch
    AUTHORIZATION merch_service_user;