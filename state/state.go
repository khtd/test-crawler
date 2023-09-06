package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"test-crawler/logging"
)

type Link struct {
	Taken      bool     `json:"taken"`
	Processed  bool     `json:"processed"`
	Discovered []string `json:"discovered"`
	Err        string   `json:"err,omitempty"`
}

type State struct {
	sync.Mutex
	logger     *logging.Logger
	knownLinks map[string]*Link
}

func NewState(logger *logging.Logger) *State {
	return &State{
		logger:     logger,
		knownLinks: make(map[string]*Link),
	}
}

func (s *State) Load() {
	s.logger.Info("start loading previous data")
	if _, err := os.Stat("db.json"); errors.Is(err, os.ErrNotExist) {
		err := os.WriteFile("db.json", []byte("{}"), 0644)
		if err != nil {
			s.logger.Error(err.Error())
		}
	}

	jsonFile, err := os.Open("db.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	v, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(v, &s.knownLinks)

	for key := range s.knownLinks {
		s.knownLinks[key].Taken = false
	}

	s.logger.Info("loaded - " + s.Print())
}

func (s *State) Save() map[string]*Link {
	jsonStr, err := json.Marshal(s.knownLinks)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		err := os.WriteFile("db.json", jsonStr, 0644)
		if err != nil {
			panic(err)
		}
	}
	return s.knownLinks
}

func (s *State) Print() string {
	jsonStr, err := json.Marshal(s.knownLinks)
	if err != nil {
		return fmt.Sprintf("Error: %s", err.Error())
	} else {
		return string(jsonStr)
	}
}

func (s *State) GetUnprocessed(startingLink string) []string {
	s.Lock()
	defer s.Unlock()

	unprocessed := make([]string, 0)
	s.collectUnprocessed(startingLink, func(link string) {
		unprocessed = append(unprocessed, link)
	})

	return unprocessed
}

func (s *State) collectUnprocessed(link string, mark func(string)) {
	data, known := s.knownLinks[link]
	if !known || !data.Processed {
		mark(link)
		return
	}
	for _, nestedLink := range data.Discovered {
		s.collectUnprocessed(nestedLink, mark)
	}
}

func (s *State) TryToTake(link string) bool {
	s.Lock()
	defer s.Unlock()
	data, known := s.knownLinks[link]
	if known && (data.Processed || data.Taken) {
		return false
	}
	s.knownLinks[link] = &Link{Taken: true}
	return true
}

func (s *State) MarkProcessed(link string, discovered []string, err string) {
	s.Lock()
	defer s.Unlock()
	data, known := s.knownLinks[link]
	if !known || !data.Taken || data.Processed {
		s.logger.Error(fmt.Sprintf("try to markProcessed when known: %t, taken: %t, processed: %t", known, data.Taken, data.Processed))
	}
	if err != "" {
		s.knownLinks[link].Err = err
	}
	s.knownLinks[link].Discovered = discovered
	s.knownLinks[link].Processed = true
}

func (s *State) IsKnown(link string) bool {
	s.Lock()
	defer s.Unlock()
	_, known := s.knownLinks[link]
	return known
}
