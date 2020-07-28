package srsd

import (
	"fmt"
	"strings"
	"testing"
)

func TestStart(t *testing.T) {
	str := "/srsd/services/zacyuan.qq.com/"

	key := strings.Replace(str, "/srsd/services/", "", 1)
	fmt.Println(key)

	index := strings.LastIndex(key, "/")

	key2 := key[:index]
	fmt.Println(key2)

	key3 := key[index+1:]
	fmt.Println(key3)
}
