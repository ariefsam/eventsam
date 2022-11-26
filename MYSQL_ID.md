# Eventsam MySQL
Eventsam MySQL adalah versi dari Eventsam yang menggunakan MySQL sebagai engine database. Eventsam MySQL dapat digunakan untuk membuat aplikasi dengan paradigma event sourcing secara mudah.

# Instalasi Server
Untuk instalasi server, silahkan ikuti langkah-langkah berikut:
1. Clone repository ini
2. Masuk ke folder cmd/mysql `cd cmd/mysql`
3. Jalankan perintah `go build -o ./eventsam ./server` di dalam directory
4. Jalankan perintah `./eventsam` untuk pengecekan instalasi
5. Copy file `eventsam` ke folder yang diinginkan
6. Tambahkan file `.env` di folder yang sama dengan file `eventsam` yang sudah di-copy, atau bisa juga menjalankan binary ini dengan menambahkan sendiri enviornment di command line
7. Jalankan `./eventsam` di folder yang telah di-copy
8. Jika sudah berjalan, maka Anda dapat menggunakan systemd atau supervisor untuk menjalankan server ini

## Environment
### Sebagai Server Master
Gunakan beberapa environment database. Contoh:
```
PORT=8009

MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=password
MYSQL_DATABASE=eventsam
```


### Sebagai Server Slave
Sebagai server slave, tambahkan `MASTER_ADDRESS`, yaitu alamat eventsam master tanpa `/`. Contoh:
```
PORT=8009

MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=password
MYSQL_DATABASE=eventsam_secondary
```
Gunakan database yang berbeda dengan database master.


