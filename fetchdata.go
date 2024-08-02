package main

import (
    "encoding/csv"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "os"
    "regexp"
    "strings"
    "golang.org/x/net/html"
)

func fetchAndParseData(baseURL, downloadURL string, headers map[string]string) ([][]string, error) {
    client := &http.Client{}

    // Fetch the base URL page
    resp, err := client.Get(baseURL)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    pageContent, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    // Extract standard IDs from HTML
    standardIDs := extractStandardIDs(pageContent)

    // Post to download URL
    jsonData, err := json.Marshal(standardIDs)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequest("POST", downloadURL, strings.NewReader(string(jsonData)))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")
    for key, value := range headers {
        req.Header.Set(key, value)
    }

    resp, err = client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    csvContent, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    // Read CSV content
    reader := csv.NewReader(strings.NewReader(string(csvContent)))
    reader.Comma = ',' // assuming CSV delimiter is comma
    reader.LazyQuotes = true
    records, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }

    return records, nil
}

func extractStandardIDs(pageContent []byte) []string {
    var standardIDs []string

    tokenizer := html.NewTokenizer(strings.NewReader(string(pageContent)))
    for {
        tokenType := tokenizer.Next()
        switch tokenType {
        case html.ErrorToken:
            return standardIDs
        case html.StartTagToken, html.SelfClosingTagToken:
            token := tokenizer.Token()
            if token.Data == "div" {
                for _, attr := range token.Attr {
                    if attr.Key == "data-standard" {
                        var data map[string]interface{}
                        if err := json.Unmarshal([]byte(attr.Val), &data); err == nil {
                            if id, ok := data["id"].(string); ok {
                                standardIDs = append(standardIDs, id)
                            }
                        }
                    }
                }
            }
        }
    }
}

func main() {
    baseURL := "https://www.instituteforapprenticeships.org/apprenticeship-standards/"
    downloadURL := baseURL + "download"
    headers := map[string]string{
        "Content-Type": "application/json",
        "Accept":       "application/json",
    }

    data, err := fetchAndParseData(baseURL, downloadURL, headers)
    if err != nil {
        fmt.Println("Error fetching and parsing data:", err)
        return
    }

    // Save data to CSV file
    file, err := os.Create("apprenticeship_data.csv")
    if err != nil {
        fmt.Println("Error creating CSV file:", err)
        return
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    for _, record := range data {
        if err := writer.Write(record); err != nil {
            fmt.Println("Error writing record to CSV file:", err)
        }
    }
}
