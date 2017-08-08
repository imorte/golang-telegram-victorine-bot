CREATE TABLE available
(
    id INTEGER PRIMARY KEY,
    group_telega TEXT,
    flag TEXT,
    current TEXT
);

CREATE TABLE pidors
(
    id INTEGER PRIMARY KEY NOT NULL,
    pidor TEXT,
    wich_group TEXT,
    score TEXT
);