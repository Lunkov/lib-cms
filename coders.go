package cms

import (
  "encoding/base64"
  "crypto/sha1"
)

func SHA1(bv string) string {
  hasher := sha1.New()
  hasher.Write([]byte(bv))
  return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
