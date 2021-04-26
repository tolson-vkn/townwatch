package watcher

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
	"text/template"
	"text/template/parse"

	"github.com/sirupsen/logrus"
)

func ListTemplFields(t *template.Template) []string {
	return listNodeFields(t.Tree.Root, nil)
}

func listNodeFields(node parse.Node, res []string) []string {
	if node.Type() == parse.NodeAction {
		logrus.Info(res)
		res = append(res, strings.Trim(node.String(), "{}. "))
	}

	if ln, ok := node.(*parse.ListNode); ok {
		for _, n := range ln.Nodes {
			res = listNodeFields(n, res)
		}
	}
	return res
}

func InspectExampleLines(regEx string, lines []string) error {
	for _, line := range lines {
		r := regexp.MustCompile(regEx)
		if inspectLine(r, line) == nil {
			return errors.New("Regular expression did not yield any matches.")
		}
	}
	return nil
}

// Take compiled regex and compare against line. Return a map of subexpressions
// and values.
func inspectLine(regEx *regexp.Regexp, line string) (regexMap map[string]string) {

	match := regEx.FindStringSubmatch(line)
	if match == nil {
		return nil
	}
	regexMap = make(map[string]string)
	for i, name := range regEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			regexMap[name] = match[i]
		}
	}

	return
}

// Would be nice to figure out a better solution. If err in here it panics. Thought Must
// was to blame but it's somewhere else in the template package.
func (watcher *Watcher) parseTemplates(caps map[string]string) (title, message string, err error) {
	var buf bytes.Buffer

	titleTmpl, err := template.New("NotifTitle").Parse(watcher.Title)
	if err != nil {
		return "", "", err
	}
	err = titleTmpl.Execute(&buf, caps)
	if err != nil {
		return "", "", err
	}
	title = buf.String()
	// Clear buffer
	buf.Reset()

	messageTmpl, err := template.New("NotifMessage").Parse(watcher.Message)
	if err != nil {
		return "", "", err
	}
	err = messageTmpl.Execute(&buf, caps)
	if err != nil {
		return "", "", err
	}
	message = buf.String()

	return title, message, nil
}
