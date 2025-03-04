package business

import (
	filterBusiness "crave/miner/cmd/api/domain/business/filter"
	page "crave/miner/cmd/api/infrastructure/externalApi/page"
	"crave/miner/cmd/model"
	craveModel "crave/shared/model"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type NamuBusiness struct {
	clnt page.IClient
}

func NewNamuBusiness(clnt page.IClient) *NamuBusiness {
	return &NamuBusiness{clnt: clnt}
}

func (biz *NamuBusiness) MakeFrontUrl(name string) string {
	return "https://namu.wiki/w/" + name
}

func (biz *NamuBusiness) MakeBackUrl(name string) string {
	return "https://namu.wiki/backlink/" + name
}

func (biz *NamuBusiness) ParseNextTargets(step craveModel.Step, filter filterBusiness.IBusiness, name string) ([]model.ParsedTarget, error) {

	if craveModel.Front == step {
		url := biz.MakeFrontUrl(name)
		html, err := biz.clnt.GetHtml(url)
		if err != nil {
			return nil, err
		}
		doc, err := biz.GetDocument(html)
		if err != nil {
			return nil, err
		}
		return biz.ExtractFrontTargets(doc, filter, name)
	}
	if craveModel.Back == step {
		url := biz.MakeBackUrl(name)
		html, err := biz.clnt.GetHtml(url)
		if err != nil {
			return nil, err
		}
		doc, err := biz.GetDocument(html)
		if err != nil {
			return nil, err
		}
		return biz.ExtractBackTargets(doc, filter, name)
	}
	return nil, nil
}

func (biz *NamuBusiness) GetDocument(html *string) (*goquery.Document, error) {
	// Parse the HTML content
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(*html))
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %w", err)
	}

	return doc, nil

}

func (biz *NamuBusiness) ExtractFrontTargets(doc *goquery.Document, filter filterBusiness.IBusiness, name string) ([]model.ParsedTarget, error) {
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
		if decodedString == name {
			return
		}

		indexOfHash := strings.Index(decodedString, "#")
		result := decodedString
		if indexOfHash != -1 {
			result = decodedString[:indexOfHash]
		}
		if filter.Apply(&result) != 0 {
			return
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

func (biz *NamuBusiness) ExtractBackTargets(doc *goquery.Document, filter filterBusiness.IBusiness, name string) ([]model.ParsedTarget, error) {
	targetMap := make(map[string]*model.ParsedTarget)
	firstPage := true
	for {
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
			if decodedString == name {
				return
			}

			indexOfHash := strings.Index(decodedString, "#")
			result := decodedString
			if indexOfHash != -1 {
				result = decodedString[:indexOfHash]
			}

			if filter.Apply(&result) != 0 {
				return
			}

			if _, exists := targetMap[result]; !exists {
				targetMap[result] = &model.ParsedTarget{
					Name:       result,
					Context:    "",
					Appearance: 1,
				}
				return
			}
		})

		nextLink := doc.Find("a[href^='/backlink/']")
		if nextLink.Length() == 4 || firstPage {
			nextLink = nextLink.Last()
			firstPage = false
			href, exists := nextLink.Attr("href")
			if exists {
				input := strings.TrimPrefix(href, "/backlink/")
				if strings.Contains(input, "?from=") {
					url := biz.MakeBackUrl(input)
					html, err := biz.clnt.GetHtml(url)
					if err != nil {
						break
					}
					doc, _ = biz.GetDocument(html)
				}
				continue
			}
			break
		}
		break
	}
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

func (biz *NamuBusiness) ApplyFilter(name string, filterChain filterBusiness.FilterChain) (craveModel.Filter, error) {
	url := biz.MakeFrontUrl(name)

	html, err := biz.clnt.GetHtml(url)
	if err != nil {
		return -1, err
	}
	return filterChain.Apply(html), nil
}
