package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

// item represents one scraped book/product from the website.
// The struct tags (json:"...") control how the fields will appear
// in the generated JSON file.
type item struct {
	DisplayTitle string `json:"display_title"` // Visible short title from the page
	Price        string `json:"price"`         // Product price
	ImgURL       string `json:"img_url"`       // Image URL/path
	FullTitle    string `json:"full_title"`    // Full title stored in the title attribute
}

func main() {
	// Create a new Colly collector.
	// AllowedDomains limits requests so the scraper only visits this domain.
	c := colly.NewCollector(
		colly.AllowedDomains("books.toscrape.com"),
	)

	// items will store all scraped results before writing them to JSON.
	var items []item

	// Register a callback for every product card on the page.
	// Each matching <article class="product_pod"> element represents one book.
	c.OnHTML("article[class=product_pod]", func(e *colly.HTMLElement) {
		// Build one item by extracting text/attributes from the current product element.
		item := item{
			DisplayTitle: e.ChildText("h3 a"),                     // Text inside the book link
			Price:        e.ChildText("div p[class=price_color]"), // Price text
			ImgURL:       e.ChildAttr("div a img", "src"),         // src attribute of the image
			FullTitle:    e.ChildAttr("h3 a", "title"),            // Full title from the title attribute
		}

		// Add the scraped item to the slice.
		items = append(items, item)
	})

	// Register a callback for the "next" pagination button.
	// This lets the scraper continue through all pages automatically.
	c.OnHTML("[class=next]", func(h *colly.HTMLElement) {
		// Get the relative URL of the next page.
		next_page := h.ChildAttr("a", "href")

		// Convert the relative URL to an absolute URL and visit it.
		c.Visit(h.Request.AbsoluteURL(next_page))
	})

	// This runs before every request.
	// Useful for seeing which pages the scraper is visiting.
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping from the homepage.
	c.Visit("https://books.toscrape.com/")

	// Convert the items slice into JSON format.
	content, err := json.Marshal(items)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Save the JSON output into a file named items.json.
	err = os.WriteFile("items.json", content, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
}
