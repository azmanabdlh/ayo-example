# API Documentation

## Overview
Semua endpoint REST di aplikasi ini diletakkan di prefix `/api`. Dokumentasi ini dibagi per domain berdasarkan struktur kode:
- Authentication
- Team Management
- League Management
- Health Check

Gunakan header berikut pada request yang memerlukan autentikasi:
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

Token diambil dari respon endpoint `POST /api/login`.

---

## Health Check

### `GET /api/ping`
- Deskripsi: Mengecek apakah service berjalan.
- Header: tidak diperlukan.
- Request payload: tidak ada.
- Response contoh:
```json
{
  "message": "Pong"
}
```

---

## Authentication

### `POST /api/login`
- Deskripsi: Login pengguna dan mendapatkan JWT Bearer token untuk endpoint yang diproteksi.
- Headers:
  - `Content-Type: application/json`
- Request payload:
```json
{
    "email": "ayo@example.com",
    "password": "password"
}
```
- Response sukses contoh:
```json
{
  "code": 200,
  "data": "<jwt-token-string>"
}
```
- Catatan: Gunakan nilai `data` sebagai token di header Authorization:
  `Authorization: Bearer <jwt-token-string>`.

---

## Team Management

### `GET /api/teams`
- Deskripsi: Ambil daftar tim.
- Header: tidak diperlukan.
- Query params:
  - `page`: halaman, default `1`.
  - `limit`: jumlah baris per halaman, default `10`, maksimum `100`.
- Response contoh:
```json
{
  "data": [
    {
      "ID": 1,
      "Name": "Team A",
      "LogoURL": "https://example.com/logo.png",
      "FoundedYear": 1990,
      "Address": "Stadium Street 1",
      "City": "Jakarta",
      "CreatedAt": "2026-06-18T00:00:00Z",
      "UpdatedAt": "2026-06-18T00:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10
  }
}
```

### `GET /api/teams/:id`
- Deskripsi: Ambil informasi detail tim berdasarkan ID.
- Header: tidak diperlukan.
- Response contoh:
```json
{
  "data": {
    "ID": 1,
    "Name": "Team A",
    "LogoURL": "https://example.com/logo.png",
    "FoundedYear": 1990,
    "Address": "Stadium Street 1",
    "City": "Jakarta",
    "Player": [
      {
        "ID": 11,
        "Name": "Player 1",
        "Height": 180,
        "Weight": 70,
        "Position": "GK",
        "BackNumber": 1,
        "TeamID": 1
      }
    ]
  }
}
```

### `POST /api/teams`
- Deskripsi: Buat tim baru.
- Headers:
  - `Content-Type: application/json`
  - `Authorization: Bearer <jwt-token>`
- Request payload:
```json
{
  "name": "Team A",
  "logo_url": "https://example.com/logo.png",
  "founded_year": 1990,
  "address": "Stadium Street 1",
  "city": "Jakarta"
}
```
- Response sukses contoh:
```json
{
  "code": 201,
  "message": "successfully add team"
}
```

### `PUT /api/teams/:id`
- Deskripsi: Ubah data tim berdasarkan ID.
- Headers:
  - `Content-Type: application/json`
  - `Authorization: Bearer <jwt-token>`
- Request payload:
```json
{
  "name": "Team A Updated",
  "logo_url": "https://example.com/logo-updated.png",
  "founded_year": 1990,
  "address": "Stadium Street 2",
  "city": "Jakarta"
}
```
- Response sukses contoh:
```json
{
  "code": 200,
  "message": "successfully modify team by id: 1"
}
```

### `DELETE /api/teams/:id`
- Deskripsi: Hapus tim berdasarkan ID.
- Headers:
  - `Authorization: Bearer <jwt-token>`
- Response sukses contoh:
```json
{
  "code": 200,
  "message": "successfully remove team by id: 1"
}
```

---

## Player Management

### `GET /api/players/:id`
- Deskripsi: Ambil detail pemain berdasarkan ID.
- Header: tidak diperlukan.
- Response contoh:
```json
{
  "data": {
    "ID": 11,
    "Name": "Player 1",
    "Height": 180,
    "Weight": 70,
    "Position": "MIDFIELDER",
    "BackNumber": 1,
    "TeamID": 1
  }
}
```

### `POST /api/players`
- Deskripsi: Assign pemain ke tim.
- Headers:
  - `Content-Type: application/json`
  - `Authorization: Bearer <jwt-token>`
- Request payload:
```json
{
  "player_id": 1,
  "team_id": 2
}
```
- Response sukses contoh:
```json
{
  "code": 201,
  "message": "successfully assign player:1  to team id: 2"
}
```

### `DELETE /api/players/:id`
- Deskripsi: Hapus pemain berdasarkan ID.
- Headers:
  - `Authorization: Bearer <jwt-token>`
- Response sukses contoh:
```json
{
  "code": 200,
  "message": "successfully remove player by id: 1"
}
```

---

## League Management

### `POST /api/matches`
- Deskripsi: Buat pertandingan baru.
- Headers:
  - `Content-Type: application/json`
  - `Authorization: Bearer <jwt-token>`
- Request payload:
```json
{
  "match_date": "2026-07-01T15:00:00Z",
  "title": "Final Match",
  "home_team_id": 1,
  "away_team_id": 2,
  "home_score": 0,
  "away_score": 0,
  "phase": 1,
  "venue_id": 1,
  "venue_name": "National Stadium"
}
```
- Response sukses contoh:
```json
{
  "message": "success to add match"
}
```

### `GET /api/matches/:id`
- Deskripsi: Ambil highlight pertandingan berdasarkan ID.
- Header: tidak diperlukan.
- Response sukses contoh:
```json
{
  "data": {
    "MatchID": 1,
    "Team": {
      "home": {
        "TeamID": 1,
        "Name": "Team A",
        "LogoURL": "https://example.com/logo.png",
        "Formation": "4-3-3",
        "Player": [
          {
            "PlayerID": 11,
            "PlayerName": "Player 1",
            "Position": "GK",
            "BackNumber": 1
          }
        ]
      },
      "away": {
        "TeamID": 2,
        "Name": "Team B",
        "LogoURL": "https://example.com/logo2.png",
        "Formation": "4-4-2",
        "Player": []
      }
    },
    "Scored": {
      "home": 1,
      "away": 0
    },
    "Goal": [
      {
        "PlayerID": 101,
        "PlayerName": "Striker",
        "ScoredAtMinute": 57
      }
    ],
    "Phase": 1,
    "Venue": {
      "ID": 1,
      "Name": "National Stadium",
      "Address": "Stadium Street 1",
      "City": "Jakarta",
      "GoogleMapsURL": "https://maps.example.com/stadium"
    }
  }
}
```

### `POST /api/matches/:id/lineup`
- Deskripsi: Assign pemain ke lineup pertandingan sebelum pertandingan dimulai. Hanya dapat dilakukan pada fase pertandingan yang upcoming (akan dimulai).
- Headers:
  - `Content-Type: application/json`
  - `Authorization: Bearer <jwt-token>`
- Request payload:
```json
{
  "team_id": 1,
  "lineup": [
    {
      "player_id": 11,
      "position_slot": "GK",
      "is_starter": true
    },
    {
      "player_id": 12,
      "position_slot": "CB-L",
      "is_starter": true
    },
    {
      "player_id": 13,
      "position_slot": "CB-R",
      "is_starter": true
    },
    {
      "player_id": 14,
      "position_slot": "RB",
      "is_starter": true
    },
    {
      "player_id": 15,
      "position_slot": "LB",
      "is_starter": true
    },
    {
      "player_id": 16,
      "position_slot": "CM-L",
      "is_starter": true
    },
    {
      "player_id": 17,
      "position_slot": "CM-R",
      "is_starter": true
    },
    {
      "player_id": 18,
      "position_slot": "LW",
      "is_starter": true
    },
    {
      "player_id": 19,
      "position_slot": "RW",
      "is_starter": true
    },
    {
      "player_id": 20,
      "position_slot": "ST-L",
      "is_starter": true
    },
    {
      "player_id": 21,
      "position_slot": "ST-R",
      "is_starter": false
    }
  ]
}
```
- Keterangan field:
  - `team_id`: ID tim yang akan di-assign lineupnya.
  - `lineup`: Array dari pemain yang akan masuk ke pertandingan.
    - `player_id`: ID pemain.
    - `position_slot`: Posisi pemain di lapangan. Valid values: `GK`, `LB`, `CB-L`, `CB-R`, `RB`, `CM-L`, `CM-R`, `LW`, `RW`, `ST-L`, `ST-R`.
    - `is_starter`: Boolean, `true` jika pemain start, `false` jika pemain cadangan.
- Response sukses contoh:
```json
{
  "message": "success to assign match player lineup"
}
```
- Catatan error:
  - Jika fase pertandingan bukan "upcoming", akan return error `"invalid match phase"`.
  - Jika tim tidak berpartisipasi di pertandingan, akan return error `"invalid teamID"`.
  - Jika ada pemain yang bukan member dari tim tersebut, akan return error `"invalid player team"`.

### `POST /api/matches/:id/goals`
- Deskripsi: Catat gol pada pertandingan.
- Headers:
  - `Content-Type: application/json`
  - `Authorization: Bearer <jwt-token>`
- Request payload:
```json
{
  "player_id": 10,
  "scored_at_minute": 57
}
```
- Response sukses contoh:
```json
{
  "message": "success add goal info to match_id: 1"
}
```

### P`OST /api/matches/:id/finish`
- Deskripsi: Tandai pertandingan sebagai selesai.
- Headers:
  - `Authorization: Bearer <jwt-token>`
- Request payload: tidak ada.
- Response sukses contoh:
```json
{
  "message": "success set match to finish"
}
```

### `POST /api/matches/:id/substitutions`
- Deskripsi: Lakukan substitusi pemain dalam pertandingan.
- Headers:
  - `Content-Type: application/json`
  - `Authorization: Bearer <jwt-token>`
- Request payload:
```json
{
  "team_id": 1,
  "minute": 60,
  "player_id": 12,
  "substituted_for_player_id": 11,
  "reason": "Injury"
}
```
- Response sukses contoh:
```json
{
  "message": "success to substitute player"
}
```

