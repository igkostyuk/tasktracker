BEGIN;

CREATE TABLE IF NOT EXISTS projects (
  id UUID DEFAULT uuid_generate_v4(),
  name varchar(500),
  description varchar(1000),

  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS columns (
  id UUID DEFAULT uuid_generate_v4(),
  position SERIAL,
  name varchar(255),
  status varchar(255),
  project_id UUID NOT NULL,
  FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS tasks (
  id UUID DEFAULT uuid_generate_v4(),
  position SERIAL,
  name varchar(500),
  description varchar(5000),
  colum_id UUID NOT NULL,
  FOREIGN KEY (colum_id) REFERENCES columns(id) ON DELETE CASCADE,

  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS comments (
  id UUID DEFAULT uuid_generate_v4(),
  text varchar(5000),
  task_id UUID NOT NULL,
  FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,

  PRIMARY KEY (id)
);

COMMIT;
