package data

const Schema = `
create table if not exists messages (
	id serial,
	version integer not null,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW(),
    message text not null
);
`
