# Eventsam SQLite
Eventsam SQLite adalah versi dari Eventsam yang menggunakan SQLite sebagai engine database. Eventsam SQLite dapat digunakan untuk membuat aplikasi dengan paradigma event sourcing secara mudah.

# Instalasi Server
Untuk instalasi server, silahkan ikuti langkah-langkah berikut:
1. Clone repository ini
2. Masuk ke folder cmd/mysql `cd cmd/sqlite`
3. Jalankan perintah `go build -o ./eventsam ./server` di dalam directory
4. Jalankan perintah `./eventsam` untuk pengecekan instalasi
5. Copy file `eventsam` ke folder yang diinginkan
6. Tambahkan file `.env` di folder yang sama dengan file `eventsam` yang sudah di-copy, atau bisa juga menjalankan binary ini dengan menambahkan sendiri enviornment di command line
7. Jalankan `./eventsam` di folder yang telah di-copy
8. Jika sudah berjalan, maka Anda dapat menggunakan systemd atau supervisor untuk menjalankan server ini

## Environment
### Sebagai Server Master
Gunakan DB_FILEPATH dengan akhiran db. Nama file bebas. Contoh:
```
PORT=8009
DB_FILEPATH=/var/data/eventsam/event_primary.db
```
### Sebagai Server Slave
Sebagai server slave, tambahkan `MASTER_ADDRESS`, yaitu alamat eventsam master tanpa `/` Gunakan DB_FILEPATH dengan akhiran db. Nama file bebas. Contoh:
```
PORT=8010
MASTER_ADDRESS=http://localhost:8009
DB_FILEPATH=/var/data/eventsam/event_secondary.db
```
