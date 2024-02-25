CREATE TABLE "queues" (
  "id" int GENERATED ALWAYS AS IDENTITY,
  "title" varchar(64) NOT NULL,
  "message" varchar(128) NOT NULL,
  "open" boolean NOT NULL,
  PRIMARY KEY("id")
);

CREATE TABLE "users" (
  "id" int GENERATED ALWAYS AS IDENTITY,
  "username" varchar(32) UNIQUE NOT NULL,
  "password" varchar(60) NOT NULL,
  PRIMARY KEY("id")
);

CREATE TABLE "messages" (
  "id" int GENERATED ALWAYS AS IDENTITY,
  "fromid" int NOT NULL,
  "toid" int NOT NULL,
  "time" timestamp NOT NULL,
  "text" varchar(128) NOT NULL,
  PRIMARY KEY("id"),
  CONSTRAINT "FK_messages.from_id"
    FOREIGN KEY("fromid")
      REFERENCES "users"("id"),
  CONSTRAINT "FK_messages.to_id"
    FOREIGN KEY("toid")
      REFERENCES "users"("id")
);

CREATE TABLE "queueitems" (
  "id" int GENERATED ALWAYS AS IDENTITY,
  "userid" int NOT NULL,
  "queueid" int NOT NULL,
  "location" varchar(64) NOT NULL,
  "active" boolean NOT NULL,
  "comment" varchar(64) NOT NULL,
  PRIMARY KEY("id"),
  CONSTRAINT "FK_queueitems.user_id"
    FOREIGN KEY("userid")
      REFERENCES "users"("id"),
  CONSTRAINT "FK_admins.queue_id"
    FOREIGN KEY("queueid")
      REFERENCES "queues"("id")
);

CREATE TABLE "admins" (
  "userid" int NOT NULL,
  "queueid" int NOT NULL,
  PRIMARY KEY("userid", "queueid"),
  CONSTRAINT "FK_admins.user_id"
    FOREIGN KEY("userid")
      REFERENCES "users"("id"),
  CONSTRAINT "FK_admins.queue_id"
    FOREIGN KEY("queueid")
      REFERENCES "queues"("id")
);

CREATE TABLE "accesses" (
  "userid" int NOT NULL,
  "queueid" int NOT NULL,
  PRIMARY KEY("userid", "queueid"),
  CONSTRAINT "FK_admins.user_id"
    FOREIGN KEY("userid")
      REFERENCES "users"("id"),
  CONSTRAINT "FK_admins.queue_id"
    FOREIGN KEY("queueid")
      REFERENCES "queues"("id")
);
