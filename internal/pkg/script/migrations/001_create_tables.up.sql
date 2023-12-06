CREATE TYPE "role_type" AS ENUM (
    'ADMIN',
    'AUTHOR'
);

CREATE TYPE "gender_type" AS ENUM (
    'M',
    'F'
);

CREATE TABLE "users" (
      "id" uuid PRIMARY KEY,
      "avatar" varchar,
      "username" varchar UNIQUE NOT NULL,
      "password" varchar NOT NULL,
      "role" role_type NOT NULL,
      "full_name" varchar,
      "gender" gender_type,
      "birth_date" date,
      "status" bool DEFAULT true,
      "phone" varchar,
      "created_at" timestamp NOT NULL DEFAULT (now()),
      "created_by" uuid REFERENCES "users" ("id"),
      "updated_at" timestamp,
      "updated_by" uuid REFERENCES "users" ("id"),
      "deleted_at" timestamp,
      "deleted_by" uuid REFERENCES "users" ("id")
);
