create extension if not exists pgcrypto;
create extension if not exists moddatetime;
create extension if not exists "uuid-ossp";

create table if not exists events(
  id uuid primary key default uuid_generate_v4(),

  sessionid uuid primary key default uuid_generate_v4(),
  visitorid uuid primary key default uuid_generate_v4(),

  etype varchar(24) not null,
  params jsonb not null default '{}'::jsonb,
  createdat timestamptz not null,

  cat timestamptz default now(),
  uat timestamptz default now()
);

create index e_sessionid on events (sessionid);
create index e_visitorid on events (visitorid);
create index e_etype on events (etype);
create index e_createdat on events (createdat);
create index e_params on events using gin (params);

drop trigger if exists uat_events on events;

create trigger uat_events
before update on events
for each row
  execute procedure moddatetime(uat);

