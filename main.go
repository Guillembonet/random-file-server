package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

const oneMB = 1024 * 1024 // 1 MB

var (
	FlagAddress   = flag.String("address", ":8080", "Address to listen on")
	FlagMaxSizeMB = flag.Int("maxSizeMB", 100000, "Maximum size of file to serve in MB")
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
		if sizeMB > *FlagMaxSizeMB {
			http.Error(w, fmt.Sprintf("size_mb must be less than or equal to %d", *FlagMaxSizeMB), http.StatusBadRequest)
			return
		}
		MBsString := r.URL.Query().Get("mbs")
		MBs := 0
		limitedSpeed := false
		if MBsString != "" {
			limitedSpeed = true
			MBs, err = strconv.Atoi(MBsString)
			if err != nil {
				http.Error(w, "invalid mbs", http.StatusBadRequest)
				return
			}
		}
		filename := r.URL.Query().Get("filename")
		if filename == "" {
			filename = "file.bin"
		}

		log.Printf("serving %d MB file...\n", sizeMB)

		w.Header().Set("Content-Length", fmt.Sprint(sizeMB*oneMB))
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

		randomBytes := make([]byte, sizeMB*oneMB)
		_, err = rand.Read(randomBytes)
		if err != nil {
			http.Error(w, "error generating random data", http.StatusInternalServerError)
			return
		}

		if !limitedSpeed {
			_, err := w.Write(randomBytes)
			if err != nil {
				log.Println(err)
			}
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
			n := MBs * oneMB / 100
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

	if FlagMaxSizeMB == nil {
		panic("flagMaxSizeMB is nil")
	}

	log.Printf("listening on %s...\n", *FlagAddress)
	err := http.ListenAndServe(*FlagAddress, nil)
	if err != nil {
		panic(err)
	}
}
