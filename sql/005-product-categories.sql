ALTER TABLE products ADD COLUMN IF NOT EXISTS category_id INTEGER REFERENCES categories(id);

UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'clothing')
WHERE code IN ('PROD001', 'PROD004', 'PROD007');

UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'shoes')
WHERE code IN ('PROD002', 'PROD006');

UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'accessories')
WHERE code IN ('PROD003', 'PROD005', 'PROD008');

ALTER TABLE products ALTER COLUMN category_id SET NOT NULL;