-- 001_words.sql
-- Seed words + symmetric word-related pairs

WITH new_words AS (
    INSERT INTO words (word)
    VALUES
        ('coffee'),
        ('tea'),
        ('beach'),
        ('mountain'),
        ('dog'),
        ('cat'),
        ('sun'),
        ('moon'),
        ('car'),
        ('bicycle'),
        ('pizza'),
        ('salad'),
        ('phone'),
        ('laptop'),
        ('rain'),
        ('snow'),
        ('river'),
        ('ocean'),
        ('bread'),
        ('butter')
    ON CONFLICT DO NOTHING
    RETURNING id, word
)
INSERT INTO word_related (word_id_1, word_id_2)
SELECT LEAST(a.id, b.id), GREATEST(a.id, b.id)
FROM new_words a
JOIN new_words b ON a.word <> b.word
WHERE (a.word, b.word) IN (
    ('coffee', 'tea'),
    ('beach', 'mountain'),
    ('dog', 'cat'),
    ('sun', 'moon'),
    ('car', 'bicycle'),
    ('pizza', 'salad'),
    ('phone', 'laptop'),
    ('rain', 'snow'),
    ('river', 'ocean'),
    ('bread', 'butter')
)
ON CONFLICT DO NOTHING;
