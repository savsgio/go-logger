package logger

import (
	"io/ioutil"
	"testing"
)

func BenchmarkInfo(b *testing.B) {
	l := New(INFO, ioutil.Discard)
	l.SetEncoder(NewEncoderJSON())
	l.SetFlags(Ldatetime | Ltimestamp)

	b.Run("lineal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l.Info("hello %s", "world")
		}
	})

	b.Run("parallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Info("hello %s", "world")
			}
		})
	})
}
