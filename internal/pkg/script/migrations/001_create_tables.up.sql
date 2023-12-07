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

CREATE TABLE "posts" (
    "id" UUID PRIMARY KEY,
    "title" JSONB NOT NULL,
    "content" JSONB NOT NULL,
    "status" BOOLEAN DEFAULT TRUE,
    "pub_date" DATE,
    "author_id" UUID NOT NULL REFERENCES "users"("id"),
    "created_at" TIMESTAMP DEFAULT now() NOT NULL,
    "created_by" UUID NOT NULL REFERENCES "users"("id"),
    "updated_at" TIMESTAMP,
    "updated_by" UUID REFERENCES "users"("id"),
    "deleted_at" TIMESTAMP,
    "deleted_by" UUID REFERENCES "users"("id")
);
