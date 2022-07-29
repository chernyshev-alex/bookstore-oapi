
create table books (
	id          INTEGER primary key autoincrement,
    authorid    int not null,
	bookid      int not null,
	publisher   varchar(255) null,
	publisherid int not null,
	title       varchar(1024) null,
	year        int not null,
	descr       varchar(10000) null,
	isbn        varchar(32)     not null,
    created_at  timestamp null default current_timestamp,
    updated_at  timestamp null default current_timestamp
);

