package main

import (
	"context"
	"embed"
	"fmt"
	"slices"
	"strings"

	"github.com/mpetavy/common"
)

//go:embed go.mod
var resources embed.FS

func init() {
	common.Init("", "", "", "", "", "", "", "", &resources, nil, nil, run, 0)
}

type Leaf struct {
	Name  string
	Leafs []*Leaf
}

func (leaf *Leaf) Find(name string) *Leaf {
	for _, child := range leaf.Leafs {
		if child.Name == name {
			return child
		}
	}

	return nil
}

func (leaf *Leaf) ToString(prefixLen int) string {
	prefix := strings.Repeat(" ", prefixLen*2)
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("%s%s\n", prefix, leaf.Name))

	for _, child := range leaf.Leafs {
		sb.WriteString(child.ToString(prefixLen + 1))
	}

	return sb.String()
}

func SliceFind[S ~[]E, E comparable](s S, e E) *E {
	index := slices.Index(s, e)
	if index == -1 {
		return nil
	}

	return &s[index]
}

func SliceFindFunc[S ~[]E, E any](s S, f func(E) bool) *E {
	index := slices.IndexFunc(s, f)
	if index == -1 {
		return nil
	}

	return &s[index]
}

func GetOrCreate(root *Leaf, name string) *Leaf {
	leaf := root.Find(name)
	if leaf == nil {
		leaf = &Leaf{Name: name}

		root.Leafs = append(root.Leafs, leaf)
	}

	return leaf
}

func GetPathPart(path string, index int) string {
	splits := common.Split(path, "/")

	return splits[index]
}

func run() error {
	tasks := common.NewTasks(context.Background())

	subscriptions := make([]Subscription, 0)
	tasks.Add(func(ctx context.Context) error {
		return fmt.Errorf("test-error")

		var err error

		subscriptions, err = ReadSubscriptions()
		if common.Error(err) {
			return err
		}

		return nil
	})

	resources := make([]Resource, 0)
	tasks.Add(func(ctx context.Context) error {
		var err error

		resources, err = ReadResources()
		if common.Error(err) {
			return err
		}

		return nil
	})

	err := tasks.Wait()
	common.Error(err)
	if common.Error(err) {
		return err
	}

	root := &Leaf{Name: "root"}

	for _, resource := range resources {
		subscriptionId := GetPathPart(resource.Id, 2)
		subscription := SliceFindFunc(subscriptions, func(s Subscription) bool {
			return s.Id == subscriptionId
		})

		subscriptionLeaf := GetOrCreate(root, subscription.Name+"(subscription)")
		resourceGroupLeaf := GetOrCreate(subscriptionLeaf, resource.ResourceGroup+"(Resourcegroup)")
		GetOrCreate(resourceGroupLeaf, resource.Name)
	}

	fmt.Printf("%s\n", root.ToString(0))

	return nil
}

func main() {
	common.Run(nil)
}
