package main

import (
	"crypto/rand"
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"
)

const oneMB = 1024 * 1024 // 1 MB

var (
	FlagAddress = flag.String("address", ":8080", "Address to listen on")
)

func main() {
	flag.Parse()

	http.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		sizeMBString := r.URL.Query().Get("size_mb")
		sizeMB, err := strconv.Atoi(sizeMBString)
		if err != nil {
			http.Error(w, "invalid size_mb", http.StatusBadRequest)
			return
		}
		MBpsString := r.URL.Query().Get("mbs")
		MBps, err := strconv.Atoi(MBpsString)
		if err != nil {
			http.Error(w, "invalid mbs", http.StatusBadRequest)
			return
		}
		log.Printf("serving %d MB file...\n", sizeMB)
		randomBytes := make([]byte, sizeMB*oneMB)
		_, err = rand.Read(randomBytes)
		if err != nil {
			http.Error(w, "error generating random data", http.StatusInternalServerError)
			return
		}

		tick := time.Tick(10 * time.Millisecond)
		bytesWritten := 0
		for {
			<-tick
			if bytesWritten >= sizeMB*oneMB {
				log.Printf("finished serving %d MB file successfully.\n", sizeMB)
				return
			}
			n := MBps * oneMB / 100
			if n > sizeMB*oneMB-bytesWritten {
				n = sizeMB*oneMB - bytesWritten
			}
			_, err := w.Write(randomBytes[bytesWritten : bytesWritten+n])
			if err != nil {
				log.Println(err)
				return
			}
			bytesWritten += n
		}
	})

	if FlagAddress == nil {
		panic("flagAddress is nil")
	}

	log.Printf("listening on %s...\n", *FlagAddress)
	err := http.ListenAndServe(*FlagAddress, nil)
	if err != nil {
		panic(err)
	}
}
