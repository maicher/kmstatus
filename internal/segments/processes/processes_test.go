package processes

import "testing"

func Test_ParseTemplate(t *testing.T) {
	d, err := parseTemplate("F firefox\n\nS Socket Process\n")
	if err != nil {
		t.Fatalf("Error: %s, want: nil", err)
	}

	if len(d) != 2 {
		t.Fatalf("Entries: %d, want: 2", len(d))
	}

	if d[1].icon != "S" || d[1].phrase != "Socket Process" {
		t.Fatalf("Entry: %+v, want: icon S, phrase Socket Process", d[1])
	}
}

func Test_ParseTemplate_MissingProcessName(t *testing.T) {
	_, err := parseTemplate("F firefox\nX\n")
	if err == nil {
		t.Fatalf("Error nil, want: error")
	}
}
