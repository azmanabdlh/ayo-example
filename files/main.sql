DO $$
DECLARE
    r RECORD;
BEGIN
    -- Mencari semua tabel buatan user di skema 'public'
    FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
        -- Mengeksekusi TRUNCATE + RESET ID + CASCADE untuk setiap tabel
        EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || ' RESTART IDENTITY CASCADE';
    END LOOP;
END $$;


-- ======================================
-- VENUES
-- ======================================
INSERT INTO venues (
    name,
    address,
    city,
    capacity,
    created_at,
    updated_at
)
VALUES
(
    'Stadion Gelora Bandung Lautan Api',
    'Jl. Gerbang Biru',
    'Bandung',
    38000,
    NOW(),
    NOW()
),
(
    'Jakarta International Stadium',
    'Jl. RE Martadinata',
    'Jakarta',
    82000,
    NOW(),
    NOW()
),
(
    'Stadion Kanjuruhan',
    'Jl. Trunojoyo',
    'Malang',
    42000,
    NOW(),
    NOW()
),
(
    'Stadion Gelora Bung Tomo',
    'Jl. Jawar',
    'Surabaya',
    46000,
    NOW(),
    NOW()
);

-- ======================================
-- USERS
-- password = admin123
-- ======================================
INSERT INTO users (email, password_hash)
VALUES
(
    'ayo@example.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy'
);

-- ======================================
-- TEAMS
-- ======================================
INSERT INTO teams (
    name,
    logo_url,
    founded_year,
    address,
    city
)
VALUES
('Persib Bandung','ayo.co.id/example.png', 1933,'Jl. Sulanjana','Bandung'),
('Persija Jakarta','ayo.co.id/example.png', 1928,'Jl. Stadion','Jakarta'),
('Arema FC','ayo.co.id/example.png', 1987,'Jl. Kertanegara','Malang'),
('Persebaya Surabaya','ayo.co.id/example.png', 1927,'Jl. Ahmad Yani','Surabaya');

-- ======================================
-- PLAYERS
-- ======================================
INSERT INTO players (
    name,
    height,
    weight,
    position,
    back_number
)
VALUES
('Ciro Alves',172,68,'FORWARD',77),
('David Da Silva',183,82,'FORWARD',19),
('Marc Klok',177,73,'MIDFIELDER',23),

('Riko Simanjuntak',168,62,'FORWARD',25),
('Ondrej Kudela',182,80,'DEFENDER',17),
('Hanif Sjahbandi',180,75,'MIDFIELDER',19),

('Charles Lokolingoy',178,75,'FORWARD',9),
('Sergio Silva',183,82,'DEFENDER',4),
('Arkhan Fikri',170,65,'MIDFIELDER',8),

('Bruno Moreira',177,73,'FORWARD',10),
('Dusan Stevanovic',186,85,'DEFENDER',5),
('Song Ui-young',174,69,'MIDFIELDER',7);

-- ======================================
-- PLAYER MEMBERSHIP
-- ======================================
INSERT INTO player_members (
    player_id,
    team_id,
    joined_at
)
VALUES
(1,1,NOW()),
(2,1,NOW()),
(3,1,NOW()),

(4,2,NOW()),
(5,2,NOW()),
(6,2,NOW()),

(7,3,NOW()),
(8,3,NOW()),
(9,3,NOW()),

(10,4,NOW()),
(11,4,NOW()),
(12,4,NOW());

-- ======================================
-- MATCHES
-- phase
-- 1 active
-- 2 cancelled
-- 3 finished
-- ======================================
INSERT INTO matches (
    title,
    match_date,
    home_team_id,
    away_team_id,
    home_score,
    away_score,
    phase,
    venue_id
)
VALUES
(
    'Persib Bandung vs Persija Jakarta',
    '2026-06-18 19:00:00',
    1,
    2,
    2,
    1,
    3,
    1
),
(
    'Arema FC vs Persebaya Surabaya',
    '2026-06-20 20:00:00',
    3,
    4,
    0,
    0,
    1,
    2
);

-- ======================================
-- GOALS
-- ======================================
INSERT INTO goals (
    match_id,
    player_id,
    team_id,
    scored_at_minute
)
VALUES
(1, 1, 1, 13),
(1, 2, 1, 52),
(1, 4, 2, 77);

-- ======================================
-- MATCH PLAYERS
-- ======================================
INSERT INTO match_players (
    match_id,
    team_id,
    player_id,
    position_slot,
    is_starter,
    minute_in,
    minute_out,
    substituted_for_player_id,
    reason,
    created_at,
    updated_at
)
VALUES
(1,1,1,'ST-L',FALSE,0,72,NULL,'TACTICAL',NOW(),NOW()),
(1,1,2,'ST-R',TRUE,0,0,NULL,NULL,NOW(),NOW()),
(1,1,3,'CM-L',TRUE,0,0,NULL,NULL,NOW(),NOW()),
(
    1,
    1,
    4,
    'ST-L',
    TRUE,
    72,
    0,
    1,
    'TACTICAL',
    NOW(),
    NOW()
);