# Pipeline
This Go package provides a versatile implementation of a pipeline pattern, enabling the construction of complex processing sequences in a concurrent and efficient manner. It's designed to offer flexibility, allowing users to define custom nodes and control the flow of data through the pipeline.

## Features

- **Versatile**: It allows you to construct complex processing sequences in a concurrent and efficient manner.
- **Flexible**: You can define custom nodes and control the flow of data through the pipeline.
- **Efficient**: It is designed to maximize the use of available resources by running tasks concurrently.

## Key Components

- **INode Interface**: This is the core interface that defines essential methods for pipeline nodes.
- **Node[T any]**: This is a generic node structure that represents a single stage in the pipeline.
- **Pipeline**: This is the orchestrator for managing the pipeline execution flow.

## Usage

### Creating Nodes

To create a new node, you can use the `NewNode` or `NewDefaultNode` function. Here is an example:

```go
workFunc := func(job int) (any, error) {
    // Your processing logic here
    return nil, nil
}
node := pipeline.NewDefaultNode("node1", workFunc)
```

### Constructing a Pipeline

After creating the nodes, you can construct a pipeline. Here is an example:

```go
rootNode := pipeline.NewDefaultNode("root", workFunc)
node1 := pipeline.NewDefaultNode("node1", workFunc)
node2 := pipeline.NewDefaultNode("node2", workFunc)

rootNode.SetNext(node1)
node1.SetNext(node2)

p, err := pipeline.NewPipeline(rootNode)
if err != nil {
    // Handle error
}
```

### Running the Pipeline

To run the pipeline, you can use the `Start` method. Here is an example:

```go
doneCh := p.Start()
```

To stop the pipeline, you can use the `Stop` method. Here is an example:

```go
p.Stop()
```

To wait for all nodes to finish, you can use the `Wait` method. Here is an example:

```go
p.Wait()
```

### Passing Jobs to the Pipeline

To pass a job to the pipeline, you can use the `JobReceiver` method. Here is an example:

```go
err := p.JobReceiver(job)
if err != nil {
    // Handle error
}
```

### Getting the Latest Job Processed by the Pipeline

To get the latest job processed by the pipeline, you can use the `LatestJob` method. Here is an example:

```go
latestJob := p.LatestJob()
```