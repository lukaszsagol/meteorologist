package main

import (
  "time"
  "log"
  "strings"
  "strconv"
  "flag"

  "github.com/lukaszsagol/meteorologist/forecast"
  "github.com/lukaszsagol/meteorologist/output"
)

func mapPeopleIds(ppl *string)(ret []int) {
  people := strings.Split(*ppl, ",")
  ret = make([]int, len(people))

  for i, p := range people {
    id, err := strconv.Atoi(p)
    if (err != nil) {
      log.Fatal(err)
    }
    ret[i] = id
  }

  return
}

func parseArguments()(int, string, string, []int, string, string) {
  intervalPtr := flag.Int("interval", 3, "Interval between checks in hours (default: 3)")
  accIdPtr := flag.String("account_id", "0", "Forecast Account ID")
  tokenPtr := flag.String("token", "0", "Forecast Token")
  peoplePtr := flag.String("people", "", "Comma separated list of people IDs")
  slackPtr := flag.String("slack", "", "Slack token for notifications")
  channelPtr := flag.String("channel", "", "Slack channel to send notifications")
  flag.Parse()

  //  Errors
  if (*peoplePtr == "") {
    log.Fatal("Missing people IDs")
  }

  if (*accIdPtr == "0" || *tokenPtr == "0") {
    log.Fatal("Missing Forecast API configuration.")
  }

  if (*slackPtr == "" || *channelPtr == "") {
    log.Fatal("Missing Slack configuration.")
  }

  people := mapPeopleIds(peoplePtr)

  if (len(people) == 0) {
    log.Fatal("No people IDs to watch.")
  }

  // All good
  return *intervalPtr, *accIdPtr, *tokenPtr, people, *slackPtr, *channelPtr
}

func main() {
  i, accId, token, people, slackToken, channel := parseArguments()

  lastCheckedAt := time.Now().Add(time.Hour * time.Duration(i * -1))

  api := forecast.CreateApi(accId, token)
  assignments := forecast.FetchAssignments(api, people, lastCheckedAt)
  projects := forecast.FetchProjects(api)

  output.SlackNotify(assignments, projects, slackToken, channel)
}
