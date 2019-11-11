package slack

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Attain your signing secret by going to https://api.slack.com/apps, clicking on your app, and scrolling down to the
// section named App Credentials.
var signingSecret = []byte("fillmein")

// VerifyRequestSignature checks the signature of the request from Slack to guarantee authenticity
func VerifyRequestSignature(r *http.Request) bool {
	srt := r.Header.Get("X-Slack-Request-Timestamp")
	defer r.Body.Close()
	// Always v0 for now for some reason
	version := "v0"
	bod, _ := ioutil.ReadAll(r.Body)
	raw := fmt.Sprintf("%s:%s:%s", version, srt, bod)
	hm := hmac.New(sha256.New, signingSecret)
	hm.Write([]byte(raw))
	rawBytes := hm.Sum(nil)
	encoded := make([]byte, hex.EncodedLen(len(rawBytes)))
	hex.Encode(encoded, rawBytes)

	// Get the signature from the request for comparison
	signedSig := r.Header.Get("X-Slack-Signature")
	withV0 := []byte(fmt.Sprintf("%s=%s", version, string(encoded)))
	if valid := hmac.Equal(withV0, []byte(signedSig)); valid {
		// Rewind the body so that the other handlers can read off of it as they expect
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bod))
		return true
	}
	return false
}
