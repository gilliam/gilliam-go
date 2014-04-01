package gilliam

import "fmt"

type Instance struct {
    Formation string
    Instance string
    Service string
    Name string
    Image string
    Command string
    Release string
    State string
    Status string
    Env map[string]string
    Ports []int
}

func (c *Client) FormationInstances(formation string) (insts []Instance , err error) {
    _, err = c.queryCollection(fmt.Sprintf(
        "http://api.scheduler.service/formation/%s/instances", formation), &insts)
    fmt.Println(insts)
    return
}
