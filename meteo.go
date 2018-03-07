package main

import (
  "fmt"
  "time"
  "log"
  "strings"
  "strconv"
  "flag"

  "github.com/joefitzgerald/forecast"
  "github.com/fatih/color"
)

type Filter = forecast.AssignmentFilter

func fetchAssignments(api *forecast.API, filter Filter) (assignments forecast.Assignments) {
  assignments, err := api.AssignmentsWithFilter(filter)

  if (err != nil) {
    log.Fatal(err)
  }
  return
}

func filterUpdatedSince(assignments forecast.Assignments, lastUpdated time.Time) (ret forecast.Assignments) {
  for _, a := range assignments {
    if lastUpdated.Before(a.UpdatedAt) {
      ret = append(ret, a)
    }
  }
  return
}

func printAssignments(api *forecast.API, assignments forecast.Assignments) {
  for _, assignment := range assignments {
    printAssignment(api, assignment)
  }
  fmt.Printf("\n")
}

func fetchProject(api *forecast.API, assignment forecast.Assignment) (project forecast.Project) {
  projects, err := api.Projects()
  if (err != nil) {
    log.Fatal(err)
  }

  for _, project = range projects {
    if project.ID == assignment.ProjectID {
      break
    }
  }

  return
}

func printAssignment(api *forecast.API, assignment forecast.Assignment) {
  project := fetchProject(api, assignment)

  green := color.New(color.FgGreen).PrintfFunc()
  yellow := color.New(color.FgYellow).PrintfFunc()
  green("%s", project.Name)
  fmt.Print(" (")
  yellow("%s - %s", assignment.StartDate, assignment.EndDate)
  fmt.Print(")\n")
}

func fetchPersonAssignmentsSince(api *forecast.API, pid int, lastChck time.Time) {
  person, err := api.Person(pid)
  if (err != nil) {
    log.Fatal(err)
  }

  fmt.Printf("=== %s %s ===\n", person.FirstName, person.LastName)

  today := time.Now().Format("2006-01-02")
  filter := Filter{StartDate: today, PersonID: pid}
  assgns := fetchAssignments(api, filter)
  assgns = filterUpdatedSince(assgns, lastChck)
  if len(assgns) > 0 {
    if (err != nil) {
      log.Fatal(err)
    }

    printAssignments(api, assgns)
  } else {
    fmt.Printf("Same old, same old\n\n")
  }
}

func createApi(acc *string, t *string)(api *forecast.API) {
  api = forecast.New("https://api.forecastapp.com/", *acc, *t)
  return
}

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

func main() {
  intervalPtr := flag.Int("interval", 3, "Interval between checks in hours (default: 3)")
  fcastAccIdPtr := flag.String("account_id", "0", "Forecast Account ID")
  fcastTokenPtr := flag.String("token", "0", "Forecast Token")
  peoplePtr := flag.String("people", "", "People to watch")
  flag.Parse()
  // Check Forecast API configuration values
  if (*fcastAccIdPtr == "0" || *fcastTokenPtr == "0") {
    log.Fatal("Missing Forecast API configuration.")
  }

  // Check People IDs
  people := mapPeopleIds(peoplePtr)
  if (len(people) == 0) {
    log.Fatal("No people IDs to watch.")
  }

  // Prepare datetime
  lastCheckedAt := time.Now().Add(time.Hour * time.Duration(*intervalPtr * -1))

  // Create API
  api := createApi(fcastAccIdPtr, fcastTokenPtr)

  // Go
  for _, personId := range people {
    fetchPersonAssignmentsSince(api, personId, lastCheckedAt)
  }

}
