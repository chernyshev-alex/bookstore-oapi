
drop table if exists users;
create table users (
	id integer primary key autoincrement,
	name varchar(255)
);

	insert into users(name) values ('Olga');

drop table if exists books;

create table books (
	id          integer primary key autoincrement,
    authorid    integer 			not null,
	bookid      integer 			not	null,
	publisher   varchar(255) 		null,
	publisherid integer not 		null,
	title       varchar(1024) 		null,
	year        integer 			not null,
	descr       varchar(10000) 		null,
	isbn        varchar(32)     	not null,
    created_at  timestamp 			null default current_timestamp,
    updated_at  timestamp 			null default current_timestamp
);

insert into books(authorid, bookid, publisherid, year, title, isbn) 
	values (1000, 1, 100, 1990, 'title_1', 'isbn_1');

