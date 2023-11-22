package parse

import (
	"cian-parse/internals/models"
	"cian-parse/internals/parse/processors"
	"cian-parse/pkg/random"
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"
)

// TODO: create a url for all requests
const mainUrl = "https://krasnodar.cian.ru/kupit-kvartiru/"
const testUrl = "https://krasnodar.cian.ru/kupit-kvartiru-krasnodarskiy-kray-krasnodar-komsomolskiy-04577/"

type Parser struct {
	log       *slog.Logger
	wg        sync.WaitGroup
	processor *processors.ParseProcessor
}

func NewParser(logger *slog.Logger) *Parser {
	parser := new(Parser)
	parser.log = logger
	parser.wg = sync.WaitGroup{}
	parser.processor = processors.NewParseProcessor()
	return parser
}

func (p *Parser) ParseImmovable() error {
	var result [][]models.Immovable
	timer := time.Now()

	p.log.Info("parsing all data")

	p.parseAll(&result)

	p.wg.Wait()
	fmt.Println(result)
	p.log.Info("parsing ended")

	p.log.Info("time after start parsing", time.Since(timer).String())

	return nil
}

func (p *Parser) parseFirstPage(result *[][]models.Immovable) {
	go func() {
		p.wg.Add(1)
		//TODO: remake url
		links, err := p.processor.ParseLinks(mainUrl)
		if err != nil {
			p.log.Error("error to parse in func parseFirstPage", err)
		}

		res, err := p.processor.Parse(links)
		if err != nil {
			p.log.Error("error to parse in func parseFirstPage", err)
		}
		*result = append(*result, res)

		//TODO: remake, too many
		p.log.Info("1" + " page has been parsed with " + strconv.Itoa(len(links)) + " cards")
		p.wg.Done()
	}()

	time.Sleep(random.RandomSec(4, 6) * time.Second)
}

func (p *Parser) parseAll(result *[][]models.Immovable) {
	p.parseFirstPage(result)

	done := make(chan struct{})
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for page := 2; page <= 55; page++ { //in real 55
		//p.wg.Add(1)
		func(page int) {
			//TODO: create a url for all requests
			//https://krasnodar.cian.ru/cat.php?deal_type=sale&engine_version=2&offer_type=flat&p=
			//https://krasnodar.cian.ru/cat.php?deal_type=sale&district%5B0%5D=577&engine_version=2&offer_type=flat&p=
			links, err := p.processor.ParseLinks(fmt.Sprintf("%s%d%s", "https://krasnodar.cian.ru/cat.php?deal_type=sale&engine_version=2&offer_type=flat&p=", page, "&region=4820"))
			if err != nil {
				p.log.Error("error to parse in func parseAll", err)
			}

			res, err := p.processor.Parse(links)
			if err != nil {
				p.log.Error("error to parse in func parseAll", err)
			}

			select {
			case <-done:
			default:
				*result = append(*result, res)
				//TODO: remake, too many
				p.log.Info(strconv.Itoa(page) + " page has been parsed with " + strconv.Itoa(len(res)) + " cards")
			}

			if len(res) < 28 && len(res) != 0 {
				cancel()
				done <- struct{}{}
			}

			//p.wg.Done()
		}(page)

		select {
		case <-ctx.Done():
			return
		default:
			//time.Sleep(random.RandomSec(4, 6) * time.Second)
		}
	}
}
