package page

import (
    "bytes"
    "net/url"
    "path"
    "strconv"
    "strings"
)

func (self *Page) Path(name string) *Result {
    pathString := "/"
    queryString := ""
    if uri, err := url.ParseRequestURI(self.result.CurrentURI); err == nil {
        pathString = uri.EscapedPath()
        if self.result.Exists {
            pathString = path.Dir(pathString)
        } else {
            strings.TrimRight(pathString, "/")
        }
        queryString = uri.Query().Encode()
        if queryString != "" {
            queryString = "?" + queryString
        }
    }

    name = escapePath(name)

    for k, v := range self.result.List {
        if v["num"] == "1" {
            self.result.List[k]["url"] = pathString + "/" + queryString
        } else {
            self.result.List[k]["url"] = pathString + "/" + name + v["num"] + "/" + queryString
        }
    }

    if self.result.Current == 1 || self.result.Prev == 1 {
        self.result.PrevURI = pathString + "/" + queryString
    } else {
        self.result.PrevURI = pathString + "/" + name + strconv.Itoa(self.result.Prev) + "/" + queryString
    }

    if self.result.Current == 1 {
        self.result.CurrentURI = pathString + "/" + queryString
    } else {
        self.result.CurrentURI = pathString + "/" + name + strconv.Itoa(self.result.Current) + "/" + queryString
    }

    if self.result.Current == self.result.Total {
        self.result.NextURI = pathString + "/" + name + strconv.Itoa(self.result.Current) + "/" + queryString
    } else {
        self.result.NextURI = pathString + "/" + name + strconv.Itoa(self.result.Next) + "/" + queryString
    }

    self.result.FirstURI = pathString + "/" + queryString

    if self.result.Total > 1 {
        self.result.LastURI = pathString + "/" + name + strconv.Itoa(self.result.Total) + "/" + queryString
    } else {
        self.result.LastURI = self.result.FirstURI
    }

    return self.result
}

func escapePath(s string) string {
    fun := func(c byte) bool {
        switch {
        case 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z':
            return false
        case '0' <= c && c <= '9':
            return false
        }
        switch c {
        case '-', '_', '.', '~':
            return false
        case '$', '&', '+', ',', ':', ';', '=', '@':
            return false
        }
        return true
    }
    buf := bytes.Buffer{}
    for i := 0; i < len(s); i++ {
        c := s[i]
        if fun(c) {
            buf.WriteByte('%')
            buf.WriteByte("0123456789ABCDEF"[c>>4])
            buf.WriteByte("0123456789ABCDEF"[c&15])
        } else {
            buf.WriteByte(c)
        }
    }
    return buf.String()
}
