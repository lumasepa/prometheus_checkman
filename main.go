package main

import "log"

func main() {

	conf, err := ParseConf("./checkman.yml")

	if err != nil {
		log.Fatal(err)
	}

	resultsChan := make(chan CheckResult, 4096)

	go mainExporter(conf, resultsChan)
	mainScheduler(conf.Checks, resultsChan)
}


