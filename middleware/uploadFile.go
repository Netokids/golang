package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func UploadFile(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read the body of the request
		file, handler, err := r.FormFile("image")
		// Check if there is an error
		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode("Error Retrieving the File")
			return
		}
		// Close the file when the function returns
		defer file.Close()

		// Read the content of the file
		fmt.Printf("Uploaded File: %+v\n", handler.Filename)

		// Create a temporary file within our temp-images directory that follows
		tempFile, err := ioutil.TempFile("uploads", "image-*"+handler.Filename)
		// Check if there is an error
		if err != nil {
			fmt.Println(err)
			fmt.Println("path upload error")
			json.NewEncoder(w).Encode(err)
			return
		}
		// Close the file when the function returns
		defer tempFile.Close()

		// Read all of the contents of our uploaded file into a
		fileBytes, err := ioutil.ReadAll(file)
		// Check if there is an error
		if err != nil {
			fmt.Println(err)
		}

		// Write this byte array to our temporary file
		tempFile.Write(fileBytes)
		// return that we have successfully uploaded our file!
		data := tempFile.Name()
		filename := data[8:]

		// Set the filename in the context
		ctx := context.WithValue(r.Context(), "dataFile", filename)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
