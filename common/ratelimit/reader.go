// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3 with static-linking exception.
// See LICENCE file for details.

package ratelimit

import (
	"io"
	"net"
	"time"
)

type reader struct {
	r      io.Reader
	bucket *Bucket
}

// Reader returns a reader that is rate limited by
// the given token bucket. Each token in the bucket
// represents one byte.
func Reader(r io.Reader, bucket *Bucket) io.Reader {
	return &reader{
		r:      r,
		bucket: bucket,
	}
}

func (r *reader) Read(buf []byte) (int, error) {
	n, err := r.r.Read(buf)
	if n <= 0 {
		return n, err
	}
	r.bucket.Wait(int64(n))
	return n, err
}

type writer struct {
	w      io.Writer
	bucket *Bucket
}

// Writer returns a reader that is rate limited by
// the given token bucket. Each token in the bucket
// represents one byte.
func Writer(w io.Writer, bucket *Bucket) io.Writer {
	return &writer{
		w:      w,
		bucket: bucket,
	}
}

func (w *writer) Write(buf []byte) (int, error) {
	w.bucket.Wait(int64(len(buf)))
	return w.w.Write(buf)
}

type readWriteCloser struct {
	w      io.ReadWriteCloser
	bucket *Bucket
}

func ReadWriteCloser(w io.ReadWriteCloser, bucket *Bucket) io.ReadWriteCloser {
	return &readWriteCloser{
		w:      w,
		bucket: bucket,
	}
}

func (w *readWriteCloser) Write(buf []byte) (int, error) {
	w.bucket.Wait(int64(len(buf)))
	return w.w.Write(buf)
}

func (w *readWriteCloser) Read(buf []byte) (int, error) {
	n, err := w.w.Read(buf)
	if n <= 0 {
		return n, err
	}
	w.bucket.Wait(int64(n))
	return n, err
}

func (w *readWriteCloser) Close() error {
	return w.w.Close()
}

type pipe struct {
	w      net.Conn
	bucket *Bucket
}

func Conn(w net.Conn, bucket *Bucket) net.Conn {
	return &pipe{
		w:      w,
		bucket: bucket,
	}
}

func (w *pipe) Write(buf []byte) (int, error) {
	w.bucket.Wait(int64(len(buf)))
	return w.w.Write(buf)
}

func (w *pipe) Read(buf []byte) (int, error) {
	n, err := w.w.Read(buf)
	if n <= 0 {
		return n, err
	}
	w.bucket.Wait(int64(n))
	return n, err
}

func (w *pipe) Close() error {
	return w.w.Close()
}

func (w *pipe) LocalAddr() net.Addr {
	return w.w.LocalAddr()
}

func (w *pipe) RemoteAddr() net.Addr {
	return w.w.LocalAddr()
}

func (w *pipe) SetDeadline(t time.Time) error {
	return w.w.SetDeadline(t)
}

func (w *pipe) SetReadDeadline(t time.Time) error {
	return w.w.SetDeadline(t)
}

func (w *pipe) SetWriteDeadline(t time.Time) error {
	return w.w.SetWriteDeadline(t)
}
