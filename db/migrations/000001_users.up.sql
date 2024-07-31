CREATE TABLE "users" (
  "id" varchar(50) PRIMARY KEY,
  "lastname" varchar(50) NOT NULL,
  "firstname" varchar(50) NOT NULL,
  "hashed_password" varchar NOT NULL,
  "phone" varchar(11) UNIQUE NOT NULL,
  "address" varchar(300) NOT NULL,
  "email" varchar(200) UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("email");