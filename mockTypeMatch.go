package mockaroo

import (
	"fmt"
	"regexp"
	"strings"
)

var mockRegex map[string]MockType = map[string]MockType{
	"age": MockType{
		Type: "Number",
		Min:  1,
		Max:  100,
	},
	"app[-_]*bundle[-_]*id": MockType{
		Type: "App Bundle ID",
	},
	"app[-_]*name": MockType{
		Type: "App Name",
	},
	"app[-_]*version": MockType{
		Type: "App Version",
	},
	"bitcoin[-_]*address": MockType{
		Type: "Bitcoin Address",
	},
	"buzzword": MockType{
		Type: "Buzzword",
	},
	"catch[-_]*phrase": MockType{
		Type: "Catch Phrase",
	},
	"city": MockType{
		Type: "City",
	},
	"color": MockType{
		Type: "Color",
	},
	"company[-_]*name": MockType{
		Type: "Company Name",
	},
	"country[-_]*code": MockType{
		Type: "Contry Code",
	},
	"credit[-_]*card[-_]*(#){0,1}(number){0,1}": MockType{
		Type: "Credit Card #",
	},
	"credit[-_]*card[-_]*type": MockType{
		Type: "Credit Card Type",
	},
	"currency": MockType{
		Type: "Currency",
	},
	"currency[-_]*code": MockType{
		Type: "Currency Code",
	},
	"domain[-_]*name": MockType{
		Type: "Domain Name",
	},
	"drug[-_]*company": MockType{
		Type: "Drug Company",
	},
	"drug[-_]*name": MockType{
		Type: "Drug Name (brand)",
	},
	"email[-_]*(address){0,1}": MockType{
		Type: "Email Address",
	},
	"file[-_]*name": MockType{
		Type: "File Name",
	},
	"first[-_]*name": MockType{
		Type: "First Name",
	},
	"frequency": MockType{
		Type: "Frequency",
	},
	"full[-_]*name": MockType{
		Type: "Full Name",
	},
	"gender": MockType{
		Type: "Gender",
	},
	"ip[-_]*address[-_]*(v4){0,1}": MockType{
		Type: "IP Addredd v4",
	},
	"ip[-_]*address[-_]*v6": MockType{
		Type: "IP Address v6",
	},
	"isbn": MockType{
		Type: "ISBN",
	},
	"job[-_]*title": MockType{
		Type: "Job Title",
	},
	"language": MockType{
		Type: "Language",
	},
	"last[-_]*name": MockType{
		Type: "Last Name",
	},
	"latitude": MockType{
		Type: "Latitude",
	},
	"linkedin[-_]*skill": MockType{
		Type: "LinkedIn Skill",
	},
	"longitude": MockType{
		Type: "Longitude",
	},
	"mac[-_]*address": MockType{
		Type: "Mac Address",
	},
	"md5mime[-_]*type": MockType{
		Type: "MD5MIME Type",
	},
	"mongodb[-_]*object[-_]*id": MockType{
		Type: "MongoDB ObjectID",
	},
	"password": MockType{
		Type: "Password",
	},
	"phone": MockType{
		Type: "Phone",
	},
	"postal[-_]*code": MockType{
		Type: "Postal Code",
	},
	"race": MockType{
		Type: "Race",
	},
	"sha1": MockType{
		Type: "SHA1",
	},
	"sha256": MockType{
		Type: "SHA256",
	},
	"shirt[-_]*size": MockType{
		Type: "Shirt Size",
	},
	"ssn": MockType{
		Type: "SSN",
	},
	"state": MockType{
		Type: "State",
	},
	"street[-_]*(address){0,1}": MockType{
		Type: "Street Address",
	},
	"street[-_]*name": MockType{
		Type: "Street Name",
	},
	"street[-_]*(number|#)": MockType{
		Type: "Street Number",
	},
	"street[-_]*suffix": MockType{
		Type: "Street Suffix",
	},
	"suffix": MockType{
		Type: "Suffix",
	},
	"time[-_]*zone": MockType{
		Type: "Time Zone",
	},
	"title": MockType{
		Type: "title",
	},
	"user[-_]*agent": MockType{
		Type: "User Agent",
	},
}

func matchMochType(name string, mockType *MockType) bool {
	var match bool
	fmt.Println(mockType.Name)
	for exp, m := range mockRegex {
		if match, _ = regexp.MatchString(exp, strings.ToLower(mockType.Name)); match {
			mockType.Type = m.Type
			mockType.Min = m.Min
			mockType.Max = m.Max
			break
		}
	}
	return match
}
