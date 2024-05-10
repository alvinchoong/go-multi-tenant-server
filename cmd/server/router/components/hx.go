package components

type HXMethod string

const (
	HXMethodGet    HXMethod = "hx-get"
	HXMethodPost   HXMethod = "hx-post"
	HXMethodPut    HXMethod = "hx-put"
	HXMethodDelete HXMethod = "hx-delete"
	HXMethodPatch  HXMethod = "hx-patch"
)

func (e HXMethod) String() string {
	return string(e)
}
