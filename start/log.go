package start

import "github.com/sirupsen/logrus"

func InitLog() {
	logrus.SetLevel(logrus.InfoLevel)
}
