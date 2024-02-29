package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micmonay/keybd_event"
	"github.com/skip2/go-qrcode"
)

func main() {

	r := gin.Default()
	cmd := exec.Command("cmd", "/c", "start", "http://localhost:8080/powerclick/app")

	if err := cmd.Run(); err != nil {
		panic(err)
	}
	r.GET("/data/left", left)
	r.GET("/data/right", right)
	r.GET("/test", test)
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			err := qrcode.WriteFile(ipv4.String(), qrcode.High, 500, "./static/qr.png")
			if err != nil {
				fmt.Println("QR kodu dosyaya yazılamadı:", err)
				return
			}
			break
		}
	}

	r.LoadHTMLGlob("templates/*")

	r.GET("/powerclick/app", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"htmx": "../static/qr.png",
		})
	})
	r.Static("/static", "./static")

	r.Run(":8080")

}

func test(c *gin.Context) {
	c.JSON(200, gin.H{
		"test ": "basarılı",
	})
}

func left(c *gin.Context) {
	c.JSON(200, gin.H{
		"state": "left",
	})
	automate("0")
}
func right(c *gin.Context) {
	c.JSON(200, gin.H{
		"state": "right",
	})
	automate("1")
}

func automate(direction string) error {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return err
	}
	defer kb.Clear()

	if direction == "1" {
		kb.SetKeys(keybd_event.VK_RIGHT)
	} else if direction == "0" {
		kb.SetKeys(keybd_event.VK_LEFT)
	}

	err = kb.Launching()
	if err != nil {
		return err
	}

	time.Sleep(100 * time.Millisecond)
	return nil

}
