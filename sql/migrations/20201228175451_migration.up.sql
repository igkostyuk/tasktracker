BEGIN;

CREATE TABLE IF NOT EXISTS projects (
  id UUID DEFAULT uuid_generate_v4(),

  name varchar,
  description varchar,

  PRIMARY KEY (id)
);

COMMIT;
