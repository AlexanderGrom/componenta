package hash

import (
    "crypto/sha256"
    "fmt"
)

type Hash struct {
}

func New() *Hash {
    return &Hash{}
}

func (self *Hash) Sum(b []byte) string {
    return fmt.Sprintf("%x", sha256.Sum256(b))
}
