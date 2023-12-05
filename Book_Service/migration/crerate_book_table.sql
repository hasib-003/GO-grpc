create table books(
                      id  serial,
                      title text not null,
                      author_id int not null,
                      publication_year int not null,
                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                      update_at TIMESTAMP default null,
                      deleted_at TIMESTAMP default null,
                      created_by uuid default null,
                      updated_by uuid default null
);