package forecast
import (
  "fmt"
  "time"
  "log"

  fapi "github.com/joefitzgerald/forecast"
)

type Person struct {
  Name string
  Assignments fapi.Assignments
}

func fetchPersonalAssgns(api *fapi.API, f fapi.AssignmentFilter)(a fapi.Assignments) {
  a, err := api.AssignmentsWithFilter(f)

  if (err != nil) {
    log.Fatal(err)
  }

  return
}

func updatedSince(a fapi.Assignments, last time.Time) (ret fapi.Assignments) {
  for _, a := range a {
    if last.Before(a.UpdatedAt) {
      ret = append(ret, a)
    }
  }
  return
}

func fetchPerson(api *fapi.API, pid int, last time.Time)(p Person) {
  person, err := api.Person(pid)

  if (err != nil) {
    log.Fatal(err)
  }


  f := fapi.AssignmentFilter{
    StartDate: time.Now().Format("2006-01-02"),
    PersonID: pid,
  }

  p = Person{
    Name: fmt.Sprintf("%s %s", person.FirstName, person.LastName),
    Assignments: updatedSince(fetchPersonalAssgns(api, f), last),
  }

  return
}

func CreateApi(acc string, t string)(api *fapi.API) {
  api = fapi.New("https://api.forecastapp.com/", acc, t)
  return
}

func FetchAssignments(api *fapi.API, people []int, last time.Time)(asgns []Person) {
  asgns = make([]Person, len(people))
  for i, pid := range people {
    asgns[i] = fetchPerson(api, pid, last)
  }
  return
}

func FetchProject(api *fapi.API, assignment fapi.Assignment) (project fapi.Project) {
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

func FetchProjects(api *fapi.API)(p map[int]string) {
  p = make(map[int]string)

  projects, err := api.Projects()
  if (err != nil) {
    log.Fatal(err)
  }

  for _, pro := range projects {
    p[pro.ID] = pro.Name
  }

  return
}
