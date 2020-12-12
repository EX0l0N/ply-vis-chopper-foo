package visio

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
)

type plyvis [][]byte

func ReadVis(in io.Reader) plyvis {
	var vislen uint64

	bin := bufio.NewReader(in)

	if err := binary.Read(bin, binary.LittleEndian, &vislen); err != nil {
		panic(fmt.Sprint("Unable to read vis: ", err))
	}

	vis := make(plyvis, vislen)

	for c := uint64(0); c < vislen; c++ {
		var num_ref_imgs uint32

		if err := binary.Read(bin, binary.LittleEndian, &num_ref_imgs); err != nil {
			panic(fmt.Sprint("Could not retrieve number of images for this point: ", err))
		}

		vis[c] = make([]byte, 4*num_ref_imgs)

		if _, err := io.ReadFull(bin, vis[c]); err != nil {
			panic(fmt.Sprint("Could not read images references: ", err))
		}
	}

	return vis
}

func WriteVis(pv plyvis, out io.Writer) {
	but := bufio.NewWriter(out)
	defer but.Flush()

	if err := binary.Write(but, binary.LittleEndian, uint64(len(pv))); err != nil {
		panic(fmt.Sprint("Unable to write vis: ", err))
	}

	for _, ri := range pv {
		if err := binary.Write(but, binary.LittleEndian, uint32(len(ri)/4)); err != nil {
			panic(fmt.Sprint("Could not write number of reference images for this point: ", err))
		}

		if _, err := but.Write(ri); err != nil {
			panic(fmt.Sprint("Could not write image references: ", err))
		}
	}
}
