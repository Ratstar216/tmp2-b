CREATE TABLE tweets (
    id char(26) NOT NULL PRIMARY KEY,
    user_id varchar(50) NOT NULL,
    parent_id char(26),
    username varchar(50) NOT NULL,
    created_at varchar(50) NOT NULL,
    likes int(5) NOT NULL,
    content varchar(300) NOT NULL
);

   -- retweet int(4) NOT NULL,