// Code generated by reactGen. DO NOT EDIT.

package react

// NavProps defines the properties for the <nav> element
type NavProps struct {
	ClassName               string
	DangerouslySetInnerHTML *DangerousInnerHTMLDef
	ID                      string
	Key                     string

	OnChange
	OnClick

	Role  string
	Style *CSS
}

func (n *NavProps) assign(v *_NavProps) {

	v.ClassName = n.ClassName

	v.DangerouslySetInnerHTML = n.DangerouslySetInnerHTML

	if n.ID != "" {
		v.ID = n.ID
	}

	if n.Key != "" {
		v.Key = n.Key
	}

	if n.OnChange != nil {
		v.o.Set("onChange", n.OnChange.OnChange)
	}

	if n.OnClick != nil {
		v.o.Set("onClick", n.OnClick.OnClick)
	}

	v.Role = n.Role

	// TODO: until we have a resolution on
	// https://github.com/gopherjs/gopherjs/issues/236
	v.Style = n.Style.hack()

}
