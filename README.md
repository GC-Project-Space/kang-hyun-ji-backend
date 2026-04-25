# 강현지 백엔드 (Ocean Diary API)

2026.04.25 Flutter Seoul Vibe Coding Hackathon

바다 일기 앱의 백엔드 서버입니다. 일기를 작성하면 현재 수심에 맞는 바다 생명체 알이 배정되고, 24시간 후 부화시켜 도감을 채워나가는 게임입니다.

---

## 기술 스택

- **Language**: Go
- **Framework**: Gin
- **Database**: SQLite (modernc.org/sqlite)
- **API Docs**: Swagger (swaggo/swag)

---

## 실행 방법

### 바이너리 직접 실행 (macOS)

```bash
./kang-hyun-ji-backend
```

### 소스에서 빌드 후 실행

```bash
go run .
```

### macOS 바이너리 빌드

```bash
GOOS=darwin GOARCH=arm64 go build -o kang-hyun-ji-backend .
```

서버는 **http://localhost:8080** 에서 실행됩니다.

---

## API 문서

서버 실행 후 Swagger UI에서 확인할 수 있습니다.

```
http://localhost:8080/swagger/index.html
```

---

## API 엔드포인트

### 유저

| Method | Path | 설명 |
|--------|------|------|
| GET | `/api/users/me` | 내 정보 조회 |
| PATCH | `/api/users/me/depth` | 수심 변경 (0~1000m) |

### 일기

| Method | Path | 설명 |
|--------|------|------|
| GET | `/api/diaries` | 일기 목록 조회 |
| POST | `/api/diaries` | 일기 작성 (알 생성) |
| GET | `/api/diaries/:id` | 일기 단건 조회 |
| POST | `/api/diaries/:id/hatch` | 알 부화 (24시간 후 가능) |

### 도감 / 수집

| Method | Path | 설명 |
|--------|------|------|
| GET | `/api/creatures` | 전체 생명체 목록 |
| GET | `/api/collection` | 내 수집 도감 |

### 업적

| Method | Path | 설명 |
|--------|------|------|
| GET | `/api/achievements` | 업적 목록 및 달성 여부 |

---

## 일기 작성 예시

```json
{
  "title": "오늘의 일기",
  "content": "오늘은 바다가 맑았다.",
  "diary_date": "2026-04-25",
  "category": "일상"
}
```

**category 가능값**: `감사` `일상` `감정` `목표` `기타`

---

## 생명체 희귀도

| 희귀도 | 설명 |
|--------|------|
| COMMON | 흔함 |
| UNCOMMON | 약간 희귀 |
| RARE | 희귀 |
| LEGENDARY | 전설 |

수심이 깊어질수록 희귀한 생명체가 출현합니다.

---

## 프로젝트 구조

```
.
├── main.go              # 서버 진입점, 라우터 설정
├── db/
│   ├── db.go            # DB 초기화 및 마이그레이션
│   └── seed.go          # 초기 데이터 시드
├── handlers/            # API 핸들러
│   ├── user.go
│   ├── diary.go
│   ├── collection.go
│   └── achievement.go
├── models/
│   └── models.go        # 데이터 모델
├── docs/                # Swagger 자동 생성 문서
└── assets/
    └── images/
        └── creatures/   # 생명체 이미지 파일
```

---

## DB 초기화

기존 데이터를 초기화하려면 DB 파일을 삭제 후 서버를 재시작합니다.

```bash
rm -f ocean_diary.db ocean_diary.db-shm ocean_diary.db-wal
./kang-hyun-ji-backend
```
