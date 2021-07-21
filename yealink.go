package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

type YealinkRoot struct {
	XMLName          xml.Name        `xml:"bisdIPPhoneDirectory"`
	Clearlight       string          `xml:"clearlight,attr"`
	Title            string          `xml:"Title"`
	Prompt           string          `xml:"Prompt"`
	DirectoryEntries []*YealinkEntry `xml:"DirectoryEntry"`
}

type YealinkEntry struct {
	Name      string `xml:"Name"`
	Telephone string `xml:"Telephone"`
}

func YealinkXML(entries []*Entry) ([]byte, error) {
	yeaEntries := make([]*YealinkEntry, 0, len(entries))
	for _, e := range entries {
		yeaEntries = append(yeaEntries, &YealinkEntry{Name: e.User, Telephone: e.Extension})
	}

	root := &YealinkRoot{Clearlight: "true", Title: "Phonelist", Prompt: "Prompt", DirectoryEntries: yeaEntries}

	header := []byte(xml.Header)
	data, err := xml.MarshalIndent(root, "", "\t")
	if err != nil {
		return nil, fmt.Errorf("could not marshal XML: %w", err)
	}

	return append(header, data...), nil
}

func (c *Config) YealinkHandler(w http.ResponseWriter, r *http.Request) {
	entries, err := GetEntries(c.SQLDSN)
	if err != nil {
		log.Println("ERROR: could not get entries:", err)
		WriteError(w, http.StatusInternalServerError)
		return
	}

	buf, err := YealinkXML(entries)
	if err != nil {
		log.Println("ERROR: could not generate XML:", err)
		WriteError(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}
