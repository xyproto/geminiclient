# Simple Gemini

A simple way to use Gemini.

### Main features

* Supports handing over a prompt and receiving a response.
* Can run both locally and in ie. a Google Cloud Run instance.
* Supports multi-modal prompts (prompts where you can add text, images or data to the prompt).
* Supports tool / function calling where you can supply custom Go functions to the Gemini client, and Gemini can call the functions as needed.

### Example use

1. Run `gcloud auth application-default login`, if needed.
2. Get the Google Project ID at https://console.cloud.google.com/.
3. `export GCP_PROJECT=123`, where "123" is your own Google Project ID.
4. (optionally) `export GCP_LOCATON=us-west1`, where "us-west1" is the location you prefer.
5. Create a directory for this experiment, for instance: `mkdir ~/geminitest` and then `cd ~/geminitest`.
6. Create a `main.go` file that looks like this:


```go
package main

import (
    "fmt"

    "github.com/xyproto/simplegemini"
)

func main() {
    fmt.Println(simplegemini.MustAsk("Write a haiku about cows.", 0.4))
}
```

7. Prepare a `go.mod` file for with `go mod init cows`
8. Fetch the dependencies (this simplegemini package) with `go mod tidy`
9. Build and run the executable: `go build && ./cows`
10. Observe the output, that should look a bit like this:

```
Black and white patches,
Chewing grass in sunlit fields,
Mooing gentle song.
```

### Environment variables

These environment variables are supported:

* `GCP_PROJECT` or `PROJECT_ID` for the Google Cloud Project ID
* `GCP_LOCATION` or `PROJECT_LOCATION` for the Google Cloud Project location (like `us-west1`)
* `MODEL_NAME` for the Gemini model name (like `gemini-1.5-flash` or `gemini-1.5-pro`)
* `MULTI_MODAL_MODEL_NAME` for the Gemini multi-modal name (like `gemini-1.0-pro-vision`)

### General info

* Version: 1.0.0
* License: Apache 2
* Author: Alexander F. Rødseth
