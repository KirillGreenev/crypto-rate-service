CREATE TABLE ask(
                    id SERIAL PRIMARY KEY,
                    price DECIMAL(8, 2),
                    volume DECIMAL(8, 2),
                    amount DECIMAL(10, 2),
                    factor FLOAT,
                    type VARCHAR(20)
);

CREATE TABLE bid(
                    id SERIAL PRIMARY KEY,
                    price DECIMAL(8, 2),
                    volume DECIMAL(8, 2),
                    amount DECIMAL(10, 2),
                    factor FLOAT,
                    type VARCHAR(20)
);

CREATE TABLE rates(
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMP,
    ask_id INT REFERENCES ask(id) ON DELETE CASCADE,
    bid_id INT REFERENCES bid(id) ON DELETE CASCADE
);



