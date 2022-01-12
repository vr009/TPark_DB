CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE users(
                      email citext UNIQUE NOT NULL,
                      fullname varchar NOT NULL,
                      nickname citext COLLATE "C" UNIQUE PRIMARY KEY,
                      about text NOT NULL DEFAULT ''
);

CREATE UNIQUE INDEX users_nickname ON users(nickname);
CREATE UNIQUE INDEX users_email ON users(email);
CREATE INDEX users_full ON users(email, nickname);

CREATE UNLOGGED TABLE forums (
    title varchar NOT NULL,
    author citext COLLATE "C" ,
    slug citext PRIMARY KEY,
    posts int DEFAULT 0,
    threads int DEFAULT 0
);
CREATE unique INDEX forums_slug ON forums(slug);

CREATE UNLOGGED TABLE forum_users (
    nickname citext  COLLATE "C",
    forum citext COLLATE "C" ,
    CONSTRAINT fk UNIQUE(nickname, forum)
);

CREATE INDEX fu_forum ON forum_users(forum);
CREATE INDEX fu_full ON forum_users(nickname,forum);

CREATE UNLOGGED TABLE threads (
    id serial PRIMARY KEY,
    author  citext COLLATE "C" ,
    message citext NOT NULL,
    title citext NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE,
    forum  citext COLLATE "C" ,
    slug citext,
    votes int
);

CREATE INDEX IF NOT EXISTS threads_slug ON threads(slug);
CREATE INDEX IF NOT EXISTS threads_id ON threads(id);
CREATE INDEX  IF NOT EXISTS cluster_thread ON threads(id, forum);
CREATE INDEX IF NOT EXISTS threads_forum ON threads(forum);
CREATE INDEX IF NOT EXISTS created_forum_index ON threads(forum, created_at);
CREATE INDEX ON threads(slug, id, forum);

CREATE UNLOGGED TABLE posts (
    id serial  PRIMARY KEY ,
    author citext COLLATE "C",
    post text NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE,
    forum citext COLLATE "C",
    isEdited bool,
    parent int,
    thread int,
    path integer []
);
CREATE INDEX IF NOT EXISTS posts_thread ON posts(thread);
CREATE INDEX IF NOT EXISTS posts_parent_thread_index ON posts(parent, thread);
CREATE INDEX IF NOT EXISTS  parent_tree_index ON posts ((path[1]), path, id);
CREATE unique INDEX IF NOT EXISTS posts_id ON posts(id);

CREATE UNLOGGED TABLE votes (
    author citext COLLATE "C",
    vote int,
    thread int,
    CONSTRAINT checks UNIQUE(author, thread)
);
CREATE INDEX votes_full ON votes(author, vote, thread);

CREATE OR REPLACE FUNCTION update_path() RETURNS TRIGGER AS
$update_path$
DECLARE
parent_path integer[];
    parent_thread int;
BEGIN
SELECT path FROM posts WHERE id = new.parent INTO parent_path;
NEW.path := parent_path || new.id;
RETURN new;
END
$update_path$ LANGUAGE plpgsql;


CREATE TRIGGER path_update_trigger
    BEFORE INSERT
    ON posts
    FOR EACH ROW
    EXECUTE PROCEDURE update_path();

VACUUM;
VACUUM ANALYSE;

