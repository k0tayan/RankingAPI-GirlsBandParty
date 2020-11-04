package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"crypto/aes"
	"crypto/cipher"
	"os"
	"github.com/golang/protobuf/proto"
	"github.com/k0tayan/ranking_gbp/protos"
)

func UnpadByPKCS7(data []byte) []byte {
	padSize := int(data[len(data) - 1])
	return data[:len(data) - padSize]
}

func main(){
	decryptKey := []byte("")
	decryptIV := []byte("")
	host := "https://api.star.craftegg.jp"

	signature := ""
	userID := ""
	clientVer := ""
	masterDBVer := ""
	eventID := ""

	rankingURL := host + "/api/user/" + userID + "/event/" + eventID + "/versus/ranking"

	req, _ := http.NewRequest("GET", rankingURL, nil)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Accept", "application/octet-stream")
	req.Header.Set("Host", "api.star.craftegg.jp")
	req.Header.Set("X-Unity-Version", "2018.4.1f1")
	req.Header.Set("X-Signature", signature)
	req.Header.Set("X-ClientVersion", clientVer)
	req.Header.Set("X-ClientPlatform", "iOS")
	req.Header.Set("User-Agent", "band/3 CFNetwork/1121.2.2 Darwin/19.2.0")
	req.Header.Set("X-MasterDataVersion", masterDBVer)
	
	client := new(http.Client)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)

	c, err := aes.NewCipher(decryptKey)
	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s.", len(decryptKey), err)
		os.Exit(-1)
	}
	cbcdec := cipher.NewCBCDecrypter(c, decryptIV)
	plainText := make([]byte, len(byteArray))

	cbcdec.CryptBlocks(plainText, byteArray)
	plainText = UnpadByPKCS7(plainText)

	ranking := &UserEventRanking.UserVersusEventRankingResponse{}
	if err := proto.Unmarshal(plainText, ranking); err != nil {
		fmt.Println("Error: Failed to parse.", err)
	}
	for _, user := range ranking.EventPointTopUsers.GetEntries() {
		fmt.Println(user.GetRank(), user.GetRankLevel(), user.GetName(), user.GetPoint())
	}
}
