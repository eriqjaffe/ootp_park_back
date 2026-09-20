package main

import (
	"bytes"
	"fmt"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/webp"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Converting in: %s\n", root)

	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))

		switch ext {
		case ".webp":
			if err := convertAndDelete(path); err != nil {
				fmt.Printf("[-] Image Error on %s: %v\n", filepath.Base(path), err)
			}
		case ".mtl":
			if err := updateMtlFile(path); err != nil {
				fmt.Printf("[-] MTL Error on %s: %v\n", filepath.Base(path), err)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error during folder traversal: %v\n", err)
	} else {
		fmt.Println("\nAll operations completed successfully!")
	}
}

func convertAndDelete(webpPath string) error {
	inFile, err := os.Open(webpPath)
	if err != nil {
		return fmt.Errorf("ERROR: could not open file: %w", err)
	}

	img, err := webp.Decode(inFile)
	if err != nil {
		inFile.Close()
		return fmt.Errorf("ERROR: could not decode webp: %w", err)
	}
	inFile.Close()

	pngPath := webpPath[:len(webpPath)-len(filepath.Ext(webpPath))] + ".png"
	outFile, err := os.Create(pngPath)
	if err != nil {
		return fmt.Errorf("ERROR: could not create png: %w", err)
	}

	err = png.Encode(outFile, img)
	if err != nil {
		outFile.Close()
		return fmt.Errorf("ERROR: could not encode png: %w", err)
	}
	outFile.Close()

	fmt.Printf("[+] Converted: %s\n", filepath.Base(webpPath))

	err = os.Remove(webpPath)
	if err != nil {
		return fmt.Errorf("ERROR: could not delete .webp file: %w", err)
	}

	return nil
}

func updateMtlFile(mtlPath string) error {
	content, err := os.ReadFile(mtlPath)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	hasWebpLower := bytes.Contains(content, []byte(".webp"))
	hasWebpHigher := bytes.Contains(content, []byte(".WEBP"))

	if !hasWebpLower && !hasWebpHigher {
		return nil
	}

	updatedContent := bytes.ReplaceAll(content, []byte(".webp"), []byte(".png"))
	updatedContent = bytes.ReplaceAll(updatedContent, []byte(".WEBP"), []byte(".PNG"))

	err = os.WriteFile(mtlPath, updatedContent, 0644)
	if err != nil {
		return fmt.Errorf("could not save updated file: %w", err)
	}

	fmt.Printf("[*] Updated textures inside %s\n", filepath.Base(mtlPath))
	return nil
}
