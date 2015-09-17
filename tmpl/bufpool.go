package tmpl

import (
    "bytes"
)

type bufferPool struct {
    ch chan *bytes.Buffer
}

func (bp *bufferPool) get() *bytes.Buffer {
    var b *bytes.Buffer
    select {
        case b = <-bp.ch:
        default:
            b = bytes.NewBuffer([]byte{})
    }
    return b
}

func (bp *bufferPool) reset(b *bytes.Buffer) {
    b.Reset()
    select {
        case bp.ch <-b:
        default:
            // discard buffer if pool is full
    }
}
