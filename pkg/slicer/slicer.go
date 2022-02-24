package slicer

func filesPerChunk(files []string, workers int) int {
	var filesPerChunk = len(files) / workers
	if filesPerChunk == 0 {
		filesPerChunk = 1
	}
	return filesPerChunk
}

func ChunkSlice(files []string, workers int) [][]string {
	var chunkSize = filesPerChunk(files, workers)
	var chunks [][]string
	for i := 0; i < len(files); i += chunkSize {
		end := i + chunkSize

		if end > len(files) {
			end = len(files)
		}

		chunks = append(chunks, files[i:end])
	}

	return chunks
}
