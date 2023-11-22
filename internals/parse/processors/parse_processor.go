package processors

import (
	"cian-parse/internals/models"
	"cian-parse/pkg/conv"
	"github.com/gocolly/colly"
)

type ParseProcessor struct {
	collector *colly.Collector
}

func NewParseProcessor() *ParseProcessor {
	processor := new(ParseProcessor)
	processor.collector = colly.NewCollector()
	return processor
}

func (p *ParseProcessor) ParseLinks(link string) ([]string, error) {
	var links []string

	p.collector.OnHTML("div._93444fe79c--wrapper--W0WqH", func(e *colly.HTMLElement) {
		links = e.ChildAttrs("a._93444fe79c--media--9P6wN", "href")
	})

	err := p.collector.Visit(link)
	if err != nil {
		return nil, err
	}

	//TODO: check DB
	//for i := 0; i < len(links); i++ { }

	return links, nil
}

func (p *ParseProcessor) Parse(links []string) ([]models.Immovable, error) {
	var result []models.Immovable
	var isCheckedTitle, isCheckedLink, isCheckedPrice bool

	for index, link := range links {
		result = append(result, models.Immovable{})
		result[index].Link = link

		p.collector.OnHTML("h1.a10a3f92e9--title--vlZwT", func(e *colly.HTMLElement) {
			if !isCheckedTitle {
				result[index].Title = e.Text
				isCheckedTitle = true
			} else {
				return
			}
		})

		p.collector.OnHTML("div.a10a3f92e9--amount--ON6i1 span.a10a3f92e9--color_black_100--Ephi7", func(e *colly.HTMLElement) {
			if !isCheckedLink {
				result[index].Price, _ = conv.StrtoIntWithoutSpace(e.Text)
				isCheckedLink = true
			} else {
				return
			}
		})

		p.collector.OnHTML("table.a10a3f92e9--history--JRbxR", func(e *colly.HTMLElement) {
			if !isCheckedPrice {
				selection := e.DOM
				result[index].PriceInitially, _ = conv.StrtoIntWithoutSpace(selection.Find("td.a10a3f92e9--event-price--xNv2v").Text())
				result[index].Data, _ = conv.StrToStrLastElement(selection.Find("td.a10a3f92e9--event-date--BvijC").Text())
				isCheckedPrice = true
			} else {
				return
			}
		})

		isCheckedTitle = false
		isCheckedLink = false
		isCheckedPrice = false
		err := p.collector.Visit(link)
		if err != nil {
			return nil, err
		}

	}

	return result, nil
}
