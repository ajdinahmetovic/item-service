CREATE TABLE item (
	id bigserial PRIMARY KEY,
	title VARCHAR(255),
	description VARCHAR(255),
	user_id integer,
	FOREIGN KEY (user_id) REFERENCES app_user (id)

);
