package main

import (
	"io"
	//	"io/ioutil"
	"bufio"
	_ "expvar"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	connpool = []net.Conn{}
)

func conchandler(w http.ResponseWriter, r *http.Request, c *chan []byte) {
	log.Println("con")

	w.Header().Set("Content-Type", "audio/mpeg")
	for cc := range *c {
		log.Println("read something")
		nr := len(cc)
		log.Println(nr)

		nw, ew := w.Write(cc[0:nr])
		if ew == nil && nw != nr {
			ew = io.ErrShortWrite
		}
		var err error
		if ew != nil {
			err = ew
		}
		if err != nil {
			log.Println(err)
		}

	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	parentFolder := "./music"
	file, err := os.Open(parentFolder)

	if err != nil {
		log.Println(err)
	}
	//log.Println(file)
	files, err := file.Readdir(-1)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "audio/mpeg")

	playlist := []string{}
	for _, v := range files {
		curr := parentFolder + "/" + v.Name()
		playlist = append(playlist, curr)
	}

	//	i := 0
	for i := 0; i < len(playlist); i++ {
		f, err := os.Open(playlist[i])
		log.Println(playlist[i])
		if err != nil {
			log.Println(err)
		}

		buf := make([]byte, 100000)
		offset := int64(0)
		for {
			nr, err := f.ReadAt(buf, offset)
			if err != nil {
				log.Println(err)
			}
			if nr == 0 {
				//return 0, err
				log.Printf("read %d", nr)
				break
			}
			nw, ew := w.Write(buf[0:nr])
			if ew == nil && nw != nr {
				ew = io.ErrShortWrite
			}
			if ew != nil {
				err = ew
			}

			offset = int64(nr) + int64(offset)
			log.Println(offset)
		}
	} //	return int64(nw), err
}

func concplayer(c *[]*chan []byte) {
	//log.Println("dfgdfg handler")
	parentFolder := "./music"
	file, err := os.Open(parentFolder)

	if err != nil {
		log.Println(err)
	}
	//log.Println(file)
	files, err := file.Readdir(-1)
	if err != nil {
		log.Println(err)
	}

	playlist := []string{}
	for _, v := range files {
		curr := parentFolder + "/" + v.Name()
		playlist = append(playlist, curr)
	}

	//	i := 0
	for i := 0; i < len(playlist); i++ {
		f, err := os.Open(playlist[i])
		log.Println(playlist[i])
		if err != nil {
			log.Println(err)
		}

		buf := make([]byte, 128000)
		offset := int64(0)
		for {
			nr, err := f.ReadAt(buf, offset)
			if err != nil {
				log.Println(err)
			}
			if nr == 0 {
				//return 0, err
				log.Printf("read %d", nr)
				break
			}
			log.Println(len(*c))

			for i, v := range *c {
				log.Println(i)
				//log.Println(*c[i])

				*v <- buf[0:nr]

			}
			offset = int64(nr) + int64(offset)
			log.Println(offset)
			time.Sleep(2 * time.Second)
		}
	} //	return int64(nw), err

}

func player() {
	parentFolder := "./music"
	file, err := os.Open(parentFolder)

	if err != nil {
		log.Println(err)
	}
	//log.Println(file)
	files, err := file.Readdir(-1)
	if err != nil {
		log.Println(err)
	}

	playlist := []string{}
	for _, v := range files {
		curr := parentFolder + "/" + v.Name()
		playlist = append(playlist, curr)
	}

	//	i := 0
	for i := 0; i < len(playlist); i++ {
		f, err := os.Open(playlist[i])
		log.Println(playlist[i])
		if err != nil {
			log.Println(err)
		}

		buf := make([]byte, 128000)
		offset := int64(0)
		for {
			nr, err := f.ReadAt(buf, offset)
			if err != nil {
				log.Println(err)
			}
			if nr == 0 {
				//return 0, err
				log.Printf("read %d", nr)
				break
			}

			for _, c := range connpool {
				writer := bufio.NewWriter(c)

				//	writer.Header().Set("Content-Type", "audio/mpeg")

				nw, ew := writer.Write(buf[0:nr])
				if ew == nil && nw != nr {
					ew = io.ErrShortWrite
				}
				if ew != nil {
					err = ew
				}

			}

			offset = int64(nr) + int64(offset)
			log.Println(offset)
			time.Sleep(2 * time.Second)
		}
	} //	return int64(nw), err

}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("hello world")

	x := http.NewServeMux()

	//	chanpool := make([]chan []byte, 10)
	var chanpool []*chan []byte
	count := 0
	mutex := &sync.Mutex{}

	go concplayer(&chanpool)

	x.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		channel := make(chan []byte)
		chanpool = append(chanpool, &channel)

		//chanpool[count] = make(chan []byte)

		mutex.Lock()
		count = count + 1
		mutex.Unlock()
		log.Println(len(chanpool))
		log.Print(count)
		conchandler(w, r, chanpool[count-1])

	})

	err := http.ListenAndServe(":8080", x)
	if err != nil {
		log.Println(err)
	}
}
