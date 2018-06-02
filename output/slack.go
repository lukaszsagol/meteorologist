package output

import (
  "fmt"
  "log"

  "github.com/nlopes/slack"
  "github.com/lukaszsagol/meteorologist/forecast"
)

func formatMessage(p []forecast.Person, pro map[int]string)(msg string){
  msg = ""
  for _, pi := range p {
    if (len(pi.Assignments) > 0) {
      msg += fmt.Sprintf("*Updates for %s*\n", pi.Name)
      for _, a := range pi.Assignments {
        msg += fmt.Sprintf("\t-> %s (%s - %s)\n", pro[a.ProjectID], a.StartDate, a.EndDate)
      }
      msg += "\n\n"
    }
  }
  return
}

func SlackNotify(p []forecast.Person, pro map[int]string, t string, c string) {
  api := slack.New(t)
  msg := formatMessage(p, pro)

  if msg == "" {
    return
  }

  _, _, id, err := api.OpenIMChannel(c)
  if err != nil {
    log.Fatal(err)
  }

  _, _, err = api.PostMessage(id, formatMessage(p, pro), slack.PostMessageParameters{Markdown: true})
  if err != nil {
    log.Fatal(err)
  }

}
