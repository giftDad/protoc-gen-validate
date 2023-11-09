package rule

const requiredTpl = `
		if {{ .Key }} == nil {
			return {{ .Field.Parent.GoIdent.GoName }}ValidationError {
				field:  "{{ .Field.GoName }}",
				reason: "value is required",
			}
		}
`
