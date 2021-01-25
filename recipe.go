package main

import (
	"encoding/json"
	"os"

	"github.com/jmhobbs/roasted/sr700"
)

type OpenRoastRecipe struct {
	Creator   string `json:"creator"`
	RoastName string `json:"roastName"`
	Steps     []struct {
		TargetTemp  int  `json:"targetTemp,omitempty"`
		FanSpeed    int  `json:"fanSpeed"`
		SectionTime int  `json:"sectionTime"`
		Cooling     bool `json:"cooling,omitempty"`
	} `json:"steps"`
	Bean struct {
		Region string `json:"region"`
		Source struct {
			Reseller string `json:"reseller"`
			Link     string `json:"link"`
		} `json:"source"`
		Country string `json:"country"`
	} `json:"bean"`
	TotalTime        int `json:"totalTime"`
	RoastDescription struct {
		RoastType   string `json:"roastType"`
		Description string `json:"description"`
	} `json:"roastDescription"`
}

type SimpleRecipe struct {
	Steps []struct {
		Heat     sr700.Heat  `json:"heat"`
		Fan      sr700.Speed `json:"fan"`
		Duration int         `json:"duration"`
		Cooling  bool        `json:"cooling,omitempty"`
	} `json:"steps"`
}

func LoadSimpleRecipe(path string) (*SimpleRecipe, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var recipe SimpleRecipe
	err = json.NewDecoder(f).Decode(&recipe)
	return &recipe, err
}
