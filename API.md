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
        "Name": "Goalkeeper Player",
        "Height": 190,
        "Weight": 85,
        "Position": "GOALKEEPER",
        "BackNumber": 1,
        "TeamID": 1
      },
      {
        "ID": 12,
        "Name": "Defender Player",
        "Height": 180,
        "Weight": 75,
        "Position": "DEFENDER",
        "BackNumber": 4,
        "TeamID": 1
      },
      {
        "ID": 16,
        "Name": "Midfielder Player",
        "Height": 175,
        "Weight": 70,
        "Position": "MIDFIELDER",
        "BackNumber": 6,
        "TeamID": 1
      },
      {
        "ID": 20,
        "Name": "Forward Player",
        "Height": 178,
        "Weight": 72,
        "Position": "FORWARD",
        "BackNumber": 9,
        "TeamID": 1
      }
    ]
  }
}
```
- Keterangan:
  - `Position` adalah tipe posisi pemain secara umum: GOALKEEPER (penjaga gawang), DEFENDER (pemain bertahan), MIDFIELDER (gelandang), FORWARD (penyerang)


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
    "Name": "John Doe",
    "Height": 180,
    "Weight": 70,
    "Position": "MIDFIELDER",
    "BackNumber": 6,
    "TeamID": 1
  }
}
```
- Keterangan:
  - `Position` adalah tipe posisi pemain: GOALKEEPER, DEFENDER, MIDFIELDER, atau FORWARD


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
            "PlayerName": "John Doe",
            "Position": "GOALKEEPER",
            "BackNumber": 1,
            "MinuteIn": 0,
            "MinuteOut": 0,
            "IsStarter": true,
            "SubstitutedForPlayerID": 0,
            "SubstitutedForPlayerName": "",
            "Reason": ""
          },
          {
            "PlayerID": 12,
            "PlayerName": "David Silva",
            "Position": "DEFENDER",
            "BackNumber": 2,
            "MinuteIn": 0,
            "MinuteOut": 0,
            "IsStarter": true,
            "SubstitutedForPlayerID": 0,
            "SubstitutedForPlayerName": "",
            "Reason": ""
          },
          {
            "PlayerID": 13,
            "PlayerName": "Marcus Rodriguez",
            "Position": "DEFENDER",
            "BackNumber": 3,
            "MinuteIn": 0,
            "MinuteOut": 0,
            "IsStarter": true,
            "SubstitutedForPlayerID": 0,
            "SubstitutedForPlayerName": "",
            "Reason": ""
          }
        ],
        "PlayerLineupPosition": [
          {
            "PlayerID": 11,
            "PositionSlot": "GK",
            "X": 0,
            "Y": 50
          },
          {
            "PlayerID": 12,
            "PositionSlot": "CB-L",
            "X": 20,
            "Y": 35
          },
          {
            "PlayerID": 13,
            "PositionSlot": "CB-R",
            "X": 20,
            "Y": 65
          }
        ]
      },
      "away": {
        "TeamID": 2,
        "Name": "Team B",
        "LogoURL": "https://example.com/logo2.png",
        "Formation": "4-4-2",
        "Player": [
          {
            "PlayerID": 21,
            "PlayerName": "Alex Johnson",
            "Position": "GOALKEEPER",
            "BackNumber": 1,
            "MinuteIn": 0,
            "MinuteOut": 0,
            "IsStarter": true,
            "SubstitutedForPlayerID": 0,
            "SubstitutedForPlayerName": "",
            "Reason": ""
          }
        ],
        "PlayerLineupPosition": [
          {
            "PlayerID": 21,
            "PositionSlot": "GK",
            "X": 100,
            "Y": 50
          }
        ]
      }
    },
    "Scored": {
      "home": 1,
      "away": 0
    },
    "Goal": [
      {
        "PlayerID": 13,
        "PlayerName": "Marcus Rodriguez",
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
- Keterangan struktur response:
  - `MatchID`: ID pertandingan
  - `Team.home` dan `Team.away`: Informasi tim home (di kiri) dan away (di kanan)
  - `Team[].Formation`: Formasi tim, contoh "4-3-3", "4-4-2"
  - `Team[].Player`: Daftar pemain di tim dengan detail berikut:
    - `Position`: Tipe posisi pemain secara umum (GOALKEEPER, DEFENDER, MIDFIELDER, FORWARD)
    - `PositionSlot` (di PlayerLineupPosition): Posisi spesifik di lapangan (GK, LB, CB-L, CB-R, RB, CM-L, CM-R, LW, RW, ST-L, ST-R)
    - Menit masuk/keluar, info substitusi, dan nomor punggung
  - `Team[].PlayerLineupPosition`: Posisi pemain di lapangan dengan koordinat X (0-100, horizontal) dan Y (0-100, vertikal) untuk visualisasi
    - X=0, Y=50 adalah GK home di gawang kiri
    - X=100, Y=50 adalah GK away di gawang kanan
  - `Scored`: Skor akhir atau sementara pertandingan
  - `Goal`: Daftar gol yang tercetak beserta menit tercetaknya
  - `Phase`: Fase pertandingan (1=active, 2=cancelled, 3=finished)
  - `Venue`: Informasi lokasi pertandingan



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
  "team_id": 1,
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
  "position_slot": "ST-L",
  "reason": "Ganti dulu yah"
}
```
- Response sukses contoh:
```json
{
  "message": "success to substitute player"
}
```

