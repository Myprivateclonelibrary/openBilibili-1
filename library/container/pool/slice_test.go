package pool

import (
	"context"
	"io"
	"testing"
	"time"

	xtime "go-common/library/time"

	"github.com/stretchr/testify/assert"
)

type closer struct {
}

func (c *closer) Close() error {
	return nil
}

type connection struct {
	c    io.Closer
	pool Pool
}

func (c *connection) HandleQuick() {
	//	time.Sleep(1 * time.Millisecond)
}

func (c *connection) HandleNormal() {
	time.Sleep(20 * time.Millisecond)
}

func (c *connection) HandleSlow() {
	time.Sleep(500 * time.Millisecond)
}

func (c *connection) Close() {
	c.pool.Put(context.Background(), c.c, false)
}

func TestSliceGetPut(t *testing.T) {
	// new pool
	config := &Config{
		Active:      1,
		Idle:        1,
		IdleTimeout: xtime.Duration(90 * time.Second),
		WaitTimeout: xtime.Duration(10 * time.Millisecond),
		Wait:        false,
	}
	pool := NewSlice(config)
	pool.New = func(ctx context.Context) (io.Closer, error) {
		return &closer{}, nil
	}

	// test Get Put
	conn, err := pool.Get(context.TODO())
	assert.Nil(t, err)
	c1 := connection{pool: pool, c: conn}
	c1.HandleNormal()
	c1.Close()
}

func TestSlicePut(t *testing.T) {
	var id = 0
	type connID struct {
		io.Closer
		id int
	}
	config := &Config{
		Active:      1,
		Idle:        1,
		IdleTimeout: xtime.Duration(1 * time.Second),
		//		WaitTimeout: xtime.Duration(10 * time.Millisecond),
		Wait: false,
	}
	pool := NewSlice(config)
	pool.New = func(ctx context.Context) (io.Closer, error) {
		id = id + 1
		return &connID{id: id, Closer: &closer{}}, nil
	}
	// test Put(ctx, conn, true)
	conn, err := pool.Get(context.TODO())
	assert.Nil(t, err)
	conn1 := conn.(*connID)
	// Put(ctx, conn, true) drop the connection.
	pool.Put(context.TODO(), conn, true)
	conn, err = pool.Get(context.TODO())
	assert.Nil(t, err)
	conn2 := conn.(*connID)
	assert.NotEqual(t, conn1.id, conn2.id)
}

func TestSliceIdleTimeout(t *testing.T) {
	var id = 0
	type connID struct {
		io.Closer
		id int
	}
	config := &Config{
		Active: 1,
		Idle:   1,
		// conn timeout
		IdleTimeout: xtime.Duration(1 * time.Millisecond),
	}
	pool := NewSlice(config)
	pool.New = func(ctx context.Context) (io.Closer, error) {
		id = id + 1
		return &connID{id: id, Closer: &closer{}}, nil
	}
	// test Put(ctx, conn, true)
	conn, err := pool.Get(context.TODO())
	assert.Nil(t, err)
	conn1 := conn.(*connID)
	// Put(ctx, conn, true) drop the connection.
	pool.Put(context.TODO(), conn, false)
	time.Sleep(5 * time.Millisecond)
	// idletimeout and get new conn
	conn, err = pool.Get(context.TODO())
	assert.Nil(t, err)
	conn2 := conn.(*connID)
	assert.NotEqual(t, conn1.id, conn2.id)
}

func TestSliceContextTimeout(t *testing.T) {
	// new pool
	config := &Config{
		Active:      1,
		Idle:        1,
		IdleTimeout: xtime.Duration(90 * time.Second),
		WaitTimeout: xtime.Duration(10 * time.Millisecond),
		Wait:        false,
	}
	pool := NewSlice(config)
	pool.New = func(ctx context.Context) (io.Closer, error) {
		return &closer{}, nil
	}
	// test context timeout
	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Millisecond)
	defer cancel()
	conn, err := pool.Get(ctx)
	assert.Nil(t, err)
	_, err = pool.Get(ctx)
	// context timeout error
	assert.NotNil(t, err)
	pool.Put(context.TODO(), conn, false)
	_, err = pool.Get(ctx)
	assert.Nil(t, err)
}

func TestSlicePoolExhausted(t *testing.T) {
	// test pool exhausted
	config := &Config{
		Active:      1,
		Idle:        1,
		IdleTimeout: xtime.Duration(90 * time.Second),
		//		WaitTimeout: xtime.Duration(10 * time.Millisecond),
		Wait: false,
	}
	pool := NewSlice(config)
	pool.New = func(ctx context.Context) (io.Closer, error) {
		return &closer{}, nil
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Millisecond)
	defer cancel()
	conn, err := pool.Get(context.TODO())
	assert.Nil(t, err)
	_, err = pool.Get(ctx)
	// config active == 1, so no avaliable conns make connection exhausted.
	assert.NotNil(t, err)
	pool.Put(context.TODO(), conn, false)
	_, err = pool.Get(ctx)
	assert.Nil(t, err)
}

func TestSliceStaleClean(t *testing.T) {
	var id = 0
	type connID struct {
		io.Closer
		id int
	}
	config := &Config{
		Active:      1,
		Idle:        1,
		IdleTimeout: xtime.Duration(1 * time.Second),
		//		WaitTimeout: xtime.Duration(10 * time.Millisecond),
		Wait: false,
	}
	pool := NewList(config)
	pool.New = func(ctx context.Context) (io.Closer, error) {
		id = id + 1
		return &connID{id: id, Closer: &closer{}}, nil
	}
	conn, err := pool.Get(context.TODO())
	assert.Nil(t, err)
	conn1 := conn.(*connID)
	pool.Put(context.TODO(), conn, false)
	conn, err = pool.Get(context.TODO())
	assert.Nil(t, err)
	conn2 := conn.(*connID)
	assert.Equal(t, conn1.id, conn2.id)
	pool.Put(context.TODO(), conn, false)
	// sleep more than idleTimeout
	time.Sleep(2 * time.Second)
	conn, err = pool.Get(context.TODO())
	assert.Nil(t, err)
	conn3 := conn.(*connID)
	assert.NotEqual(t, conn1.id, conn3.id)
}

func BenchmarkSlice1(b *testing.B) {
	config := &Config{
		Active:      30,
		Idle:        30,
		IdleTimeout: xtime.Duration(90 * time.Second),
		WaitTimeout: xtime.Duration(10 * time.Millisecond),
		Wait:        false,
	}
	pool := NewSlice(config)
	pool.New = func(ctx context.Context) (io.Closer, error) {
		return &closer{}, nil
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			conn, err := pool.Get(context.TODO())
			if err != nil {
				b.Error(err)
				continue
			}
			c1 := connection{pool: pool, c: conn}
			c1.HandleQuick()
			c1.Close()
		}
	})
}

func BenchmarkSlice2(b *testing.B) {
	config := &Config{
		Active:      30,
		Idle:        30,
		IdleTimeout: xtime.Duration(90 * time.Second),
		WaitTimeout: xtime.Duration(10 * time.Millisecond),
		Wait:        false,
	}
	pool := NewSlice(config)
	pool.New = func(ctx context.Context) (io.Closer, error) {
		return &closer{}, nil
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			conn, err := pool.Get(context.TODO())
			if err != nil {
				b.Error(err)
				continue
			}
			c1 := connection{pool: pool, c: conn}
			c1.HandleNormal()
			c1.Close()
		}
	})
}

func BenchmarkSlice3(b *testing.B) {
	config := &Config{
		Active:      30,
		Idle:        30,
		IdleTimeout: xtime.Duration(90 * time.Second),
		WaitTimeout: xtime.Duration(10 * time.Millisecond),
		Wait:        false,
	}
	pool := NewSlice(config)
	pool.New = func(ctx context.Context) (io.Closer, error) {
		return &closer{}, nil
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			conn, err := pool.Get(context.TODO())
			if err != nil {
				b.Error(err)
				continue
			}
			c1 := connection{pool: pool, c: conn}
			c1.HandleSlow()
			c1.Close()
		}
	})
}

func BenchmarkSlice4(b *testing.B) {
	config := &Config{
		Active:      30,
		Idle:        30,
		IdleTimeout: xtime.Duration(90 * time.Second),
		//		WaitTimeout: xtime.Duration(10 * time.Millisecond),
		Wait: false,
	}
	pool := NewSlice(config)
	pool.New = func(ctx context.Context) (io.Closer, error) {
		return &closer{}, nil
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			conn, err := pool.Get(context.TODO())
			if err != nil {
				b.Error(err)
				continue
			}
			c1 := connection{pool: pool, c: conn}
			c1.HandleSlow()
			c1.Close()
		}
	})
}

func BenchmarkSlice5(b *testing.B) {
	config := &Config{
		Active:      30,
		Idle:        30,
		IdleTimeout: xtime.Duration(90 * time.Second),
		//		WaitTimeout: xtime.Duration(10 * time.Millisecond),
		Wait: true,
	}
	pool := NewSlice(config)
	pool.New = func(ctx context.Context) (io.Closer, error) {
		return &closer{}, nil
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			conn, err := pool.Get(context.TODO())
			if err != nil {
				b.Error(err)
				continue
			}
			c1 := connection{pool: pool, c: conn}
			c1.HandleSlow()
			c1.Close()
		}
	})
}
