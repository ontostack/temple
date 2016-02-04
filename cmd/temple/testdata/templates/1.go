<%! func MyTmpl(w io.Writer, funcs []*Func) error %>

<%% import "strings" %%>

<%for _, f := range funcs { %>

<%if len(f.Owner) > 0 {%>
func (c <%=strings.TrimSpace(f.Owner)%>) <%=f.Name%>() {
<%} else {%>
func <%=f.Name%>() {
<%}%>
	return <%= f.Default%>
}

<%}%>
