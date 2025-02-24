package business

import (
	"crave/miner/cmd/model"
	craveModel "crave/shared/model"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type NamuBusiness struct {
}

func NewNamuBusiness() *NamuBusiness {
	return &NamuBusiness{}
}

func (biz *NamuBusiness) MakeUrl(step craveModel.Step, name string) string {
	baseUrl := "https://namu.wiki/"
	frontlink := "w/"
	if craveModel.Front == step {
		return baseUrl + frontlink + name
	}
	backlink := "backlink/"
	if craveModel.Back == step {
		return baseUrl + backlink + name
	}
	return ""
}

func (biz *NamuBusiness) ParseNextTargets(step craveModel.Step, name string) ([]model.ParsedTarget, error) {
	url := biz.MakeUrl(step, name)
	doc, err := biz.GetDocument(url)
	if err != nil {
		return nil, err
	}
	return biz.ExtractTargets(doc)
}

func (biz *NamuBusiness) GetDocument(url string) (*goquery.Document, error) {
	client := &http.Client{}

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("🛑error creating request: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("🛑error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Parse the HTML content
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %w", err)
	}

	return doc, nil

}

func (biz *NamuBusiness) ExtractTargets(doc *goquery.Document) ([]model.ParsedTarget, error) {
	targetMap := make(map[string]*model.ParsedTarget)
	doc.Find("a[href^='/w/']").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		input := strings.TrimPrefix(href, "/w/")

		decodedString, err := url.QueryUnescape(input)
		if err != nil {
			fmt.Printf("Error decoding URL: %v\n", err)
			return
		}

		indexOfHash := strings.Index(decodedString, "#")
		result := decodedString
		if indexOfHash != -1 {
			result = decodedString[:indexOfHash]
		}
		if _, exists := targetMap[result]; !exists {
			context := biz.ExtractContext(s, result)
			targetMap[result] = &model.ParsedTarget{
				Name:       result,
				Context:    context,
				Appearance: 1,
			}
			return
		}
		targetMap[result].Appearance++
	})
	targets := make([]model.ParsedTarget, 0, len(targetMap))
	for _, target := range targetMap {
		targets = append(targets, *target)
	}
	sort.Slice(targets, func(i, j int) bool {
		return targets[i].Appearance > targets[j].Appearance
	})
	return targets, nil
}

func (biz *NamuBusiness) ExtractContext(s *goquery.Selection, name string) string {

	fullText := s.Parent().Text()
	if len(fullText) <= len(name) {
		return fullText
	}

	linkIndex := strings.Index(fullText, name)
	if linkIndex == -1 {
		return fullText
	}

	sentences := strings.Split(fullText, ".")
	for _, sentence := range sentences {
		if strings.Contains(sentence, name) {
			return sentence
		}
	}
	return sentences[0]
}
