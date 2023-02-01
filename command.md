# Go build

- 명령어(for Mac) : CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gomovies ./cmd/api

# Postgress DB dump

- pg_dump를 설치
- 명령어 : pg_dump --no-owner -h DB주소(예: localhost) -p DB포트(예: 5432) -u 사용자명(예:user) DB명(예: movies) > 출력파일명(예: movies.sql)
