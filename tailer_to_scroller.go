package main

import (
	"github.com/hpcloud/tail"
	"io"
)

type TailerToScroller struct {
	tailer      *tail.Tail
	offset      int64
	currentLine *tail.Line
}

// Seek sets the offset for the next Read or Write on file to offset, interpreted
// according to whence: 0 means relative to the origin of the file, 1 means
// relative to the current offset, and 2 means relative to the end.
// It returns the new offset and an error, if any.
// The behavior of Seek on a file opened with O_APPEND is not specified.
func (tts *TailerToScroller) Seek(offset int64, whence int) (int64, error) {

	switch whence {
	case 0:
		return tts.offset, nil

	case 1:

		if offset > 0 {

			for offset > 0 {

				if tts.currentLine != nil {
					currentLength := int64(len(tts.currentLine.Text))
					if currentLength > offset {
						return tts.offset, nil
					}
				}

				nextLine := <-tts.tailer.Lines

				if nextLine.Err != nil {
					return tts.offset, nextLine.Err
				}

				tts.currentLine = nextLine
				length := int64(len(nextLine.Text))

				tts.offset = tts.offset + length
				offset = offset - length

			}

		}

		return tts.offset, nil
	case 2:

		if tts.offset == 0 {
			for line := range tts.tailer.Lines {

				if line.Err != nil {
					return tts.offset, line.Err
				}

				tts.offset += int64(len(line.Text))

			}
		}

		return tts.offset, nil

	}

	return tts.offset, nil

}

func (tts *TailerToScroller) Read(p []byte) (n int, err error) {

	nextLine := <-tts.tailer.Lines

	if nextLine.Err != nil {
		err = nextLine.Err
		return
	}

	n = copy(p, nextLine.Text)

	tts.offset += int64(n)
	if n == 0 {
		err = io.EOF
	}

	return

}
