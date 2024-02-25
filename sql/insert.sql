INSERT INTO queues (title, message, open) VALUES ('queue1', 'Welcome to queue1!', true);
INSERT INTO queues (title, message, open) VALUES ('queue2', 'Welcome to queue2!', false);

-- password = bcrypt of 'admin', 'pass' and 'word'
INSERT INTO users (username, password) VALUES ('admin', '$2a$10$wZV0NYmeN.2/AhnEJ7bz5Ohr1D0dayzakwwTewbIJSE7q8DqbK/Um');
INSERT INTO users (username, password) VALUES ('user1', '$2a$10$zI.vwmDSVUGzO3022zTIouPCV6dbHNq4rqCyG4mutTOmBfKjy/Zea');
INSERT INTO users (username, password) VALUES ('user2', '$2a$10$.p9jnBTMufPIh0OLk92B9.zb55BasX9lH/T5v9oBNE5Qki2CJvNzO');

INSERT INTO messages (fromid, toid, time, text) VALUES (2, 3, CURRENT_TIMESTAMP, 'hello!');
INSERT INTO messages (fromid, toid, time, text) VALUES (3, 2, CURRENT_TIMESTAMP, 'hi!');

INSERT INTO queueitems (userid, queueid, location, active, comment) VALUES (1, 1, 'chair', true, 'help!');
INSERT INTO queueitems (userid, queueid, location, active, comment) VALUES (2, 1, 'room', true, 'present!');

INSERT INTO admins (userid, queueid) values (1, 1);
INSERT INTO admins (userid, queueid) values (1, 2);

INSERT INTO accesses (userid, queueid) values (2, 1);
INSERT INTO accesses (userid, queueid) values (3, 1);
