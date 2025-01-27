-- +goose Up
CREATE TABLE IF NOT EXISTS tracks (
    id integer PRIMARY KEY,
    group_name text NOT NULL,
    song text  NOT NULL,
    text text  NOT NULL,
    realise_date text  NOT NULL,
    link text  NOT NULL
);

INSERT INTO tracks (group_name, song, text, realise_date, link) VALUES ('Survivor', 'Burning heart', 'bla-bla-bla', date('1985-04-12'), 'https://www.youtube.com/watch?v=Kc71KZG87X4');

-- +goose Down
DROP TABLE tracks;
