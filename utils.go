package main

import (

)

func find_delimiter(entry string, delimiters []string) string {
  for _, e := range entry {
    for _, d := range delimiters {
      if d == string(e) {
        return string(d)
      }
    }
  }

  return ""
}
