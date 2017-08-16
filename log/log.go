package log

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

// Config set log output
func Config(server string) {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// prefix := time.Now().Format("2006-01-02")
	// suffix := "-" + server + ".log"

	// file, err := os.OpenFile("/data/mj/"+prefix+suffix, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	// if err != nil {
	// 	log.WithFields(log.Fields{
	// 		"error": err,
	// 	}).Error(def.ErrLogToFile)
	// 	fmt.Println("Failed to log to file, using default stderr")
	// 	return
	// }
	// log.SetOutput(file)

	// // 清理文件
	// delLogFiles(suffix)
}

func delLogFiles(suffix string) {
	diffTime := int64(3600 * 24 * 30)
	nowTime := time.Now().Unix()
	filepath.Walk("/data/mj/", func(path string, file os.FileInfo, err error) error {
		if file == nil {
			return err
		}

		if !strings.HasSuffix(file.Name(), suffix) {
			return nil
		}

		modTime := file.ModTime().Unix()
		if (nowTime - modTime) > diffTime {
			os.Remove("/data/mj/" + file.Name())
		}
		return nil
	})
}
