package main

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"
	"github.com/nfnt/resize"
	"os"
	"strconv"
	"sync/atomic"
)

// fileID - initial ID for saved files (thumb).
// ! Is not saved between program session. //TODO
var fileID uint64 = 1  

const (
	width = 100   // thumb size
	height = 100
)

func(s *Server) saveThumb(buf []byte) (uint64, error) { 
	
	// create 100x100 thumb
	image, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		log.Println("saveThumb. Cannot decode image")
		return 0, err
	}

	newImage := resize.Resize(width, height, image, resize.NearestNeighbor)

	// save thumb into jpeg

	// increase file name/ID
	newID := atomic.AddUint64(&fileID, 1)

	newFile := s.workDir + "/" + strconv.FormatUint(newID, 10) + ".jpeg"
	f, err := os.Create(newFile)
	if err != nil {
		log.Println("saveThumb. Cannot create file: ", newFile)
		return 0, err
	}

	err = jpeg.Encode(f, newImage, nil)
	if err != nil {
		log.Println("saveThumb. Cannot save file")
		return 0, err
	}

	return newID, nil
}