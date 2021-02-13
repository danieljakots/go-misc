// Copyright (c) 2021 Daniel Jakots

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package main

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"os"
)

func rightCert(myname string, names []string) bool {
	for _, name := range names {
		if myname == name {
			return true
		}
	}
	return false
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("USAGE: %s domain\n", os.Args[0])
		os.Exit(1)
	}
	peerName := os.Args[1]
	peer := fmt.Sprintf("%s:443", peerName)
	conf := &tls.Config{}
	conn, err := tls.Dial("tcp", peer, conf)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for _, cert := range conn.ConnectionState().PeerCertificates {
		if !rightCert(peerName, cert.DNSNames) {
			continue
		}
		fmt.Println("-----BEGIN CERTIFICATE-----")
		i := 1
		for _, char := range base64.StdEncoding.EncodeToString(cert.Raw) {
			fmt.Printf("%c", char)
			if i == 64 {
				i = 0
				fmt.Println()
			}
			i += 1
		}
		fmt.Println()
		fmt.Println("-----END CERTIFICATE-----")
	}
}
