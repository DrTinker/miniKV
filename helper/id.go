package helper

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
)

func GenClientId() string {
	u := uuid.NewV4()
	str := u.String()
	id := md5.Sum([]byte(str))

	return fmt.Sprintf("%x", id)
}

func GenSeqId() string {
	u := uuid.NewV4()
	str := u.String() + strconv.Itoa(int(time.Now().Unix()))
	id := md5.Sum([]byte(str))

	return fmt.Sprintf("%x", id)
}
