CREATE TABLE "cats" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "years_of_experience" int NOT NULL,
  "breed" varchar NOT NULL,
  "salary" decimal NOT NULL,
  "status" varchar NOT NULL DEFAULT 'available',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "missions" (
  "id" bigserial PRIMARY KEY,
  "cat_id" bigint UNIQUE, -- A cat can only be on one mission at a time
  "completed" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "targets" (
  "id" bigserial PRIMARY KEY,
  "mission_id" bigint NOT NULL,
  "name" varchar NOT NULL,
  "country" varchar NOT NULL,
  "notes" text NOT NULL DEFAULT '',
  "completed" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "missions" ADD FOREIGN KEY ("cat_id") REFERENCES "cats" ("id") ON DELETE SET NULL;

ALTER TABLE "targets" ADD FOREIGN KEY ("mission_id") REFERENCES "missions" ("id") ON DELETE CASCADE;

CREATE INDEX ON "cats" ("status");
CREATE INDEX ON "missions" ("cat_id");
CREATE INDEX ON "targets" ("mission_id");
