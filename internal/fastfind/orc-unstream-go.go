/*
 * orc-unstream - Convert ORC stream/journal format to standard 7zip files
 *
 * Usage:
 *   orc-unstream <input_stream> <output_file>
 *   Example: ./orc-unstream ORC_Foo_Quick.7zs ORC_Foo_Quick.7z
 *
 * This utility transforms ORC output files in 'stream/journal' format
 * into standard files using pure Go with no third-party dependencies.
 */

package fastfind

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

/*
 * Constants and structures
 */

// Constants for the file format
const (
	JrnlVersion = 2
	BufferSize  = 1024 * 1024 // 1MB buffer size for read/write operations
)

// Magic constants for operations
var (
	MagicJrnl  = [4]byte{'J', 'R', 'N', 'L'}
	MagicWrite = [4]byte{'W', 'R', 'I', 'T'}
	MagicSeek  = [4]byte{'S', 'E', 'E', 'K'}
	MagicClose = [4]byte{'C', 'L', 'O', 'S'}
)

// JournalHeader represents the header of a journal file
type journalHeader struct {
	Magic   [4]byte
	Version uint32
}

// Operation represents a journal operation
type operation struct {
	Magic  [4]byte
	Param  uint32
	Length uint64
}

// BytesEqual compares two byte arrays for equality
func bytesEqual(a, b [4]byte) bool {
	return a == b
}

// ProcessStream handles the conversion from journal stream to regular 7zip file
func processStream(in io.ReadSeeker, out io.WriteSeeker) error {
	// Read and verify the journal header
	var header journalHeader
	if err := binary.Read(in, binary.LittleEndian, &header); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	if !bytesEqual(header.Magic, MagicJrnl) {
		return fmt.Errorf("invalid file format: expected JRNL magic, got %s", string(header.Magic[:]))
	}

	if header.Version != JrnlVersion {
		return fmt.Errorf("unsupported journal version: got %d, expected %d", header.Version, JrnlVersion)
	}

	log.Tracef("Unstream: Valid JRNL header found, version %d", header.Version)

	// Process operations until CLOSE is encountered
	buffer := make([]byte, BufferSize)
	for {
		var op operation
		if err := binary.Read(in, binary.LittleEndian, &op); err != nil {
			return fmt.Errorf("failed to read operation: %w", err)
		}

		log.Tracef("Unstream: READ operation: %s", string(op.Magic[:]))

		// Process operation based on magic value
		if bytesEqual(op.Magic, MagicClose) {
			log.Trace("Unstream: CLOSE operation - processing complete")
			break
		} else if bytesEqual(op.Magic, MagicSeek) {
			whence := 0 // Default to SEEK_SET (io.SeekStart)
			switch op.Param {
			case 0:
				whence = io.SeekStart
			case 1:
				whence = io.SeekCurrent
			case 2:
				whence = io.SeekEnd
			default:
				return fmt.Errorf("unstream: invalid seek mode: %d", op.Param)
			}

			log.Tracef("Unstream: SEEK operation: offset %d, mode %d", op.Length, op.Param)
			if _, err := out.Seek(int64(op.Length), whence); err != nil {
				return fmt.Errorf("seek operation failed: %w", err)
			}
		} else if bytesEqual(op.Magic, MagicWrite) {
			remaining := op.Length
			log.Tracef("Unstream: WRITE operation: %d bytes", remaining)

			for remaining > 0 {
				bytesToRead := remaining
				if bytesToRead > BufferSize {
					bytesToRead = BufferSize
				}

				bytesRead, err := io.ReadFull(in, buffer[:bytesToRead])
				if err != nil {
					return fmt.Errorf("failed to read data: %w", err)
				}

				bytesWritten, err := out.Write(buffer[:bytesRead])
				if err != nil {
					return fmt.Errorf("failed to write data: %w", err)
				}

				if bytesWritten != bytesRead {
					return fmt.Errorf("write error: wrote %d of %d bytes", bytesWritten, bytesRead)
				}

				remaining -= uint64(bytesRead)
			}
		} else {
			return fmt.Errorf("unknown operation type: %s", string(op.Magic[:]))
		}
	}

	return nil
}

/*
 * Unstream receives a serialized filename and an output filename and converts the serialized file to a regular 7zip file
 */
func Unstream(inputPath string, outputPath string) error {

	var inputFile io.ReadSeeker

	var err error
	inputFile, err = os.Open(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input file '%s': %v\n", inputPath, err)
		return err
	}
	defer inputFile.(io.Closer).Close()
	log.Debugf("Opened input file: %s", inputPath)

	// Open output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file '%s': %v\n", outputPath, err)
		return err
	}
	defer outputFile.Close()
	log.Debugf("Opened output file: %s", outputPath)

	// Process the file
	if err := processStream(inputFile, outputFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error processing stream: %v\n", err)
		return err
	}

	log.Debugf("Successfully converted journal stream '%s' to file '%s.", inputPath, outputPath)

	return err
}

/*
 * Unstream receives serialized in a buffer and the output filename and converts the serialized file to a regular 7zip file
 */
func UnstreamBuffer(inputData []byte, outputPath string) error {

	input := bytes.NewReader(inputData)

	// Open output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file '%s': %v\n", outputPath, err)
		return err
	}
	defer outputFile.Close()
	log.Tracef("Opened output file: %s", outputPath)

	// Process the file
	if err := processStream(input, outputFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error processing stream: %v\n", err)
		return err
	}

	log.Tracef("Successfully converted journal stream  to file '%s.", outputPath)

	return err
}
