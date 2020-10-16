package process

import (
	"fmt"
	"log"
	"regexp"
)

func FindInfo(text, infoRe string) ([]string, error) {
	matchResult := regexp.MustCompile(infoRe).FindStringSubmatch(text)
	if len(matchResult) == 0 {
		log.Printf("FindStringSubmatch failed matchResult:%v\n", matchResult)
		return nil, fmt.Errorf("FindItemInfo failed regexp:%s", infoRe)
	}
	return matchResult[1:], nil
}

func FindAllInfo(text, allInfoRe string) ([][]string, error) {
	matchResults := regexp.MustCompile(allInfoRe).FindAllStringSubmatch(text, -1)
	if len(matchResults) == 0 {
		log.Printf("FindStringSubmatch failed matchResults:%v\n", matchResults)
		return nil, fmt.Errorf("FindItemInfos failed regexp:%s", allInfoRe)
	}
	var result [][]string
	for _, match := range matchResults {
		result = append(result, match[1:])
	}
	return result, nil
}
