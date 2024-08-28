-- +goose Up
ALTER TABLE stock
ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE pokemon
ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now(); 
   RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_stock_updated_at BEFORE UPDATE
ON stock FOR EACH ROW EXECUTE PROCEDURE 
update_updated_at_column();

CREATE TRIGGER update_pokemon_updated_at BEFORE UPDATE
ON pokemon FOR EACH ROW EXECUTE PROCEDURE 
update_updated_at_column();
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_pokemon_updated_at ON pokemon;

DROP TRIGGER update_stock_updated_at ON stock;

DROP FUNCTION update_updated_at_column()
-- +goose StatementEnd

ALTER TABLE pokemon
DROP COLUMN created_at,
DROP COLUMN updated_at;

ALTER TABLE stock
DROP COLUMN created_at,
DROP COLUMN updated_at;