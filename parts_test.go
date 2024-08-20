package simplegemini_test

import (
	"os"
	"testing"

	"github.com/xyproto/simplegemini"
)

func TestAddImage(t *testing.T) {
	gc := simplegemini.MustNew()

	// Create a temporary image file for testing
	tmpfile, err := os.CreateTemp("", "testimage.png")
	if err != nil {
		t.Fatal("Failed to create temporary image file:", err)
	}
	defer os.Remove(tmpfile.Name())

	// Write some dummy data to the image file
	if _, err := tmpfile.Write([]byte("PNG DATA")); err != nil {
		t.Fatal("Failed to write to temporary image file:", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal("Failed to close temporary image file:", err)
	}

	// Test adding the image
	err = gc.AddImage(tmpfile.Name())
	if err != nil {
		t.Fatalf("Expected to successfully add image, but got error: %v", err)
	}

	// Check if the part was added
	if len(gc.Parts) != 1 {
		t.Fatalf("Expected 1 part to be added, but got %d", len(gc.Parts))
	}
}

func TestMustAddImageInvalidPath(t *testing.T) {
	gc := simplegemini.MustNew()

	// Attempt to add an image from an invalid path
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic when adding image from invalid path, but no panic occurred")
		}
	}()
	gc.MustAddImage("/non/existent/path.png")
}

func TestAddURI(t *testing.T) {
	gc := simplegemini.MustNew()

	gc.AddURI("gs://generativeai-downloads/images/scones.jpg")

	// Check if the part was added
	if len(gc.Parts) != 1 {
		t.Fatalf("Expected 1 part to be added, but got %d", len(gc.Parts))
	}
}

func TestAddData(t *testing.T) {
	gc := simplegemini.MustNew()

	// Test adding data
	data := []byte("Some data")
	gc.AddData("text/plain", data)

	// Check if the part was added
	if len(gc.Parts) != 1 {
		t.Fatalf("Expected 1 part to be added, but got %d", len(gc.Parts))
	}
}

func TestAddText(t *testing.T) {
	gc := simplegemini.MustNew()

	// Test adding text
	gc.AddText("This is a prompt")

	// Check if the part was added
	if len(gc.Parts) != 1 {
		t.Fatalf("Expected 1 part to be added, but got %d", len(gc.Parts))
	}
}

func TestClearParts(t *testing.T) {
	gc := simplegemini.MustNew()

	// Add some parts
	gc.AddText("Text part")
	gc.AddData("text/plain", []byte("Some data"))

	// Clear the parts
	gc.ClearParts()

	// Check if parts were cleared
	if len(gc.Parts) != 0 {
		t.Fatalf("Expected 0 parts after clearing, but got %d", len(gc.Parts))
	}
}
