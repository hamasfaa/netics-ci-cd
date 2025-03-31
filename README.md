| Nama                      |    NRP     |
| ------------------------- | :--------: |
| Hamasah Fatiy Dakhilullah | 5025231139 |

# Penugasan Modul 1 Open Recruitment NETICS 2025

## Tujuan

Membuat sebuah API publik dengan endpoint `/health` dengan menerapkan konsep CI/CD menggunakan GitHub Actions.

## Link

- Published Docker Image <br/>
  [Docker Image](https://hub.docker.com/r/hamasfa/health/tags)

- Deployed API <br/>
  [URL API](http://128.199.145.27/health)

## Struktur Direktori

```
.
│   Dockerfile
│   go.mod
│   go.sum
│   main.go
│   Readme.md
│
├───.github
│   └───workflows
│           docker-image.yml
│
├───handler
│       health.go
│       index.go
│
└───response
        health_response.go
```

## Penjelasan Kode

- .github/workflows/docker-image.yml

```
name: Docker Image CI

on:
  push:
    branches: ["main"]
  pull_request:
    branches: ["main"]

jobs:
  build-and-test:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout the repository
        uses: actions/checkout@v4

      - name: Cache
        uses: actions/cache@v4
        with:
          path: |
            ~/.cache/go-build
            ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          restore-keys: |
            ${{ runner.os }}-go-

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: "1.23.2"

      - name: Download modules
        run: go mod download

      - name: Build Go
        run: go build -o api main.go

      - name: Test Go
        run: go test -v ./...

  docker:
    runs-on: ubuntu-latest
    needs: build-and-test
    steps:
      - name: Set up QEMU
        uses: docker/setup-qemu-action@v3

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          push: true
          tags: ${{ secrets.DOCKERHUB_USERNAME }}/health:latest
          cache-from: type=gha
          cache-to: type=gha,mode=max

  deploy:
    runs-on: ubuntu-latest
    needs: docker
    steps:
      - name: Deploy
        uses: appleboy/ssh-action@master
        with:
          host: ${{ secrets.HOST }}
          username: ${{ secrets.USERNAME }}
          key: ${{ secrets.SSH_KEY }}
          script: |
            docker pull ${{ secrets.DOCKERHUB_USERNAME }}/health:latest
            docker stop health || echo "Gak ada yang bisa distop"
            docker rm health || echo "Gak ada yang bisa dihapus"
            docker run -d --name health -p 80:8080 ${{ secrets.DOCKERHUB_USERNAME }}/health:latest
            docker image prune -f
```

Pada workflow ini terbagi menjadi 3 job, yaitu :

1. build-and-test <br/>
   Job ini melakukan tugas sebagai berikut :
   - Melakukan checkout untuk mengambil kode dari repositori
   - Menyimpan module-module yang digunakan untuk mempercepat build selanjutnya
   - Menyiapkan env Go dengan versi 1.23.2
   - Mendownload modules yang dibutuhkan
   - Melakukan compile dan melakukan testing
2. docker <br/>
   Job ini melakukan tugas sebagai berikut :
   - Menyiapkan Qemu dan Docker Buildx untuk build multi-arsitektur
   - Melakukan login ke docker hub
   - Melakukan build dan push ke docker hub
   - Menggunakan caching dari GitHub Actions untuk mempercepat build
3. deploy <br/>
   Job ini melakukan tugas sebagai berikut :
   - Melakukkan ssh ke vps
   - Menjalankan script dengan tujuan :
     1. Mengambil docker images health terbaru dari sumber yang telah diatur
     2. Memberhentikan container lama health, jika ada
     3. Menghapus container lama health, jika ada
     4. Menjalankan container baru dengan port 80:8080
     5. Menghapus image lama yang tidak terpakai untuk menghemat penyimpanan

- Dockerfile

```
FROM golang:1.23.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

Terdapat dua fase, yaitu :

1. Fase 1 (Builder) <br/>
   Fase ini melakukan tugas sebagai berikut :
   - Menggunakan image golang:1.23.2-alpine sebagai env
   - Menyalin go.mod dan go.sum, lalu mendownload modules yang diperlukan
   - Menyalin semua yang ada di dalam workspace dan melakukan compile
2. Fase 2 (Runtime) <br/>
   Fase ini melakukan tugas sebagai berikut :
   - Menggunakan image alpine:latest sebagai base image agar memiliki ukuran yang kecil
   - Menyalin hasil compile dari fase builder
   - Membuka port 8080
   - Menjalankan server

- response/health_response.go

```
package response

type HealthResponse struct {
	Nama      string `json:"nama"`
	NRP       string `json:"nrp"`
	Status    string `json:"status"`
	TimeStamp string `json:"timestamp"`
	Uptime    string `json:"uptime"`
}
```

Kode di atas adalah sebuah struktur data yang nantinya akan dikirim sebagai respons JSON.

- handler/index.go

```
package handler

import (
	"net/http"
)

func IndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, Docker Tes! <3"))
	}
}
```

Kode di atas adalah handler untuk endpoint `/` yang hanya menerima metode GET requests. Ketika mendapat request yang valid maka akan mengembalikan respons `Hello, Docker Tes! <3`

- handler/health.go

```
package handler

import (
	"encoding/json"
	"learn-ci-cd/response"
	"net/http"
	"time"
)

func HealthHandler(timeUp string, timeZone *time.Location) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
			return
		}

		currentTime := time.Now().In(timeZone).Format("2006-01-02 15:04:05")

		response := response.HealthResponse{
			Nama:      "Tunas Bimatara Chrisnanta Budiman",
			NRP:       "5025231999",
			Status:    "UP",
			TimeStamp: currentTime,
			Uptime:    timeUp,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
```

Kode di atas adalah handler untuk endpoint `/health` yang hanya menerima metode GET requests. Ketika mendapat request yang valid maka akan mengembalikan respons yang diambil dari `response/health_response.go`
