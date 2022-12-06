package main

import (
  "context"
	"io"
	"os"
  "encoding/json"
  "flag"
  "log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
  "github.com/docker/docker/api/types/network"
)

var containerConfig string
var hostConfig string
var name string

type State struct{
  ContainerConfig *container.Config
  HostConfig *container.HostConfig
}

func init(){
  flag.StringVar(&containerConfig, "c", "./container.json", "path to the container config as JSON")
  flag.StringVar(&hostConfig, "h", "./host.json", "path to the host config as JSON")
  flag.StringVar(&name, "n", "foobar", "container name")
}


func (s *State)getContainerConfig()(error){
  log.Print("Reading container config from: ", containerConfig)
  ccbytes, err:=os.ReadFile(containerConfig)
  if err!=nil{
    return err
  }
  return json.Unmarshal(ccbytes, s.ContainerConfig)
}
func (s *State)getHostConfig()(error){
  log.Print("Reading host config from: ", hostConfig)
  ccbytes, err:=os.ReadFile(hostConfig)
  if err!=nil{
    return err
  }
  return json.Unmarshal(ccbytes, s.HostConfig)
}

func(s *State)StartContainer(containerName string)(error){
  ctx := context.Background()
  cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
  if err != nil {
    log.Fatal(err)
  }
  defer cli.Close()
  reader, err := cli.ImagePull(ctx, s.ContainerConfig.Image, types.ImagePullOptions{})
  if err != nil {
    log.Fatal(err)
  }
  defer reader.Close()
  io.Copy(os.Stdout, reader)
  resp, err := cli.ContainerCreate(ctx, s.ContainerConfig, s.HostConfig, &network.NetworkingConfig{}, nil, containerName)
  if err != nil {
    log.Fatal(err)
  }
  if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
    log.Fatal(err)
  }
  return nil
}
func (s *State)RemoveContainer(containerName string)(error){
  ctx := context.Background()
  cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
  if err != nil {
    log.Fatal(err)
  }
  defer cli.Close()
  containerList, err := cli.ContainerList(ctx, types.ContainerListOptions{
    All: true,
  })
  if err != nil {
    return err
  }
  containerId:=""
  for _, container := range containerList{
    if findInList(container.Names, "/"+containerName) != -1{
      containerId=container.ID
      break
    }
  }
  if containerId==""{
    return nil
  } 
  return cli.ContainerRemove(ctx, containerId, types. ContainerRemoveOptions{
    RemoveVolumes: true,
    Force: true,
  })
}

func findInList(list []string, element string)(int){
  for index, el:=range list{
    if el == element{
      return index
    }
  }
  return -1
}

func main(){
  flag.Parse()
  s:=&State{
    ContainerConfig: &container.Config{},
    HostConfig: &container.HostConfig{},
  }
  err:=s.getContainerConfig()
  if err!=nil{
    log.Fatal(err)
  }
  err=s.getHostConfig()
  if err!=nil{
    log.Fatal(err)
  }
  err=s.RemoveContainer(name)
  if err!=nil{
    log.Fatal(err)
  }
  err=s.StartContainer(name)
  if err!=nil{
    log.Fatal(err)
  }
  log.Print("Container started successfully")
}
