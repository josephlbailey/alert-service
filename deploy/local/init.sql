revoke all on schema public from public;

-- create user
create user alert_service_owner with password 'alert_service_owner';

-- grant connect on database alert_service to alert_service_owner
grant connect, create on database alert_service to alert_service_owner;

-- grant usage and create privileges on the schema
grant usage, create on schema public to alert_service_owner;

-- grant all privileges on all tables in the schema
grant all on all tables in schema public to alert_service_owner;

-- to automatically grant all privileges on future tables in the schema
alter default privileges in schema public grant all on tables to alert_service_owner;

-- create alert_service_user user
create user alert_service_user with password 'alert_service_user';

-- grant connect on database alert_service to alert_service_user
grant connect on database alert_service to alert_service_user;

-- grant usage on the schema
grant usage on schema public to alert_service_user;

-- grant select, insert, update, delete privileges on all tables in the schema
grant select, insert, update, delete on all tables in schema public to alert_service_user;

-- automatically grant select, insert, update, delete privileges on future tables in the schema
alter default privileges for user alert_service_owner in schema public grant select, insert, update, delete on tables to alert_service_user;
