-- LE VESTIAIRE — Football Categories

-- Remove old generic categories
DELETE FROM post_categories WHERE category_id IN (SELECT id FROM categories WHERE name IN ('General', 'Technology', 'Gaming', 'Music', 'Sport', 'Science', 'Art', 'Other', 'Tactique', 'Humour', 'Coup de gueule', 'Débat'));
DELETE FROM categories WHERE name IN ('General', 'Technology', 'Gaming', 'Music', 'Sport', 'Science', 'Art', 'Other', 'Tactique', 'Humour', 'Coup de gueule', 'Débat');

-- Football categories
INSERT OR IGNORE INTO categories (name) VALUES ('Actualité');
INSERT OR IGNORE INTO categories (name) VALUES ('Transferts');
INSERT OR IGNORE INTO categories (name) VALUES ('Résultats');
