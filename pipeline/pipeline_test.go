package pipeline

import (
	"encoding/json"
	"fmt"
	"github.com/huahuayu/kit/logger"
	"io"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

var (
	welcomer = NewDefaultNode[string]("welcomer", welcome)
	nodeA    = NewDefaultNode[string]("nodeA", nodeAImpl)
	nodeB    = NewDefaultNode[string]("nodeB", nodeBImpl)
	nodeC    = NewDefaultNode[string]("nodeC", nodeCImpl)
	nodeD    = NewDefaultNode[string]("nodeD", nodeDImpl)
	nodeE    = NewDefaultNode[string]("nodeE", nodeEImpl)
	nodeF    = NewDefaultNode[string]("nodeF", nodeFImpl)
	nodeG    = NewDefaultNode[string]("nodeG", nodeGImpl)
	nodeH    = NewDefaultNode[string]("nodeH", nodeHImpl)
	printer  = NewDefaultNode[string]("printer", print)
	printer1 = NewDefaultNode[string]("printer1", print1)
	printer2 = NewDefaultNode[string]("printer2", print2)
)

/*
*
Here is a text graph to visualize the pipeline:
welcomer

	|
	v

nodeA

	|\
	| v
	|nodeE
	|  |
	|  v
	|nodeF
	|  |
	|  v
	|printer1
	v

nodeB

	|
	v

nodeC

	|\
	| v
	|nodeG
	|  |
	|  v
	|nodeH
	|  |
	|  v
	|printer2
	v

nodeD

	|
	v

printer
*/
func TestPipeline_Start(t *testing.T) {
	logger.Logger.SetLevel(logger.LevelDebug)
	welcomer.SetNext(nodeA)
	nodeA.SetNext(nodeB)
	nodeB.SetNext(nodeC)
	nodeC.SetNext(nodeD)
	nodeD.SetNext(printer)
	nodeA.SetNext(nodeE)
	nodeE.SetNext(nodeF)
	nodeF.SetNext(printer1)
	nodeC.SetNext(nodeG)
	nodeG.SetNext(nodeH)
	nodeH.SetNext(printer2)
	p, err := NewPipeline(welcomer)
	if err != nil {
		t.Fatal(err)
	}
	doneSignal := p.Start()
	go time.AfterFunc(time.Second*5, func() {
		p.Stop()
	})
	producer(10, doneSignal)
	p.Wait()
	fmt.Println("latestJob processed: ", p.LatestJob())
	fmt.Println("exit!")
}

func producer(num int, stopSignal chan struct{}) {
	fmt.Println("producer started")
	for i := 0; i < num; i++ {
		select {
		case <-stopSignal:
			fmt.Println("producer stop because of pipeline stop signal")
			return
		default:
			if name, _ := nameGenerator(); name != "" {
				go func(i int, name string) {
					welcomer.jobReceiver(name)
					fmt.Printf("producer send: %d %s\n", i, name)
				}(i, name)
			}
		}
	}
	fmt.Println("producer finished")
}

func welcome(name string) (any, error) {
	time.Sleep(time.Millisecond * time.Duration(1000+rand.Intn(2000)))
	return "Hello " + name, nil
}

func print(str string) (any, error) {
	fmt.Println(str)
	return nil, nil
}

func print1(str string) (any, error) {
	fmt.Println(str)
	return nil, nil
}
func print2(str string) (any, error) {
	fmt.Println(str)
	return nil, nil
}

func nodeAImpl(str string) (any, error) {
	time.Sleep(time.Millisecond * time.Duration(1000+rand.Intn(2000)))
	return str + " from nodeA", nil
}

func nodeBImpl(str string) (any, error) {
	time.Sleep(time.Millisecond * time.Duration(1000+rand.Intn(2000)))
	return str + " from nodeB", nil
}

func nodeCImpl(str string) (any, error) {
	time.Sleep(time.Millisecond * time.Duration(1000+rand.Intn(2000)))
	return str + " from nodeC", nil
}

func nodeDImpl(str string) (any, error) {
	time.Sleep(time.Millisecond * time.Duration(1000+rand.Intn(2000)))
	return str + " from nodeD", nil
}

func nodeEImpl(str string) (any, error) {
	time.Sleep(time.Millisecond * time.Duration(1000+rand.Intn(2000)))
	return str + " from nodeE", nil
}

func nodeFImpl(str string) (any, error) {
	time.Sleep(time.Millisecond * time.Duration(1000+rand.Intn(2000)))
	return str + " from nodeF", nil
}

func nodeGImpl(str string) (any, error) {
	time.Sleep(time.Millisecond * time.Duration(1000+rand.Intn(2000)))
	return str + " from nodeG", nil
}
func nodeHImpl(str string) (any, error) {
	time.Sleep(time.Millisecond * time.Duration(1000+rand.Intn(2000)))
	return str + " from nodeH", nil
}

func nameGenerator() (string, error) {
	if resp, err := http.Get("https://api.namefake.com/"); err != nil {
		return "", err
	} else {
		if body, err := io.ReadAll(resp.Body); err != nil {
			return "", err
		} else {
			var res = struct {
				Name string `json:"name"`
			}{}
			err = json.Unmarshal(body, &res)
			if err != nil {
				return "", err
			}
			return res.Name, nil
		}
	}
}
func TestPipelineCycleCheck(t *testing.T) {
	// Create a pipeline with a cycle
	node1 := NewDefaultNode[string]("node1", nodeAImpl)
	node2 := NewDefaultNode[string]("node2", nodeBImpl)
	node2_1 := NewDefaultNode[string]("node2_1", nodeBImpl)
	node3 := NewDefaultNode[string]("node3", nodeCImpl)
	node1.SetNext(node2)
	node2.SetNext(node2_1)
	node2_1.SetNext(node3)
	node3.SetNext(node2) // This creates a cycle
	_, err := NewPipeline(node1)
	if err.Error() != "cycle not allowed in pipeline" {
		t.Fatal(err)
	}

	// Create a pipeline without a cycle
	node4 := NewDefaultNode[string]("node4", nodeAImpl)
	node5 := NewDefaultNode[string]("node5", nodeBImpl)
	node6 := NewDefaultNode[string]("node6", nodeCImpl)
	node4.SetNext(node5)
	node5.SetNext(node6)
	_, err = NewPipeline(node4)
	if err != nil {
		t.Fatal(err)
	}
}
