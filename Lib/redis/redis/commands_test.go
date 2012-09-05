package redis

import (
    "log"
    "reflect"
    "testing"
    "time"
)

func TestGeneric(t *testing.T) {
    c := New("", 0, "")
    if r := SendStr(c.Rw, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }

    if res, err := c.Randomkey(); res != "" && err != nil {
        error_(t, "randomkey", "", res, err)
    }

    c.Set("foo", "foo")

    if res, err := c.Randomkey(); res != "foo" {
        error_(t, "randomkey", "foo", res, err)
    }

    ex, _ := c.Exists("key")
    nr, _ := c.Del("key")

    if (ex && nr != 1) || (!ex && nr != 0) {
        error_(t, "del", "unknown", nr, nil)
    }

    c.Set("foo", "foo")
    c.Set("bar", "bar")
    c.Set("baz", "baz")

    if nr, err := c.Del("foo", "bar", "baz"); nr != 3 {
        error_(t, "del", 3, nr, err)
    }

    c.Set("foo", "foo")

    if res, err := c.Expire("foo", 10); !res {
        error_(t, "expire", true, res, err)
    }
    if res, err := c.Persist("foo"); !res {
        error_(t, "persist", true, res, err)
    }
    if res, err := c.Ttl("foo"); res == 0 {
        error_(t, "ttl", 0, res, err)
    }
    if res, err := c.Expireat("foo", time.Now().Unix()+10); !res {
        error_(t, "expireat", true, res, err)
    }
    if res, err := c.Ttl("foo"); res <= 0 {
        error_(t, "ttl", "> 0", res, err)
    }
    if err := c.Rename("foo", "bar"); err != nil {
        error_(t, "rename", nil, nil, err)
    }
    if err := c.Rename("foo", "bar"); err == nil {
        error_(t, "rename", "error", nil, err)
    }
    if res, err := c.Renamenx("bar", "foo"); !res {
        error_(t, "renamenx", true, res, err)
    }

    c.Set("bar", "bar")

    if res, err := c.Renamenx("foo", "bar"); res {
        error_(t, "renamenx", false, res, err)
    }

    c2 := New("", 1, "")
    c2.Del("foo")
    if res, err := c.Move("foo", 1); res != true {
        error_(t, "move", true, res, err)
    }
}

func TestKeys(t *testing.T) {
    c := New("", 0, "")

    if r := SendStr(c.Rw, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }

    SendStr(c.Rw, "MSET", "foo", "one", "bar", "two", "baz", "three")

    res, err := c.Keys("foo")

    if err != nil {
        error_(t, "keys", nil, nil, err)
    }

    expected := []string{"foo"}

    if len(res) != len(expected) {
        error_(t, "keys", len(res), len(expected), nil)
    }

    for i, v := range res {
        if v != expected[i] {
            error_(t, "keys", expected[i], v, nil)
        }
    }
}

func TestSort(t *testing.T) {
    c := New("", 0, "")
    if r := SendStr(c.Rw, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }
    SendStr(c.Rw, "RPUSH", "foo", "2")
    SendStr(c.Rw, "RPUSH", "foo", "3")
    SendStr(c.Rw, "RPUSH", "foo", "1")

    res, err := c.Sort("foo")

    if err != nil {
        error_(t, "sort", nil, nil, err)
    }

    expected := []int{1, 2, 3}
    if len(res.Elems) != len(expected) {
        error_(t, "sort", len(res.Elems), len(expected), nil)
    }

    for i, v := range res.Elems {
        r := int(v.Elem.Int64())
        if r != expected[i] {
            error_(t, "sort", expected[i], v, nil)
        }
    }
}

func TestString(t *testing.T) {
    c := New("", 0, "")
    if r := SendStr(c.Rw, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }

    if res, err := c.Decr("qux"); err != nil || res != -1 {
        error_(t, "decr", -1, res, err)
    }

    if res, err := c.Decrby("qux", 1); err != nil || res != -2 {
        error_(t, "decrby", -2, res, err)
    }

    if res, err := c.Incrby("qux", 1); err != nil || res != -1 {
        error_(t, "incrby", -1, res, err)
    }

    if res, err := c.Incr("qux"); err != nil || res != 0 {
        error_(t, "incrby", 0, res, err)
    }

    if res, err := c.Setbit("qux", 0, 1); err != nil || res != 0 {
        error_(t, "setbit", 0, res, err)
    }

    if res, err := c.Getbit("qux", 0); err != nil || res != 1 {
        error_(t, "getbit", 1, res, err)
    }

    if err := c.Set("foo", "foo"); err != nil {
        t.Errorf(err.Error())
    }

    if res, err := c.Append("foo", "bar"); err != nil || res != 6 {
        error_(t, "append", 6, res, err)
    }

    if res, err := c.Get("foo"); err != nil || res.String() != "foobar" {
        error_(t, "get", "foobar", res, err)
    }

    if v, err := c.Get("nil"); err != nil && v != nil {
        error_(t, "get", "error", nil, err)
    }

    if res, err := c.Getrange("foo", 0, 2); err != nil || res.String() != "foo" {
        error_(t, "getrange", "foo", res, err)
    }

    if res, err := c.Setrange("foo", 0, "qux"); err != nil || res != 6 {
        error_(t, "setrange", 6, res, err)
    }

    if res, err := c.Getset("foo", "foo"); err != nil || res.String() != "quxbar" {
        error_(t, "getset", "quxbar", res, err)
    }

    if res, err := c.Setnx("foo", "bar"); err != nil || res != false {
        error_(t, "setnx", false, res, err)
    }

    if res, err := c.Strlen("foo"); err != nil || res != 3 {
        error_(t, "strlen", 3, res, err)
    }

    if err := c.Setex("foo", 10, "bar"); err != nil {
        error_(t, "setex", nil, nil, err)
    }

    out := []string{"foo", "bar", "qux"}
    in := map[string]string{"foo": "foo", "bar": "bar", "qux": "qux"}

    if err := c.Mset(in); err != nil {
        error_(t, "mset", nil, nil, err)
    }

    if res, err := c.Msetnx(in); err != nil || res == true {
        error_(t, "msetnx", false, res, err)
    }

    res, err := c.Mget(out...)

    if err != nil || len(res.Elems) != 3 {
        error_(t, "mget", 3, len(res.Elems), err)
        t.FailNow()
    }

    for i, v := range res.StringArray() {
        if v != out[i] {
            error_(t, "mget", out[i], v, nil)
        }
    }

    out = append(out, "il")
    if res, err = c.Mget(out...); err != nil && res != nil {
        error_(t, "mget", nil, "expected error", err)
        t.FailNow()
    }
}

func TestList(t *testing.T) {
    c := New("", 0, "")
    if r := SendStr(c.Rw, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }

    if res, err := c.Lpush("foobar", "foo"); err != nil || res != 1 {
        error_(t, "LPUSH", 1, res, err)
    }

    if res, err := c.Linsert("foobar", "AFTER", "foo", "bar"); err != nil || res != 2 {
        error_(t, "Linsert", 2, res, err)
    }

    if res, err := c.Linsert("foobar", "AFTER", "qux", "bar"); err != nil || res != -1 {
        error_(t, "Linsert", -1, res, err)
    }

    if res, err := c.Llen("foobar"); err != nil || res != 2 {
        error_(t, "Llen", 2, res, err)
    }

    if res, err := c.Lindex("foobar", 0); err != nil || res.String() != "foo" {
        error_(t, "Lindex", "foo", res, err)
    }

    if res, err := c.Lpush("foobar", "qux"); err != nil || res != 3 {
        error_(t, "Lpush", 3, res, err)
    }

    if res, err := c.Lpop("foobar"); err != nil || res.String() != "qux" {
        error_(t, "Lpop", "qux", res, err)
    }

    want1 := []string{"foo", "bar"}

    if out, err := c.Lrange("foobar", 0, 1); err != nil || !reflect.DeepEqual(want1, out.StringArray()) {
        error_(t, "Lrange", nil, nil, err)
    }

    want := []string{"foo"}

    if res, err := c.Lrem("foobar", 0, "bar"); err != nil || res != 1 {
        error_(t, "Lrem", 1, res, err)
    }

    want = []string{"bar"}

    if err := c.Lset("foobar", 0, "bar"); err != nil {
        error_(t, "Lrem", nil, nil, err)
    }

    want = []string{}

    if err := c.Ltrim("foobar", 1, 0); err != nil {
        error_(t, "Ltrim", nil, nil, err)
    }

    want = []string{"foo", "bar", "qux"}
    var res int64
    var err error

    for _, v := range want {
        res, err = c.Rpush("foobar", v)
    }

    if err != nil || res != 3 {
        error_(t, "Rpush", 3, res, err)
    }

    if res, err := c.Rpushx("foobar", "baz"); err != nil || res != 4 {
        error_(t, "Rpushx", 4, res, err)
    }

    if res, err := c.Rpop("foobar"); err != nil || res.String() != "baz" {
        error_(t, "Rpop", "baz", res, err)
    }

    if res, err := c.Rpoplpush("foobar", "foobaz"); err != nil || res.String() != "qux" {
        error_(t, "Rpoplpush", "qux", res, err)
    }

    if res, err := c.Blpop([]string{"foobar", "foobaz"}, 1); err != nil || res.StringMap()["foobar"] != "foo" {
        error_(t, "Blpop", "foo", res.StringMap()["foobar"], err)
    }

    if res, err := c.Brpop([]string{"foobar"}, 1); err != nil || res.StringMap()["foobar"] != "bar" {
        error_(t, "Brpop", "bar", res, err)
    }

    if res, err := c.Blpop([]string{"foobar", "foobaz"}, 1); err != nil || res.StringMap()["foobaz"] != "qux" {
        error_(t, "Blpop", "bar", res, err)
    }

    if res, err := c.Blpop([]string{"foobar", "foobaz"}, 1); err == nil {
        error_(t, "Blpop timeout err", nil, res, err)
    }

    if res, err := c.Brpoplpush("foobar", "foobaz", 1); err == nil {
        error_(t, "brpoplpush timeout err", nil, res, err)
    }
}

func TestHash(t *testing.T) {
    c := New("", 0, "")
    if r := SendStr(c.Rw, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }

    if res, err := c.Hset("foobar", "foo", "foo"); err != nil || res != true {
        error_(t, "Hset", true, res, err)
    }

    if res, err := c.Hset("foobar", "foo", "foo"); err != nil || res != false {
        error_(t, "Hset", false, res, err)
    }

    if res, err := c.Hget("foobar", "foo"); err != nil || res.String() != "foo" {
        error_(t, "Hget", "foo", res, err)
    }

    if res, err := c.Hdel("foobar", "foo"); err != nil || res != true {
        error_(t, "Hdel", true, res, err)
    }

    if res, err := c.Hexists("foobar", "foo"); err != nil || res != false {
        error_(t, "Hexists", false, res, err)
    }

    if res, err := c.Hsetnx("foobar", "foo", 1); err != nil || res != true {
        error_(t, "Hsetnx", true, res, err)
    }
    c.Hset("foobar", "bar", 2)

    want := []string{"foo", "1", "bar", "2"}

    if res, err := c.Hgetall("foobar"); err != nil || !reflect.DeepEqual(want, res.StringArray()) {
        error_(t, "Hgetall", want, res, err)
    }

    if res, err := c.Hincrby("foobar", "foo", 1); err != nil || int64(2) != res {
        error_(t, "Hincrby", int64(2), res, err)
    }

    want1 := []string{"foo", "bar"}

    if res, err := c.Hkeys("foobar"); err != nil || !reflect.DeepEqual(want1, res) {
        error_(t, "Hkeys", want1, res, err)
    }

    if res, err := c.Hlen("foobar"); err != nil || int64(2) != res {
        error_(t, "Hlen", int64(2), res, err)
    }

    if res, err := c.Hlen("foobar"); err != nil || int64(2) != res {
        error_(t, "Hlen", int64(2), res, err)
    }

    want = []string{"2"}

    if res, err := c.Hmget("foobar", "bar"); err != nil || !reflect.DeepEqual(want, res.StringArray()) {
        error_(t, "Hmget", want, res, err)
    }

    m := map[string]interface{}{
        "foo": 1,
        "bar": 2,
        "qux": 3,
    }

    if err := c.Hmset("foobar", m); err != nil {
        error_(t, "Hmset", nil, nil, err)
    }

    want2 := []int64{1, 2, 3}
    if res, err := c.Hvals("foobar"); err != nil || !reflect.DeepEqual(want2, res.IntArray()) {
        error_(t, "Hvals", want2, res, err)
    }
}

func TestSet(t *testing.T) {
    c := New("", 0, "")
    if r := SendStr(c.Rw, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }

    if res, err := c.Sadd("foobar", "foo"); err != nil || res != true {
        error_(t, "Sadd", true, res, err)
    }

    if res, err := c.Sadd("foobar", "foo"); err != nil || res != false {
        error_(t, "Sadd", false, res, err)
    }

    if res, err := c.Scard("foobar"); err != nil || res != 1 {
        error_(t, "Scard", 1, res, err)
    }

    c.Sadd("foobar", "bar")
    c.Sadd("foobaz", "foo")

    want := []string{"bar", "foo", "bar"}

    switch res, err := c.Sunion("foobar", "foobaz"); {
    case err != nil:
    case !reflect.DeepEqual(want[:2], res.StringArray()):
    case !reflect.DeepEqual(want[1:], res.StringArray()):
        error_(t, "Sunion", want[:2], res, err)
    }

    if res, err := c.Sunionstore("fooqux", "foobar", "foobaz"); err != nil || res != 2 {
        error_(t, "Sunionstore", 2, res, err)
    }

    want = []string{"bar"}

    if res, err := c.Sdiff("foobar", "foobaz"); err != nil || !reflect.DeepEqual(want, res.StringArray()) {
        error_(t, "Sdiff", want, res, err)
    }

    if res, err := c.Sdiffstore("foobar", "foobaz"); err != nil || res != 1 {
        error_(t, "Sdiffstore", 1, res, err)
    }

    want = []string{"foo"}

    if res, err := c.Sinter("foobar", "foobaz"); err != nil || !reflect.DeepEqual(want, res.StringArray()) {
        error_(t, "Sinter", want, res, err)
    }

    if res, err := c.Sinterstore("foobar", "foobaz"); err != nil || res != 1 {
        error_(t, "Sinterstore", 1, res, err)
    }

    if res, err := c.Sismember("foobar", "qux"); err != nil || res != false {
        error_(t, "Sismember", false, res, err)
    }

    if res, err := c.Smembers("foobaz"); err != nil || !reflect.DeepEqual(want, res.StringArray()) {
        error_(t, "smembers", want, res, err)
    }

    if res, err := c.Smove("foobar", "foobaz", "foo"); err != nil || res != true {
        error_(t, "smove", true, res, err)
    }

    if res, err := c.Spop("foobaz"); err != nil || res.String() != "foo" {
        error_(t, "spop", "foo", res, err)
    }

    if res, err := c.Srandmember("foobaz"); err != nil && res != nil {
        error_(t, "srandmember", nil, res, err)
    }

    c.Sadd("foobar", "foo")
    c.Sadd("foobar", "bar")
    c.Sadd("foobar", "baz")

    if res, err := c.Srem("foobar", "baz"); err != nil || res != true {
        error_(t, "srem", nil, res, err)
    }
}

func TestSortedSet(t *testing.T) {
    c := New("", 0, "")
    if r := SendStr(c.Rw, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }

    m := map[string]float64{
        "foo": 1.0,
        "bar": 2.0,
        "baz": 3.0,
        "qux": 4.0,
    }

    for k, v := range m {
        if res, err := c.Zadd("foobar", v, k); err != nil || res != true {
            error_(t, "Zadd", true, res, err)
        }
    }

    if res, err := c.Zcard("foobar"); err != nil || res != 4 {
        error_(t, "Zcard", 4, res, err)
    }

    if res, err := c.Zcount("foobar", 1, 2); err != nil || res != 2 {
        error_(t, "Zcount", 2, res, err)
    }

    if res, err := c.Zincrby("foobar", 0.5, "foo"); err != nil || res != 1.5 {
        error_(t, "Zincrby", 1.5, res, err)
    }

    if res, err := c.Zinterstore("barbaz", []string{"foobar"}); err != nil || res != 4 {
        error_(t, "Zinterstore", 4, res, err)
    }

    want := []string{"foo", "bar", "baz", "qux"}

    if res, err := c.Zrange("foobar", 0, 4); err != nil || !reflect.DeepEqual(want, res.StringArray()) {
        error_(t, "Zrange", want, res, err)
    }

    if res, err := c.Zrangebyscore("foobar", "0", "+inf"); err != nil || !reflect.DeepEqual(want, res.StringArray()) {
        error_(t, "Zrangebyscore", want, res, err)
    }

    if res, err := c.Zrank("foobar", "baz"); err != nil || res != 2 {
        error_(t, "Zrank", 2, res, err)
    }

    // if res, err := c.Zrank("foobar", "nil"); err == nil || res != 0  {
    //    error(t, "Zrank", 0, res, err)
    // }

    if res, err := c.Zrem("foobar", "bar"); err != nil || res != true {
        error_(t, "Zrem", true, res, err)
    }

    if res, err := c.Zrem("foobar", "bar"); err != nil || res != false {
        error_(t, "Zrem", false, res, err)
    }

    if res, err := c.Zremrangebyrank("foobar", 0, 0); err != nil || res != 1 {
        error_(t, "Zremrangebyrank", 1, res, err)
    }

    if res, err := c.Zremrangebyscore("foobar", 0, 3); err != nil || res != 1 {
        error_(t, "zremrangebyscore", 1, res, err)
    }

    c.Zinterstore("foobar", []string{"barbaz"})
    want = []string{"qux", "baz", "bar", "foo"}

    if res, err := c.Zrevrange("foobar", 0, 4); err != nil || !reflect.DeepEqual(want, res.StringArray()) {
        error_(t, "Zrevrange", want, res.StringArray(), err)
    }

    want2 := map[string]string{"qux": "4", "baz": "3", "bar": "2", "foo": "1.5"}
    if res, err := c.Zrevrangebyscore("foobar", 4, 0, "WITHSCORES"); err != nil || !reflect.DeepEqual(want2, res.StringMap()) {
        error_(t, "Zrevrangebyscore", want, res.StringMap(), err)
    }

    if res, err := c.Zrevrank("foobar", "baz"); err != nil || res != 1 {
        error_(t, "Zrevrank", 1, res, err)
    }

    if res, err := c.Zscore("foobar", "foo"); err != nil || res != 1.5 {
        error_(t, "Zscore", 1.5, res, err)
    }

    if res, err := c.Zunionstore("foobar", []string{"nil"}); err == nil || res != -1 {
        error_(t, "Zunionstore", -1, res, err)
    }
}

func TestConnection(t *testing.T) {
    c := New("", 0, "")
    if r := SendStr(c.Rw, "flushall"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }

    if res, err := c.Echo("foo"); err != nil || res.String() != "foo" {
        error_(t, "Echo", "foo", res, err)
    }

    if res, err := c.Ping(); err != nil || res.String() != "PONG" {
        error_(t, "Ping", "PONG", res, err)
    }

    c.Set("foo", "foo")

    if err := c.Select(2); err != nil {
        error_(t, "select", nil, nil, err)
    }

    if res, err := c.Get("foo"); err != nil && res != nil {
        error_(t, "get-select", nil, res, err)
    }

    for i := 0; i < MaxClientConn; i++ {
        if err := c.Quit(); err != nil {
            error_(t, "quit", nil, nil, err)
        }
    }

    if err := c.Set("foo", "foo"); err != nil {
        error_(t, "quit", nil, nil, err)
    }
}

func TestServer(t *testing.T) {
    c := New("", 0, "")

    if res, err := c.Monitor(); err != nil {
        error_(t, "monitor", nil, res, err)
    }

    if err := c.ConfigSet("requirepass", "foobared"); err != nil {
        error_(t, "requirepass foobared", nil, nil, err)
    }

    c = New("", 1, "foobared")

    if _, err := c.Ping(); err != nil {
        error_(t, "ping authenticated", nil, nil, err)
    }

    if r, err := c.ConfigGet("requirepass"); r.StringMap()["requirepass"] != "foobared" || err != nil {
        error_(t, "config get", "foobared", r.StringMap(), err)
    }

    if err := c.ConfigSet("requirepass", ""); err != nil {
        error_(t, "requirepass reset", nil, nil, err)
    }

    if err := c.ConfigResetstat(); err != nil {
        error_(t, "config resetstat", nil, nil, err)
    }

    // Don't bother to test each time
    // since it changes the state of server.
    //if err := c.Slaveof("localhost", "6379"); err != nil {
    //    error_(t, "slavof", nil, nil, err)
    //}
    //
    //if err := c.Slaveof("localhost", "NO ONE"); err != nil {
    //    error_(t, "slavof", nil, nil, err)
    //}
}

func TestPubSub(t *testing.T) {
    c := New("", 0, "")

    if r := SendStr(c.Rw, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }

    sub, err := c.Subscribe("foochan", "barchan")

    if err != nil {
        t.Fatal("subscribe", nil, nil, err)
    }

    if res, err := c.Publish("foochan", "foo"); err != nil || res != 1 {
        error_(t, "publish", 1, res, err)
    }

    go func() {
        m := <-sub.Messages
        if m.Elem.String() != "foo" || m.Channel != "foochan" {
            error_(t, "res-subscribe", "foo "+"foochan", m.Elem.String()+m.Channel, nil)
        }
    }()

    time.Sleep(1e8)
    if res, err := c.Publish("barchan", "bar"); err != nil || res != 1 {
        error_(t, "publish", 1, res, err)
    }

    time.Sleep(1e8)

    if m := <-sub.Messages; m.Elem.String() != "bar" || m.Channel != "barchan" {
        error_(t, "subscribe barchan", "bar", m, nil)
    }

    if err := sub.Unsubscribe("foochan"); err != nil {
        error_(t, "unsubscribe", nil, nil, err)
    }

    if res, err := c.Publish("foochan", "foo"); err != nil || res != 0 {
        error_(t, "publish", 0, res, err)
    }

    if res, err := c.Publish("barchan", "bar"); err != nil || res != 1 {
        error_(t, "publish", 1, res, err)
    }

    if m := <-sub.Messages; m.Elem.String() != "bar" || m.Channel != "barchan" {
        error_(t, "subscribe barchan", "bar", m.Elem.String(), nil)
    }

    if err := sub.Psubscribe("*chan"); err != nil {
        error_(t, "psubscribe", nil, nil, err)
    }

    if res, err := c.Publish("bazchan", "baz"); err != nil || res != 1 {
        error_(t, "publish", 1, res, err)
    }

    if m := <-sub.Messages; m.Elem.String() != "baz" || m.Channel != "bazchan" {
        error_(t, "psubscribe bazchan", "baz", m.Elem.String(), nil)
    }

    sub.Close()
    time.Sleep(1e8)

    if _, ok := <-sub.Messages; ok != false {
        error_(t, "closed chan", false, ok, nil)
    }
}

func TestTransaction(t *testing.T) {
    c := New("", 0, "")

    if r := SendStr(c.Rw, "FLUSHDB"); r.Err != nil {
        t.Fatalf("'%s': %s", "FLUSHDB", r.Err)
    }

    c.Set("foo", "bar")

    p := NewPipeClientFromClient(c)

    if err := p.Watch("foo"); err != nil {
        error_(t, "watch", nil, nil, err)
    }

    // A different client modifies "foo" after watch was called
    c.Set("foo", "qux")

    if err := p.Multi(); err != nil {
        error_(t, "multi", nil, nil, err)
    }

    if err := p.Set("foo", "foo"); err != nil {
        error_(t, "set", nil, nil, err)
    }

    if replies := p.Exec(); len(replies) > 0 {
        error_(t, "exec watched", 0, replies[0].Elem, nil)
    }

    if err := p.Watch("foo"); err != nil {
        error_(t, "watch", nil, nil, err)
    }

    if err := p.Multi(); err != nil {
        error_(t, "multi", nil, nil, err)
    }

    if err := p.Set("foo", "foo"); err != nil {
        error_(t, "set", nil, nil, err)
    }

    // "foo" was not modified after watch was called
    if replies := p.Exec(); len(replies) != 1 {
        error_(t, "exec watched", 0, replies[0].Elem, nil)
    }

    if err := p.Watch("foo"); err != nil {
        error_(t, "watch", nil, nil, err)
    }

    if err := p.Unwatch(); err != nil {
        error_(t, "watch", nil, nil, err)
    }

    // A different client modifies "foo" after Unwatch was called
    c.Set("foo", "qux")

    if err := p.Multi(); err != nil {
        error_(t, "multi", nil, nil, err)
    }

    if err := p.Set("foo", "foo"); err != nil {
        error_(t, "set", nil, nil, err)
    }

    // "foo" was not modified after watch was called
    if replies := p.Exec(); len(replies) != 1 {
        error_(t, "exec watched", 0, replies[0].Elem, nil)
    }
}

func BenchmarkRpush(b *testing.B) {
    c := New("", 0, "")
    start := time.Now()
    for i := 0; i < b.N; i++ {
        if _, err := c.Rpush("qux", "qux"); err != nil {
            log.Println("RPUSH", err)
            return
        }
    }
    c.Del("qux")
    stop := time.Now().Sub(start)
    log.Printf("time: %.4f\n", float32(stop/1.0e+6)/1000.0)
}

func BenchmarkRpushPiped(b *testing.B) {
    c := NewPipeClient("", 0, "")
    start := time.Now()

    for i := 0; i < b.N; i++ {
        if _, err := c.Rpush("qux", "qux"); err != nil {
            log.Println("RPUSH", err)
            return
        }
    }

    c.Del("qux")
    c.Exec()
    stop := time.Now().Sub(start)
    log.Printf("time: %.4f\n", float32(stop/1.0e+6)/1000.0)
}
