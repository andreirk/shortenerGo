package link

import (
	"go/adv-demo/internal/stat"
	"gorm.io/gorm"
	"math/rand"
)

type Link struct {
	gorm.Model
	Url     string `json:"url"`
	Hash    string `json:"hash" gorm:"uniqueIndex"`
	hashLen int
	Stats   []stat.Stat `gorm:"constraint:OnUpdate:CASCADE,constraints:OnDelete:SET NULL;"`
}

func NewLink(url string, hashLen int) *Link {
	if hashLen == 0 {
		hashLen = 8
	}
	link := &Link{
		Url:     url,
		hashLen: hashLen,
	}
	link.GenerateHash()
	return link
}

func (link *Link) GenerateHash() {
	link.Hash = randStringRunes(link.hashLen)
}

var letterRunes = []rune("abcdefghiklmnoprstuvwxyzABCDEFGHIKLMOPRSTUVWZYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
