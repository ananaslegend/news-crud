create table if not exists posts (
  id serial primary key,
  title text not null,
  content text not null,
  author_id integer not null,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now()
);
