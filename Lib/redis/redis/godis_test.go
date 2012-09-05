package redis

import (
    "bufio"
    "bytes"
    "errors"
    "log"
    "reflect" //"strconv"

    "testing"
    "time"
)

func error_(t *testing.T, name string, expected, got interface{}, err error) {
    t.Errorf("`%s` expected `%v` got `%v`, err(%v)", name, expected, got, err)
}

func printCmdCount() {
    log.Println("command | count ")
    for k, v := range cmdCount {
        log.Printf("      %c | %d\n", k, v)
    }
}

func printRes(t *testing.T, r *Reply) {
    if len(r.Elems) > 0 {
        t.Logf("str arr: %q", r.StringArray())
    } else {
        t.Logf("str: %q", r.Elem.String())
    }
    if r.Err != nil {
        t.Logf("err: %q", r.Err)
    }
}

func compareReply(t *testing.T, name string, a, b *Reply) {
    if a.Err != nil && b.Err == nil {
        t.Fatalf("'%s': expected error `%v`", name, a.Err)
    } else if b.Err != nil && b.Err.Error() != b.Err.Error() {
        t.Fatalf("'%s': expected %s got %v", name, a.Err, b.Err)
    } else if b.Elem != nil {
        for i, c := range a.Elem {
            if c != b.Elem[i] {
                t.Errorf("'%s': expected %v got %v", name, b, a)
            }
        }
    } else if b.Elems != nil {
        for i, rep := range a.Elems {
            for j, e := range rep.Elem {
                if e != b.Elems[i].Elem[j] {
                    t.Errorf("expected %v got %v", b, a)
                    break
                }
            }
        }
    }
}

type simpleParserTest struct {
    in   string
    out  Reply
    name string
}

type redisReadWriter struct {
    writer *bufio.Writer
    reader *bufio.Reader
}

func dummyReadWriter(data string) *conn {
    br := bufio.NewReader(bytes.NewBufferString(data))
    bw := bufio.NewWriter(bytes.NewBufferString(data))
    return &conn{rwc: nil, r: br, w: bw}
}

var simpleParserTests = []simpleParserTest{
    {"+OK\r\n", Reply{Elem: []byte("OK")}, "ok"},
    {"-ERR\r\n", Reply{Err: errors.New("ERR")}, "err"},
    {":1\r\n", Reply{Elem: []byte("1")}, "num"},
    {"$3\r\nfoo\r\n", Reply{Elem: []byte("foo")}, "bulk"},
    {"$-1\r\n", Reply{Err: errors.New("Nonexisting Key")}, "bulk-nil"},
    {"*-1\r\n", Reply{}, "multi-bulk-nil"},
}

func TestParser(t *testing.T) {
    for _, test := range simpleParserTests {
        rw := dummyReadWriter(test.in)
        r := rw.readReply()
        compareReply(t, test.name, r, &test.out)
        t.Log(test.in, r, test.out)
    }
}

func s2MultiReply(ss ...string) []*Reply {
    var r = make([]*Reply, len(ss))
    for i := 0; i < len(ss); i++ {
        r[i] = &Reply{Elem: []byte(ss[i])}
    }
    return r
}

type SimpleSendTest struct {
    cmd  string
    args []string
    out  Reply
}

var simpleSendTests = []SimpleSendTest{
    {"FLUSHDB", []string{}, Reply{Elem: []byte("OK")}},
    {"SET", []string{"key", "foo"}, Reply{Elem: []byte("OK")}},
    {"EXISTS", []string{"key"}, Reply{Elem: []byte("1")}},
    {"GET", []string{"key"}, Reply{Elem: []byte("foo")}},
    {"GET", []string{"/dev/null"}, Reply{}},
    {"RPUSH", []string{"list", "foo"}, Reply{Elem: []byte("1")}},
    {"RPUSH", []string{"list", "bar"}, Reply{Elem: []byte("2")}},
    {"LRANGE", []string{"list", "0", "2"}, Reply{Elems: s2MultiReply("foo", "bar")}},
    {"KEYS", []string{"list"}, Reply{Elems: s2MultiReply("list")}},
}

func TestSimpleSend(t *testing.T) {
    c := New("", 0, "")
    for _, test := range simpleSendTests {
        r := SendStr(c.Rw, test.cmd, test.args...)
        compareReply(t, test.cmd, &test.out, r)
        t.Log(test.cmd, test.args)
        t.Logf("%q == %q\n", test.out, r)
    }
}

func equals(a, b []byte) bool {
    for i, c := range a {
        if c != b[i] {
            return false
        }
    }
    return true
}

func TestBinarySafe(t *testing.T) {
    c := New("", 0, "")
    want1 := make([]byte, 256)
    for i := 0; i < 256; i++ {
        want1[i] = byte(i)
    }

    Send(c.Rw, []byte("SET"), []byte("foo"), want1)

    if res := Send(c.Rw, []byte("GET"), []byte("foo")); !equals(res.Elem.Bytes(), want1) {
        error_(t, "ascii-table-Send", want1, res.Elem.Bytes(), res.Err)
    }

    SendIface(c.Rw, "SET", "bar", string(want1))

    if res := SendIface(c.Rw, "GET", "bar"); !equals(res.Elem.Bytes(), want1) {
        error_(t, "ascii-table-SendIface", want1, res.Elem.Bytes(), res.Err)
    }

    want2 := []byte("♥\r\nµs\r\n")
    Send(c.Rw, []byte("SET"), []byte("foo"), want2)

    if res := Send(c.Rw, []byte("GET"), []byte("foo")); !equals(res.Elem.Bytes(), want2) {
        error_(t, "unicode-Send", want2, res.Elem.Bytes(), res.Err)
    }

    SendIface(c.Rw, "SET", "bar", want2)

    if res := SendIface(c.Rw, "GET", "bar"); !equals(res.Elem.Bytes(), want2) {
        error_(t, "unicode-SendIface", want2, res.Elem.Bytes(), res.Err)
    }

    for _, b := range want2 {
        SendIface(c.Rw, "SET", "bar", b)
        res := SendIface(c.Rw, "GET", "bar")
        if uint8(res.Elem.Int64()) != b {
            error_(t, "unicode-SendIface", b, res.Elem, res.Err)
        }
    }
}

func TestSimplePipe(t *testing.T) {
    c := NewPipeClient("", 0, "")

    for _, test := range simpleSendTests {
        r := SendStr(c.Rw, test.cmd, test.args...)

        if r.Err != nil {
            t.Error(test.cmd, r.Err, test.args)
        }
    }

    replies := c.Exec()

    if len(replies) != len(simpleSendTests) {
        error_(t, "pipe replies len", len(simpleSendTests), len(replies), nil)
    }

    for i, test := range simpleSendTests {
        compareReply(t, test.cmd, &test.out, replies[i])
    }
}

func TestSimpleTransaction(t *testing.T) {
    c := New("", 0, "")

    if r := SendStr(c.Rw, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }

    p := NewPipeClientFromClient(c)
    p.Multi()
    p.Set("foo", "bar")
    p.Set("bar", "bar")
    p.Lpush("bar", "bar")
    p.Get("bar")
    replies := p.Exec()

    t.Log(replies)
    t.Log(replies[2].Err)

    pc := NewPipeClient("", 0, "")
    pc.Set("baz", "baz")
    pc.Exec()
}

//func TestPipeConn(t *testing.T) {
//    c := NewPipe("", 0, "")
//
//    if r := SendStr(c.Rw, "FLUSHDB"); r.Err != nil {
//        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
//    }
//
//    if r := SendStr(c.Rw, "SET", "foo", "foo"); r.Elem != nil {
//        error_(t, "PIPE-SET", nil, r.Elem, r.Err)
//    }
//
//    want := []byte("OK")
//
//    if r := c.GetReply(); !reflect.DeepEqual(r.Elem.Bytes(), want) {
//        error_(t, "PIPE-GET-FLUSHDB", want, r.Elem, r.Err)
//    }
//
//    if r := SendStr(c.Rw, "SET", "bar", "bar"); r.Elem != nil {
//        error_(t, "PIPE-SET", nil, r.Elem, r.Err)
//    }
//
//    if r := c.GetReply(); !reflect.DeepEqual(r.Elem.Bytes(), want) {
//        error_(t, "PIPE-GET-SET", want, r.Elem, r.Err)
//    }
//
//    if r := c.GetReply(); !reflect.DeepEqual(r.Elem.Bytes(), want) {
//        error_(t, "PIPE-GET-SET", want, r.Elem, r.Err)
//    }
//
//    if r := c.GetReply(); r.Err == nil {
//        error_(t, "PIPE-GET-SET", nil, r.Elem, nil)
//    }
//}

func TestMemory(t *testing.T) {
    c := New("", 0, "")
    n := 2
    Send(c.Rw, []byte("FLUSHDB"))

    for i := 0; i < 5; i++ {
        SendIface(c.Rw, "RPUSH", "list", i)
    }

    //time.Sleep(1.0e+9 * 10)
    start := time.Now()
    replies := make([]*Reply, n)

    for i := 0; i < n; i++ {
        replies[i], _ = c.Lrange("list", 0, 4)
    }

    stop := time.Now().Sub(start)
    if debug {
        log.Printf("time: %.2f\n", float32(stop/1.0e+9))
    }
    //time.Sleep(1.0e+9 * 10)
    //printCmdCount()
}

// for this test to work redis.conf has to be set timeout to 1sec
// the test return a nil pointer if failed
func TestConnTimeout(t *testing.T) {
    c := New("", 0, "")
    Send(c.Rw, []byte("FLUSHDB"))

    defer func() {
        if x := recover(); x != nil {
            t.Errorf("`conn timeout` expected got `%v`", x)
        }
    }()

    c.Set("foo", 1)
    c.Set("bar", 2)

    time.Sleep(1e+9 * 8)

    rep, err := c.Mget("foo", "bar")
    // rep.IntArray will invoke a nil-pointer panic if there was an err
    rep.IntArray()

    if err != nil {
        error_(t, "connection timeout", nil, nil, err)
    }
}

func TestReadingBulk(t *testing.T) {
    c := New("", 0, "")

    if r := SendStr(c.Rw, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }

    var want3 []int64

    for i := 0; i < 600; i++ {
        want3 = append(want3, int64(i))
        c.Rpush("foobaz", i)

        if res, err := c.Lrange("foobaz", 0, i); err != nil || !reflect.DeepEqual(want3, res.IntArray()) {
            error_(t, "Lranges", nil, nil, err)
            t.FailNow()
        }
    }
}

func BenchmarkParsing(b *testing.B) {
    c := New("", 0, "")

    for i := 0; i < 1000; i++ {
        SendStr(c.Rw, "RPUSH", "list", "foo")
    }

    start := time.Now()

    for i := 0; i < b.N; i++ {
        SendStr(c.Rw, "LRANGE", "list", "0", "50")
    }

    stop := time.Now().Sub(start)

    log.Printf("time: %.2f\n", float32(stop/1.0e+9))
    Send(c.Rw, []byte("FLUSHDB"))
}

//func TestBenchmark(t *testing.T) {
//    c := New("", 0, "")
//    c.Send("FLUSHDB")
//    start := time.Nanoseconds()
//    n := 2000000
//
//    a, b := []byte("zrs"), []byte("hi")
//    for i := 0; i < n; i++ {
//        c.Send("RPUSH", a, b)
//    }
//
//    //c.Del("zrs")
//    stop := time.Nanoseconds() - start
//
//    ti := float32(stop / 1.0e+6) / 1000.0
//    fmt.Fprintf(os.Stdout, "godis: %.2f %.2f per/s\n", ti, float32(n) / ti)
//}
