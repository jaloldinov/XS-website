CREATE TYPE "role_type" AS ENUM (
    'ADMIN',
    'AUTHOR'
);

CREATE TYPE "gender_type" AS ENUM (
    'M',
    'F'
);

CREATE TYPE "menu_type" AS ENUM (
  'MAIN',
  'EXTRA'
);

CREATE TYPE "file_type" AS ENUM (
  'IMAGE',
  'FILE'
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

CREATE TABLE "menu" (
  "id" uuid PRIMARY KEY,
  "title" jsonb NOT NULL,
  "content" jsonb NOT NULL,
  "parent_id" uuid,
  "is_static" bool DEFAULT false,
  "status" bool DEFAULT true,
  "slug" varchar NOT NULL,
  "type" menu_type NOT NULL,
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
    "slug" VARCHAR NOT NULL,
    "author_id" UUID NOT NULL REFERENCES "users"("id"),
    "menu_id" UUID NOT NULL REFERENCES "menu"("id"),
    "created_at" TIMESTAMP DEFAULT now() NOT NULL,
    "created_by" UUID NOT NULL REFERENCES "users"("id"),
    "updated_at" TIMESTAMP,
    "updated_by" UUID REFERENCES "users"("id"),
    "deleted_at" TIMESTAMP,
    "deleted_by" UUID REFERENCES "users"("id")
);

CREATE TABLE "menu_file" (
  "id" uuid PRIMARY KEY,
  "link" varchar NOT NULL,
  "type" file_type NOT NULL,
  "marked_link" varchar,
  "grouping" varchar,
  "carusel" bool DEFAULT false,
  "menu_id" uuid NOT NULL REFERENCES "menu"("id"),
  "author_id" UUID NOT NULL REFERENCES "users"("id"),
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "created_by" UUID NOT NULL REFERENCES "users"("id"),
  "updated_at" TIMESTAMP,
  "updated_by" UUID REFERENCES "users"("id"),
  "deleted_at" TIMESTAMP,
  "deleted_by" UUID REFERENCES "users"("id")
);

CREATE TABLE "post_file" (
  "id" uuid PRIMARY KEY,
  "link" varchar NOT NULL,
  "type" file_type NOT NULL,
  "marked_link" varchar,
  "grouping" varchar,
  "carusel" bool DEFAULT false,
  "post_id" uuid NOT NULL REFERENCES "posts"("id"),
  "author_id" UUID NOT NULL REFERENCES "users"("id"),
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "created_by" UUID NOT NULL REFERENCES "users"("id"),
  "updated_at" TIMESTAMP,
  "updated_by" UUID REFERENCES "users"("id"),
  "deleted_at" TIMESTAMP,
  "deleted_by" UUID REFERENCES "users"("id")
);
