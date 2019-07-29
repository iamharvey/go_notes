package mem_profiling

import (
	"bytes"
	"io/ioutil"
)

func findInFile(word string) (count int, err error) {

	data, err := ioutil.ReadFile("test.txt")
	if err != nil {
		return count, err
	}

	// create an input stream (a reader) for the content
	input := bytes.NewReader(data)

	// we just need the same size of word for reading bytes
	size := len(word)

	// make the buf
	buf := make([]byte, size)
	end := size - 1

	// read in an initial number of bytes we need to get started
	if n, err := input.Read(buf[:end]); err != nil || n < end {
		return count, err
	}

	for {
		// read in one byte from the input stream.
		if _, err := input.Read(buf[end:]); err != nil {
			return count, err
		}

		// if bytes matches, count +1
		if bytes.Compare(buf, []byte(word)) == 0 {
			count ++
		}

		// remove one that has been read
		copy(buf, buf[1:])

	}

	return count, err
}