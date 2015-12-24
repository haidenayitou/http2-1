package hpack

import (
	"encoding/hex"
	"reflect"
	"strings"
	"testing"
)

func TestStaticTableDefinition(t *testing.T) {
	a.testEncode(t, 0)
	a.testDecode(t, 0)
}

func TestHeaderFieldRepresentation(t *testing.T) {
	c2[0].testEncode(t, 4096)
	c2[0].testDecode(t, 4096)
	c2[1].testEncode(t, 4096)
	c2[1].testDecode(t, 4096)
	c2[2].testEncode(t, 4096)
	c2[2].testDecode(t, 4096)
	c2[3].testEncode(t, 4096)
	c2[3].testDecode(t, 4096)
}

func TestRequestWithoutHuffmanCoding(t *testing.T) {
	c3.testEncode(t, 4096)
	c3.testDecode(t, 4096)
}

func TestRequestWithHuffmanCoding(t *testing.T) {
	c4.testEncode(t, 4096)
	c4.testDecode(t, 4096)
}

func TestResponseWithoutHuffmanCoding(t *testing.T) {
	c5.testEncode(t, 256)
	c5.testDecode(t, 256)
}

func TestResponseWithHuffmanCoding(t *testing.T) {
	c6.testEncode(t, 256)
	c6.testDecode(t, 256)
}

type testcase []struct {
	enc       string
	huff      huffman
	idx       indexType
	headers   []headerField
	table     []headerField
	tableSize uint32
}

func (tc testcase) testEncode(t *testing.T, tableSize uint32) {
	enc, buf := NewEncoder(tableSize), []byte{}

	for i, test := range tc {
		enc.h = test.huff
		enc.i = test.idx == indexedTrue

		exp := strings.Replace(test.enc, " ", "", -1)
		exp = strings.Replace(exp, "\n", "", -1)

		buf = buf[:0]

		for _, hf := range test.headers {
			_, buf = enc.EncodeHeaderField(buf, hf.name, hf.value, test.idx == indexedNever)
		}
		if got := hex.EncodeToString(buf); got != exp {
			t.Fatalf("encode test case %d: expected hex %v; got %v", i, exp, got)
		}

		table := make([]headerField, len(enc.table.data))
		for i := range enc.table.data {
			table[i] = enc.table.data[len(enc.table.data)-1-i]
		}
		if !reflect.DeepEqual(table, test.table) {
			t.Errorf("encode test case %d: expected table %v; got %v", i, test.table, table)
		}

		if enc.table.size != test.tableSize {
			t.Errorf("encode test case %d: expected table size %v; got %v", i, test.tableSize, enc.table.size)
		}
	}
}

func (tc testcase) testDecode(t *testing.T, tableSize uint32) {
	dec := NewDecoder(tableSize)
	for i, test := range tc {
		s := strings.Replace(test.enc, " ", "", -1)
		s = strings.Replace(s, "\n", "", -1)
		enc, err := hex.DecodeString(s)
		if err != nil {
			panic(err)
		}

		headers := []headerField{}
		_, err = dec.Decode(enc, 0, func(name, value string, sensitive bool) error {
			if sensitive != (test.idx == indexedNever) {
				t.Fatalf("decode test case %d: expected sensitivity %v; got %v", i, test.idx == indexedNever, sensitive)
			}
			headers = append(headers, headerField{name, value})
			return nil
		})
		if err != nil {
			t.Fatalf("decode test case %d: got '%v' error", i, err)
		}

		if err = dec.Reset(); err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(headers, test.headers) {
			t.Fatalf("decode test case %d: expected headers %v; got %v", i, test.headers, headers)
		}

		table := make([]headerField, len(dec.table.data))
		for i := range dec.table.data {
			table[i] = dec.table.data[len(dec.table.data)-1-i]
		}
		if !reflect.DeepEqual(table, test.table) {
			t.Errorf("decode test case %d: expected table %v; got %v", i, test.table, table)
		}

		if dec.table.size != test.tableSize {
			t.Errorf("decode test case %d: expected table size %v; got %v", i, test.tableSize, dec.table.size)
		}
	}
}

var (
	a              testcase
	c2             []testcase
	c3, c4, c5, c6 testcase
)

func init() {
	a = testcase{
		{`
8182 8384 8586 8788 898a 8b8c 8d8e 8f90
9192 9394 9596 9798 999a 9b9c 9d9e 9fa0
a1a2 a3a4 a5a6 a7a8 a9aa abac adae afb0
b1b2 b3b4 b5b6 b7b8 b9ba bbbc bd
`,
			huffForceTrue, indexedTrue,
			[]headerField{
				{":authority", ""},
				{":method", "GET"},
				{":method", "POST"},
				{":path", "/"},
				{":path", "/index.html"},
				{":scheme", "http"},
				{":scheme", "https"},
				{":status", "200"},
				{":status", "204"},
				{":status", "206"},
				{":status", "304"},
				{":status", "400"},
				{":status", "404"},
				{":status", "500"},
				{"accept-charset", ""},
				{"accept-encoding", "gzip, deflate"},
				{"accept-language", ""},
				{"accept-ranges", ""},
				{"accept", ""},
				{"access-control-allow-origin", ""},
				{"age", ""},
				{"allow", ""},
				{"authorization", ""},
				{"cache-control", ""},
				{"content-disposition", ""},
				{"content-encoding", ""},
				{"content-language", ""},
				{"content-length", ""},
				{"content-location", ""},
				{"content-range", ""},
				{"content-type", ""},
				{"cookie", ""},
				{"date", ""},
				{"etag", ""},
				{"expect", ""},
				{"expires", ""},
				{"from", ""},
				{"host", ""},
				{"if-match", ""},
				{"if-modified-since", ""},
				{"if-none-match", ""},
				{"if-range", ""},
				{"if-unmodified-since", ""},
				{"last-modified", ""},
				{"link", ""},
				{"location", ""},
				{"max-forwards", ""},
				{"proxy-authenticate", ""},
				{"proxy-authorization", ""},
				{"range", ""},
				{"referer", ""},
				{"refresh", ""},
				{"retry-after", ""},
				{"server", ""},
				{"set-cookie", ""},
				{"strict-transport-security", ""},
				{"transfer-encoding", ""},
				{"user-agent", ""},
				{"vary", ""},
				{"via", ""},
				{"www-authenticate", ""},
			},
			[]headerField{},
			0,
		},
	}

	c2 = []testcase{
		{
			{"400a 6375 7374 6f6d 2d6b 6579 0d63 7573 746f 6d2d 6865 6164 6572",
				huffFalse, indexedTrue,
				[]headerField{
					{"custom-key", "custom-header"},
				},
				[]headerField{
					{"custom-key", "custom-header"},
				},
				55,
			},
		},
		{
			{"040c 2f73 616d 706c 652f 7061 7468",
				huffFalse, indexedFalse,
				[]headerField{
					{":path", "/sample/path"},
				},
				[]headerField{},
				0,
			},
		},
		{
			{"1008 7061 7373 776f 7264 0673 6563 7265 74",
				huffFalse, indexedNever,
				[]headerField{
					{"password", "secret"},
				},
				[]headerField{},
				0,
			},
		},
		{
			{"82",
				huffFalse, indexedTrue,
				[]headerField{
					{":method", "GET"},
				},
				[]headerField{},
				0,
			},
		},
	}

	c3 = testcase{
		{"8286 8441 0f77 7777 2e65 7861 6d70 6c65 2e63 6f6d",
			huffFalse, indexedTrue,
			[]headerField{
				{":method", "GET"},
				{":scheme", "http"},
				{":path", "/"},
				{":authority", "www.example.com"},
			},
			[]headerField{
				{":authority", "www.example.com"},
			},
			57,
		},
		{"8286 84be 5808 6e6f 2d63 6163 6865",
			huffFalse, indexedTrue,
			[]headerField{
				{":method", "GET"},
				{":scheme", "http"},
				{":path", "/"},
				{":authority", "www.example.com"},
				{"cache-control", "no-cache"},
			},
			[]headerField{
				{"cache-control", "no-cache"},
				{":authority", "www.example.com"},
			},
			110,
		},
		{"8287 85bf 400a 6375 7374 6f6d 2d6b 6579 0c63 7573 746f 6d2d 7661 6c75 65",
			huffFalse, indexedTrue,
			[]headerField{
				{":method", "GET"},
				{":scheme", "https"},
				{":path", "/index.html"},
				{":authority", "www.example.com"},
				{"custom-key", "custom-value"},
			},
			[]headerField{
				{"custom-key", "custom-value"},
				{"cache-control", "no-cache"},
				{":authority", "www.example.com"},
			},
			164,
		},
	}

	c4 = testcase{
		{"8286 8441 8cf1 e3c2 e5f2 3a6b a0ab 90f4 ff",
			huffForceTrue, indexedTrue,
			[]headerField{
				{":method", "GET"},
				{":scheme", "http"},
				{":path", "/"},
				{":authority", "www.example.com"},
			},
			[]headerField{
				{":authority", "www.example.com"},
			},
			57,
		},
		{"8286 84be 5886 a8eb 1064 9cbf",
			huffForceTrue, indexedTrue,
			[]headerField{
				{":method", "GET"},
				{":scheme", "http"},
				{":path", "/"},
				{":authority", "www.example.com"},
				{"cache-control", "no-cache"},
			},
			[]headerField{
				{"cache-control", "no-cache"},
				{":authority", "www.example.com"},
			},
			110,
		},
		{"8287 85bf 4088 25a8 49e9 5ba9 7d7f 8925 a849 e95b b8e8 b4bf",
			huffForceTrue, indexedTrue,
			[]headerField{
				{":method", "GET"},
				{":scheme", "https"},
				{":path", "/index.html"},
				{":authority", "www.example.com"},
				{"custom-key", "custom-value"},
			},
			[]headerField{
				{"custom-key", "custom-value"},
				{"cache-control", "no-cache"},
				{":authority", "www.example.com"},
			},
			164,
		},
	}

	c5 = testcase{
		{`
4803 3330 3258 0770 7269 7661 7465 611d
4d6f 6e2c 2032 3120 4f63 7420 3230 3133
2032 303a 3133 3a32 3120 474d 546e 1768
7474 7073 3a2f 2f77 7777 2e65 7861 6d70
6c65 2e63 6f6d
`,
			huffFalse, indexedTrue,
			[]headerField{
				{":status", "302"},
				{"cache-control", "private"},
				{"date", "Mon, 21 Oct 2013 20:13:21 GMT"},
				{"location", "https://www.example.com"},
			},
			[]headerField{
				{"location", "https://www.example.com"},
				{"date", "Mon, 21 Oct 2013 20:13:21 GMT"},
				{"cache-control", "private"},
				{":status", "302"},
			},
			222,
		},
		{"4803 3330 37c1 c0bf",
			huffFalse, indexedTrue,
			[]headerField{
				{":status", "307"},
				{"cache-control", "private"},
				{"date", "Mon, 21 Oct 2013 20:13:21 GMT"},
				{"location", "https://www.example.com"},
			},
			[]headerField{
				{":status", "307"},
				{"location", "https://www.example.com"},
				{"date", "Mon, 21 Oct 2013 20:13:21 GMT"},
				{"cache-control", "private"},
			},
			222,
		},
		{`
88c1 611d 4d6f 6e2c 2032 3120 4f63 7420
3230 3133 2032 303a 3133 3a32 3220 474d
54c0 5a04 677a 6970 7738 666f 6f3d 4153
444a 4b48 514b 425a 584f 5157 454f 5049
5541 5851 5745 4f49 553b 206d 6178 2d61
6765 3d33 3630 303b 2076 6572 7369 6f6e
3d31
`,
			huffFalse, indexedTrue,
			[]headerField{
				{":status", "200"},
				{"cache-control", "private"},
				{"date", "Mon, 21 Oct 2013 20:13:22 GMT"},
				{"location", "https://www.example.com"},
				{"content-encoding", "gzip"},
				{"set-cookie", "foo=ASDJKHQKBZXOQWEOPIUAXQWEOIU; max-age=3600; version=1"},
			},
			[]headerField{
				{"set-cookie", "foo=ASDJKHQKBZXOQWEOPIUAXQWEOIU; max-age=3600; version=1"},
				{"content-encoding", "gzip"},
				{"date", "Mon, 21 Oct 2013 20:13:22 GMT"},
			},
			215,
		},
	}

	c6 = testcase{
		{`
4882 6402 5885 aec3 771a 4b61 96d0 7abe
9410 54d4 44a8 2005 9504 0b81 66e0 82a6
2d1b ff6e 919d 29ad 1718 63c7 8f0b 97c8
e9ae 82ae 43d3
`,
			huffForceTrue, indexedTrue,
			[]headerField{
				{":status", "302"},
				{"cache-control", "private"},
				{"date", "Mon, 21 Oct 2013 20:13:21 GMT"},
				{"location", "https://www.example.com"},
			},
			[]headerField{
				{"location", "https://www.example.com"},
				{"date", "Mon, 21 Oct 2013 20:13:21 GMT"},
				{"cache-control", "private"},
				{":status", "302"},
			},
			222,
		},
		{"4883 640e ffc1 c0bf",
			huffForceTrue, indexedTrue,
			[]headerField{
				{":status", "307"},
				{"cache-control", "private"},
				{"date", "Mon, 21 Oct 2013 20:13:21 GMT"},
				{"location", "https://www.example.com"},
			},
			[]headerField{
				{":status", "307"},
				{"location", "https://www.example.com"},
				{"date", "Mon, 21 Oct 2013 20:13:21 GMT"},
				{"cache-control", "private"},
			},
			222,
		},
		{`
88c1 6196 d07a be94 1054 d444 a820 0595
040b 8166 e084 a62d 1bff c05a 839b d9ab
77ad 94e7 821d d7f2 e6c7 b335 dfdf cd5b
3960 d5af 2708 7f36 72c1 ab27 0fb5 291f
9587 3160 65c0 03ed 4ee5 b106 3d50 07
`,
			huffForceTrue, indexedTrue,
			[]headerField{
				{":status", "200"},
				{"cache-control", "private"},
				{"date", "Mon, 21 Oct 2013 20:13:22 GMT"},
				{"location", "https://www.example.com"},
				{"content-encoding", "gzip"},
				{"set-cookie", "foo=ASDJKHQKBZXOQWEOPIUAXQWEOIU; max-age=3600; version=1"},
			},
			[]headerField{
				{"set-cookie", "foo=ASDJKHQKBZXOQWEOPIUAXQWEOIU; max-age=3600; version=1"},
				{"content-encoding", "gzip"},
				{"date", "Mon, 21 Oct 2013 20:13:22 GMT"},
			},
			215,
		},
	}
}
