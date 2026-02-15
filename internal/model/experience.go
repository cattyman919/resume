package model

import (
	"fmt"
	"strings"
)

type Experience struct {
	Company   string
	Location  string
	JobType   string
	Role      string
	FromDate  string
	AfterDate string
	Points    []string
}

func (exp Experience) String() string {
	var res strings.Builder
	fmt.Fprintf(&res, "Company: %s\n", exp.Company)
	fmt.Fprintf(&res, "Location: %s\n", exp.Location)
	fmt.Fprintf(&res, "Job Type: %s\n", exp.JobType)
	fmt.Fprintf(&res, "Role: %s\n", exp.Role)
	fmt.Fprintf(&res, "From: %s\n", exp.FromDate)
	fmt.Fprintf(&res, "After: %s\n", exp.AfterDate)
	res.WriteString("Points:\n")
	for _, p := range exp.Points {
		fmt.Fprintf(&res, "- %s\n", p)
	}
	return res.String()
}
