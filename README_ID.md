# Choose Language
[Indonesia](README_ID.md)

# Eventsam
Eventsam merupakan library untuk membuat aplikasi dengan paradigma event sourcing secara mudah. 

# Instalasi Server
Untuk instalasi server, silahkan ikuti langkah-langkah berikut:
1. Clone repository ini
2. Jalankan perintah `go build -o ./eventsam ./server`
3. Jalankan perintah `./eventsam`
4. Copy file `eventsam` ke folder yang diinginkan
5. Tambahkan file `.env` di folder yang sama dengan file `eventsam` yang sudah di-copy, atau bisa juga menjalankan binary ini dengan menambahkan sendiri enviornment di command line
6. Jalankan `./eventsam` di folder yang telah di-copy

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