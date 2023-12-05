create table user_registration(
                                  user_id serial,
                                  first_name text not null,
                                  last_name text not null,
                                  occupation text not null,
                                  role text not null,
                                  email text not null,
                                  password text not null,
                                  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                  update_at TIMESTAMP default null,
                                  deleted_at TIMESTAMP default null,
                                  created_by uuid default null,
                                  updated_by uuid default null,
                                  primary key (user_id)
);