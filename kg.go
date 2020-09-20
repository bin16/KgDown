package main

import (
	"fmt"
)

var nameMap map[string]string

type kgShare struct {
	Title   string
	ImgURL  string `json:"img_url"`
	DataURL string `json:"data_url"`
}

type kgDetail struct {
	Cover    string
	PlayURL  string
	SongName string `json:"song_name"`
	KgNick   string `json:"hg_nick"`
	HcNick   string `json:"hc_nick"`
}

type kgData struct {
	ShareID   string
	Share     kgShare
	Album     []string
	ShareNick string
	Detail    kgDetail
}

func (kg kgData) songURL() string {
	if kg.Detail.PlayURL != "" {
		return kg.Detail.PlayURL
	}

	return kg.Share.DataURL
}

func (kg kgData) albumURL() string {
	if len(kg.Album) > 0 {
		return kg.Album[0]
	}

	return kg.Share.ImgURL
}

func (kg kgData) songTitle() string {
	return kg.Detail.SongName
}

func (kg kgData) artist() string {
	k := kg.Detail.KgNick
	k2 := kg.Detail.HcNick
	if k == "" {
		return nameify(kg.ShareNick)
	} else if k2 == "" || k2 == k {
		return nameify(k)
	}

	return fmt.Sprintf("%s, %s", nameify(k), nameify(k2))
}

func nameify(s string) string {
	if v, ok := nameMap[s]; ok {
		return v
	}

	return s
}
