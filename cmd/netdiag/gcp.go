package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"google.golang.org/api/iterator"

	"github.com/autonity/autonity/cmd/netdiag/ssh"
	"github.com/autonity/autonity/log"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

const (
	binaryName  = "netdiag"
	logFileName = "output.log"
)

var ZoneMap = map[string][]string{
	"c2-standard-8": []string{
		"us-east1-b",
		"us-east1-c",
		"us-east1-d",
		"us-east4-c",
		"us-east4-b",
		"us-east4-a",
		"us-central1-c",
		"us-central1-a",
		"us-central1-f",
		"us-central1-b",
		"us-west1-b",
		"us-west1-c",
		"us-west1-a",
		"europe-west4-a",
		"europe-west4-b",
		"europe-west4-c",
		"europe-west1-b",
		"europe-west1-d",
		"europe-west1-c",
		"europe-west3-c",
		"europe-west3-a",
		"europe-west3-b",
		"europe-west2-c",
		"europe-west2-b",
		"europe-west2-a",
		"asia-east1-b",
		"asia-east1-a",
		"asia-east1-c",
		"asia-southeast1-b",
		"asia-southeast1-a",
		"asia-southeast1-c",
		"asia-northeast1-b",
		"asia-northeast1-c",
		"asia-northeast1-a",
		"asia-south1-c",
		"asia-south1-b",
		"asia-south1-a",
		"australia-southeast1-b",
		"australia-southeast1-c",
		"australia-southeast1-a",
		"southamerica-east1-b",
		"southamerica-east1-c",
		"southamerica-east1-a",
		"asia-east2-a",
		"asia-east2-b",
		"asia-east2-c",
		"asia-northeast2-a",
		"asia-northeast2-b",
		"asia-northeast2-c",
		"asia-northeast3-a",
		"asia-northeast3-b",
		"asia-northeast3-c",
		"asia-south2-a",
		"asia-south2-b",
		"asia-south2-c",
		"europe-north1-a",
		"europe-north1-b",
		"europe-north1-c",
		"europe-west6-a",
		"europe-west6-b",
		"europe-west6-c",
		"me-west1-a",
		"me-west1-b",
		"me-west1-c",
		"northamerica-northeast1-a",
		"northamerica-northeast1-b",
		"northamerica-northeast1-c",
		"southamerica-west1-a",
		"southamerica-west1-b",
		"southamerica-west1-c",
		"us-east5-a",
		"us-east5-b",
		"us-east5-c",
		"us-west2-a",
		"us-west2-b",
		"us-west2-c",
		"us-west3-a",
		"us-west3-b",
		"us-west3-c",
		"us-west4-a",
		"us-west4-b",
		"us-west4-c",
	},
	"c2-standard-4": {
		"us-central1-a",
		"us-central1-b",
		"us-central1-c",
		"us-central1-f",
		"europe-west1-b",
		"europe-west1-c",
		"europe-west1-d",
		"us-west1-a",
		"us-west1-b",
		"us-west1-c",
		"asia-east1-a",
		"asia-east1-b",
		"asia-east1-c",
		"us-east1-a",
		"us-east1-b",
		"us-east1-c",
		"us-east1-d",
		"asia-northeast1-a",
		"asia-northeast1-b",
		"asia-northeast1-c",
		"asia-southeast1-a",
		"asia-southeast1-b",
		"asia-southeast1-c",
		"us-east4-a",
		"us-east4-b",
		"us-east4-c",
		"australia-southeast1-c",
		"australia-southeast1-a",
		"australia-southeast1-b",
		"europe-west2-a",
		"europe-west2-b",
		"europe-west2-c",
		"europe-west3-c",
		"europe-west3-a",
		"europe-west3-b",
		"southamerica-east1-a",
		"southamerica-east1-b",
		"southamerica-east1-c",
		"asia-south1-b",
		"asia-south1-a",
		"asia-south1-c",
		"northamerica-northeast1-a",
		"northamerica-northeast1-b",
		"northamerica-northeast1-c",
		"europe-west4-c",
		"europe-west4-b",
		"europe-west4-a",
		"europe-north1-b",
		"europe-north1-c",
		"europe-north1-a",
		"us-west2-c",
		"us-west2-b",
		"us-west2-a",
		"asia-east2-c",
		"asia-east2-b",
		"asia-east2-a",
		"europe-west6-b",
		"europe-west6-c",
		"europe-west6-a",
		"asia-northeast2-b",
		"asia-northeast2-c",
		"asia-northeast2-a",
		"asia-northeast3-a",
		"asia-northeast3-c",
		"asia-northeast3-b",
		"us-west3-a",
		"us-west3-b",
		"us-west3-c",
		"us-west4-c",
		"us-west4-a",
		"us-west4-b",
		"asia-south2-a",
		"asia-south2-c",
		"asia-south2-b",
		"southamerica-west1-a",
		"southamerica-west1-b",
		"southamerica-west1-c",
		"us-east7-c",
		"us-east5-c",
		"us-east5-b",
		"us-east5-a",
		"me-west1-b",
		"me-west1-a",
		"me-west1-c",
		"me-central2-c",
	},
}

type vm struct {
	id           int
	ip           string
	instanceName string
	zone         string
	user         string
}

// prefer using newVM over directly instantiating a vm object
func newVM(id int, ip, instanceName, zone, user string) *vm {
	return &vm{
		id:           id,
		ip:           ip,
		instanceName: instanceName,
		zone:         zone,
		user:         user,
	}
}

func deployVM(ctx context.Context, client *compute.InstancesClient, id int, projectID, instanceName, zone, instanceTemplate, user string) (*vm, error) {
	// Create a new VM instance
	if err := createInstance(ctx, client, projectID, zone, instanceName, instanceTemplate); err != nil {
		return nil, err
	}
	// Get instance external IP
	ipAddress := getInstanceExternalIP(ctx, client, projectID, zone, instanceName)
	return newVM(id, ipAddress, instanceName, zone, user), nil
}

func (vm *vm) deployRunner(configFileName string, debug bool, skipConfigDeploy bool) error {
	client := ssh.Manager(vm.ip, vm.user)
	if client == nil {
		log.Error("ssh client not found", "ip", vm.ip)
		return nil
	}

	if !skipConfigDeploy {
		log.Info("Transferring Config file to the VM...", "id", vm.id)
		// Send the binary to the VM
		err := client.Upload(configFileName, configFileName, true)
		if err != nil {
			log.Error("upload failure", "err", err, "file", configFileName)
			return err
		}
	}
	exePath, err := os.Executable()
	if err != nil {
		log.Error("os executable", "err", err)
		return err
	}
	err = client.Upload(exePath, binaryName, true)
	if err != nil {
		log.Error("upload failure", "err", err, "file", configFileName)
		return err
	}

	serverCert := "./server.crt"
	log.Info("Transferring server certificates to the VM...", "id", vm.id)
	// upload server certificate to the VM
	err = client.Upload(serverCert, serverCert, true)
	if err != nil {
		log.Error("upload cert failure", "err", err, "file", serverCert)
		return err
	}

	// upload server key to the VM
	serverKey := "./server.key"
	err = client.Upload(serverKey, serverKey, true)
	if err != nil {
		log.Error("upload key failure", "err", err, "file", serverKey)
		return err
	}
	return nil
}

func (vm *vm) startRunner(configFileName, networkMode string, optFlags string) error {
	client := ssh.Manager(vm.ip, vm.user)
	if client == nil {
		log.Error("ssh client not found", "ip", vm.ip)
		return nil
	}

	// enable groupids for icmp ping
	cmd := "sudo sysctl -w net.ipv4.ping_group_range='0 2147483647'"
	_, err := client.RunCommand(cmd)
	if err != nil {
		log.Error("run failure", "err", err, "cmd", cmd)
		return err
	}
	// Execute the binary on the VM
	log.Info("Executing the runner on the VM...", "id", vm.id)

	flags := fmt.Sprintf("run --Config %s --id %d --network %s %s", configFileName, vm.id, networkMode, optFlags)
	localCommand := fmt.Sprintf("chmod +x ~/%s && sudo -b ./%s %s > %s 2>&1 ", binaryName, binaryName, flags, logFileName)
	log.Info("local command", "cmd", localCommand)

	_, err = client.RunCommand(localCommand)
	if err != nil {
		log.Error("run failure", "err", err, "cmd", cmd)
		return err
	}

	// copy server key files
	return nil
}

func (vm *vm) killRunner(configFileName string) error {
	client := ssh.Manager(vm.ip, vm.user)
	if client == nil {
		log.Error("ssh client not found", "ip", vm.ip)
		return nil
	}
	// Execute the binary on the VM
	log.Info("Killing the runner on the VM...", "id", vm.id)
	localCommand := fmt.Sprintf("sudo killall -9 %s", binaryName)
	_, err := client.RunCommand(localCommand)
	if err != nil {
		log.Error("run failure", "err", err, "cmd", localCommand)
	}
	return nil
}

func (vm *vm) downloadLogs() {
	log.Info("Downloading the log file...", "id", vm.id)
	client := ssh.Manager(vm.ip, vm.user)
	if client == nil {
		log.Error("ssh client not found", "ip", vm.ip)
		return
	}
	err := client.Download(logFileName, ".", true)
	if err != nil {
		log.Info("Error downloading log file", "error ", err)
	}
}

func (vm *vm) deleteRunner(ctx context.Context, client *compute.InstancesClient, projectID string) {
	// Download the log file
	// Delete the VM
	fmt.Println("Deleting the VM...")
	deleteInstance(ctx, client, projectID, vm.zone, vm.instanceName)
}

func createInstance(ctx context.Context, client *compute.InstancesClient, projectID, zone, instanceName, instanceTemplate string) error {
	log.Info("Creating a new VM instance...", "name", instanceName)

	op, err := client.Insert(ctx, &computepb.InsertInstanceRequest{
		Project: projectID,
		Zone:    zone,
		InstanceResource: &computepb.Instance{
			Name: proto.String(instanceName),
		},
		SourceInstanceTemplate: proto.String(instanceTemplate),
	})

	if err != nil {
		return fmt.Errorf("unable to create instance: %v", err)
	}

	log.Info("Waiting for the vm to deploy...", "name", instanceName)
	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}
	log.Info("Instance created", "name", instanceName)
	return nil
}

func getInstanceExternalIP(ctx context.Context, client *compute.InstancesClient, projectID, zone, instanceName string) string {
	req := &computepb.GetInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instanceName,
	}
	inst, err := client.Get(ctx, req)
	if err != nil {
		log.Error("Error getting instance", "err", err)
	}

	for _, intf := range inst.GetNetworkInterfaces() {
		if ac := intf.GetAccessConfigs(); len(ac) > 0 {
			return ac[0].GetNatIP()
		}
	}
	log.Crit("No external IP found for the ineval $(ssh-agent)stance.")
	return ""
}

func deleteInstance(ctx context.Context, client *compute.InstancesClient, projectID, zone, instanceName string) error {
	op, err := client.Delete(ctx, &computepb.DeleteInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instanceName,
	})
	if err != nil {
		log.Error("Error deleting instance", "err", err)
		return err
	}

	fmt.Println("Waiting for the instance to be deleted...")
	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}
	fmt.Printf("Instance %s deleted.\n", instanceName)
	return nil
}

func listZones(projectID string) ([]*computepb.Zone, error) {
	ctx := context.Background()

	// Create a new Compute Engine client.
	c, err := compute.NewZonesRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("compute.NewZonesRESTClient: %v", err)
	}
	defer c.Close()

	// Build the request.
	req := &computepb.ListZonesRequest{
		Project: projectID,
	}

	// Send the request to list zones
	zones := make([]*computepb.Zone, 0)
	it := c.List(ctx, req)
	for {
		zone, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list zones: %v", err)
		}

		zones = append(zones, zone)
		//fmt.Println(zone.GetName())
	}

	return zones, nil
}

func filterZones(zones []*computepb.Zone, machine string) []*computepb.Zone {
	log.Info("machine", "name", machine)
	filteredZones := make([]*computepb.Zone, 0)
	if zoneList, ok := ZoneMap[machine]; ok {
		for _, zone := range zoneList {
			for _, z := range zones {
				if z.GetName() == zone {
					filteredZones = append(filteredZones, z)
				}
			}
		}
	}
	return filteredZones
}
func getInstanceTemplate(projectID, templateName string) (*computepb.InstanceTemplate, error) {
	ctx := context.Background()
	instTemplateCl, err := compute.NewInstanceTemplatesRESTClient(ctx)
	if err != nil {
		log.Error("New Instance Rest Client error", "err", err)
		return nil, err
	}
	defer instTemplateCl.Close()

	req := &computepb.GetInstanceTemplateRequest{}
	req.Project = projectID
	req.InstanceTemplate = templateName
	temp, err := instTemplateCl.Get(ctx, req)
	if err != nil {
		log.Error("Error getting instance template", "err:", err)
	}
	return temp, err
}
