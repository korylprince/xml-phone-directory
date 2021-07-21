package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

type GrandstreamRoot struct {
	XMLName  xml.Name            `xml:"AddressBook"`
	Contacts []*GrandstreamEntry `xml:"Contact"`
}

type GrandstreamEntry struct {
	FirstName    string `xml:"FirstName"`
	LastName     string `xml:"LastName"`
	PhoneNumber  string `xml:"phonenumber"`
	AccountIndex int    `xml:"accountindex"`
}

func GrandstreamXML(entries []*Entry) ([]byte, error) {
	gEntries := make([]*GrandstreamEntry, 0, len(entries))
	for _, e := range entries {
		first, last := SplitName(e.User)
		gEntries = append(gEntries, &GrandstreamEntry{FirstName: first, LastName: last, PhoneNumber: e.Extension, AccountIndex: 1})
	}

	root := &GrandstreamRoot{Contacts: gEntries}

	header := []byte(xml.Header)
	data, err := xml.MarshalIndent(root, "", "\t")
	if err != nil {
		return nil, fmt.Errorf("could not marshal XML: %w", err)
	}

	return append(header, data...), nil
}

func (c *Config) GrandstreamHandler(w http.ResponseWriter, r *http.Request) {
	entries, err := c.cache.GetEntries(c.SQLDSN)
	if err != nil {
		log.Println("ERROR: could not get entries:", err)
		WriteError(w, http.StatusInternalServerError)
		return
	}

	buf, err := GrandstreamXML(entries)
	if err != nil {
		log.Println("ERROR: could not generate XML:", err)
		WriteError(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}
