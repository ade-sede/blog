package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
)

var FLAGS = os.O_RDWR | os.O_CREATE
var MODE fs.FileMode = 0644

func buildWorkExperiences() []ExperienceEntry {
	workExperiences := make([]ExperienceEntry, 0)

	cubyn := ExperienceEntry{
		title:       "Senior Software Engineer, Tech Lead",
		company:     "Cubyn",
		begin:       "Nov. 2021",
		end:         "March 2024",
		description: "Logistics service provider, specialized in e-commerce with operations all over Europe and a mission to make logistics sustainable and socially responsible. I was recruited after the company closed its 35M Series B with a mission to automate Cubyn's warehouses in order to improve margins, quality and lead time.", bulletPoints: make([]string, 0),
	}

	cubyn.bulletPoints = append(cubyn.bulletPoints, "Led 4-person Warehouse Systems & Automation team, overseeing full software lifecycle from ideation to production.")
	cubyn.bulletPoints = append(cubyn.bulletPoints, "Mentored 2 direct reports.")
	cubyn.bulletPoints = append(cubyn.bulletPoints, "Architected and implemented automated e-commerce fulfillment system integrating AGVs (Automated Guided Vehicles).")
	cubyn.bulletPoints = append(cubyn.bulletPoints, "Scaled from 0% to 60% of company-wide order volume.")
	cubyn.bulletPoints = append(cubyn.bulletPoints, "Built using Domain-Driven Design, Hexagonal Architecture.")
	cubyn.bulletPoints = append(cubyn.bulletPoints, "Tech stack: TypeScript, Vue.js, PostgreSQL, RabbitMQ, Microservices, Kubernetes")

	cityTaps := ExperienceEntry{
		title:        "Lead Embedded Systems Developer",
		company:      "CityTaps",
		begin:        "July 2018",
		end:          "March 2021",
		description:  "Startup striving to enable fair access to running water in urban homes all over the world, using financial and technical innovation. The company’s main product is a water prepayment solution relying on a smart LoRaWAN water meter built to last 10 years on a D size battery, without maintenance. Joined at seed stage with the main mission of scaling from a prototype to a reliable production-ready system.",
		bulletPoints: make([]string, 0),
	}

	cityTaps.bulletPoints = append(cityTaps.bulletPoints, "Scaled IoT device firmware from prototype to production, growing deployment from 300 to 10,000 units across 3 countries.")
	cityTaps.bulletPoints = append(cityTaps.bulletPoints, "Developed robust embedded system ecosystem.")
	cityTaps.bulletPoints = append(cityTaps.bulletPoints, "Wrote core firmware in C for STM32 microcontrollers.")
	cityTaps.bulletPoints = append(cityTaps.bulletPoints, "Developed supporting tools and services in Node.js and Python.")
	cityTaps.bulletPoints = append(cityTaps.bulletPoints, "Reviewed and optimized PCB design, BOM.")
	cityTaps.bulletPoints = append(cityTaps.bulletPoints, "Developed custom test bench and test frameworks to guarantee quality.")
	cityTaps.bulletPoints = append(cityTaps.bulletPoints, "Led technical architecture and specification for product iterations.")

	workExperiences = append(workExperiences, cubyn)
	workExperiences = append(workExperiences, cityTaps)

	return workExperiences
}

func buildSchoolExperience() []ExperienceEntry {
	schoolExperiences := make([]ExperienceEntry, 0)

	fortyTwo := ExperienceEntry{
		title:             "Student",
		company:           "42 Paris",
		begin:             "Sep. 2016",
		end:               "July 2018",
		description:       "A very intensive programming curriculum focused on the C language, organised around successive projects and challenges that each highlight or deepen a particular programming concept. I also spent 6 month at the school's electronics lab learning about PCB design, firmware and embedded systems in general.",
		bulletPointsIntro: "Among other things",
		bulletPoints:      make([]string, 0),
	}

	fortyTwo.bulletPoints = append(fortyTwo.bulletPoints, "Developed feature rich shell in C.")
	fortyTwo.bulletPoints = append(fortyTwo.bulletPoints, "Implemented Supervisor-like utility in Rust.")
	fortyTwo.bulletPoints = append(fortyTwo.bulletPoints, "Built a wireless mouse controlled by hand gestures.")

	return append(schoolExperiences, fortyTwo)
}

func main() {
	outputDir := os.Getenv("OUTPUT_DIR")
	indexFileName := fmt.Sprintf("%s/index.html", outputDir)
	indexFile, err := os.OpenFile(indexFileName, FLAGS, MODE)
	if err != nil {
		log.Fatal("Could not open file: %w", err)
	}

	aboutFileName := fmt.Sprintf("%s/about.html", outputDir)
	aboutFile, err := os.OpenFile(aboutFileName, FLAGS, MODE)
	if err != nil {
		log.Fatal("Could not open file: %w", err)
	}

	resumeLightFileName := fmt.Sprintf("%s/resume-light.html", outputDir)
	resumeLightFile, err := os.OpenFile(resumeLightFileName, FLAGS, MODE)
	if err != nil {
		log.Fatal("Could not open file: %w", err)
	}

	home := home(make(map[string]string))
	home.Render(context.Background(), indexFile)

	workExperiences := buildWorkExperiences()
	schoolExperiences := buildSchoolExperience()

	about := about(make(map[string]string), workExperiences, schoolExperiences)
	about.Render(context.Background(), aboutFile)

	// Print resume to a standalone HTMl file so that we can easily create a PDF from it
	// Will most likely be removed before serving the content
	resumeLight := resumeLight(workExperiences, schoolExperiences)
	resumeLight.Render(context.Background(), resumeLightFile)
}
