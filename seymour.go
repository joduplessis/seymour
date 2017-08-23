package main

import (
  "bufio"
  "fmt"
  "strings"
  "log"
  "os"
  "io/ioutil"
  "reflect"
  "regexp"
)

func main() {
  // Check if the user provided arguments
  if len(os.Args) == 1 {
    fmt.Println("Sorry, you provided no arguments.")

    os.Exit(3)
  }

  // Default to no HTML
  html := false
  arguments := os.Args[1:]
  filename := arguments[0]

  // See if it's a file
  if len(arguments) < 1 {
    fmt.Println("Sorry, you need to provide a file.")

    os.Exit(3)
  }

  // See what the extension of the file is
  extension := strings.Split(arguments[0], ".")[1]

  // See if it's a file
  if extension!="html" && extension!="seymour" {
    fmt.Println("Please provide a Seymour file or a HTML file.")

    os.Exit(3)
  }

  if extension=="html" {
    html = true
  }

  // Open the file
  file, err := os.Open(filename)
  inc := 0
  output := ""
  cssIndex := []string{}
  cssClasses := []string{}

  // Start the output
  if html {
    fmt.Print("Generating from HTML...")
  } else {
    fmt.Print("Generating from Seymour...")
  }

  // If there is a fatal error
  if err != nil {
    log.Fatal(err)
  }

  // Close the file
  defer file.Close()

  // Scan the buffer
  scanner := bufio.NewScanner(file)

  // Dependant on file
  if html {
    fmt.Println("WARNING")
    fmt.Println("The HTML output is experimental...")

    // Iterate over the contents
    for scanner.Scan() {
      line := scanner.Text()
      trimmedLine := strings.Trim(line, " ")

      if trimmedLine != "" {
        spacing := ""
        cssClass := ""
        regexHtml := regexp.MustCompile(`<(\S*)`)
        regexClass := regexp.MustCompile(`class="?([^"]*)"?`)

        // Find our HTML/CSS data
        matchClass := regexClass.FindStringSubmatch(trimmedLine)
        matchHtml := regexHtml.FindStringSubmatch(trimmedLine)

        // Strip out any lingering tags on the HTML search
        htmlTagElements := strings.Split(matchHtml[1], ">")
        htmlElement := htmlTagElements[0]

        // If there is a CSS class - if not, use the tag
        if len(matchClass) > 0 {
          cssClass = htmlElement + "." + matchClass[1]
        } else {
          cssClass = htmlElement
        }

        // Increase the indentation & add the css class to the index
        if string(htmlElement[0]) != "/" {
          inc++

          if cssClass != "" {
            cssIndex = append(cssIndex, cssClass)
          }
        }

        // Adjust our indentation
        for i := 0; i < inc; i++ {
          spacing += "  "
        }

        // Decrease the indentation & remove the last css index
        if string(htmlElement[0]) == "/" {
          inc--

          // Remove last element
          if cssClass != "" {
            cssIndex = cssIndex[:len(cssIndex)-1]
          }
        }

        // Class for this iteratoin
        joinedCssIndex := strings.Join(cssIndex, " ")
        cssIndexString := ""

        // If it's not empty
        if string(joinedCssIndex) != "" {
          cssIndexString = joinedCssIndex + " {}"
        }

        // If it's already in the array, do not add it!
        if !InArray(cssIndexString, cssClasses) {
          cssClasses = append(cssClasses, cssIndexString)
        }
      }
    }

    // Create the byte objects
    outputCss := []byte(strings.Join(cssClasses, "\n"))

    // Write out the bytes for the files
    ioutil.WriteFile("seymour.css", outputCss, 0644)
  } else {
    // Iterate over the contents
    for scanner.Scan() {
      text := strings.Trim(scanner.Text(), " ")
      parts := strings.Split(text, " ")

      // If it's not empty
      if string(parts[0]) != "" {
        tag := parts[0]
        cssClass := parts[0]
        attributeParts := append(parts[:0], parts[1:]...)
        attributes := ""
        spacing := ""

        // If there attributes & it's not plain text
        if len(attributeParts) != 0 && string(tag[0])!=">" {
          for attribute := range attributeParts {
            singleAtt := attributeParts[attribute]
            singleAttSlice := strings.Split(singleAtt, ":")
            attributes = attributes + " " + singleAttSlice[0] + "=\"" + singleAttSlice[1] + "\""

            if (singleAttSlice[0]=="class") {
              cssClass = cssClass + "." + singleAttSlice[1]
            }
          }
        }

        // If there are attributes
        if (attributes!="") {
          tag = tag + attributes
        }

        // Increase the indentation & add the css class
        if string(tag[0])!="/" {
          inc++

          cssIndex = append(cssIndex, cssClass)
        }

        // Adjust our indentation
        for i := 0; i < inc; i++ {
          spacing += "  "
        }

        // Create the output for the HTML: important that this is in between inc/dec
        // Notice how we're using "text" here, not "tag": it's because we want to display
        // The whole line with indentation
        if string(cssClass[0])==">" {
          output += spacing + text[1:len(text)] + "\n"
        } else {
          output += spacing + "<" + tag + ">" + "\n"
        }

        // Decrease the indentation & remove the last css index
        if string(cssClass[0])=="/" || string(cssClass[0])==">" {
          inc--

          // Remove last element
          cssIndex = cssIndex[:len(cssIndex)-1]
        }

        // Class for this iteratoin
        joinedCssIndex := strings.Join(cssIndex, " ")
        cssIndexString := ""

        // If it's not empty
        if (string(joinedCssIndex)!="") {
          cssIndexString = joinedCssIndex + " {}"
        }

        // If it's already in the array
        if !InArray(cssIndexString, cssClasses) {
          cssClasses = append(cssClasses, cssIndexString)
        }
      }
    }

    // Create the byte objects
    outputHtml := []byte(output)
    outputCss := []byte(strings.Join(cssClasses, "\n"))

    // Write out the bytes for the files
    ioutil.WriteFile("seymour.html", outputHtml, 0644)
    ioutil.WriteFile("seymour.css", outputCss, 0644)
  }

  // Log OKAY
  fmt.Println("Done!")

  // If there's an error
  if err := scanner.Err(); err != nil {
    fmt.Println("Nope!\n")

    log.Fatal(err)
  }
}

// Checks to see if it's in the array
func InArray(val interface{}, array interface{}) (exists bool) {
    exists = false

    switch reflect.TypeOf(array).Kind() {
      case reflect.Slice:
        s := reflect.ValueOf(array)

        for i := 0; i < s.Len(); i++ {
          if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
            exists = true

            return
          }
        }
    }

    return
}
