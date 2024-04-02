package main

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	"google.golang.org/api/iterator"

	"github.com/autonity/autonity/log"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

const (
	binaryFilePath = "/root/netdiag"
	binaryName     = "netdiag"
	logFileName    = "output.log"
)

type vm struct {
	ip             string
	instanceName   string
	zone           string
	instanceClient *compute.InstancesClient // don't forget to close ! defer instancesClient.Close()
}

func deployVM(projectID, instanceName, zone, instanceTemplate string) (*vm, error) {
	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		log.Error("NewInstancesRESTClient: %v", err)
		return nil, err
	}
	// Create a new VM instance
	if err := createInstance(ctx, instancesClient, projectID, zone, instanceName, instanceTemplate); err != nil {
		return nil, err
	}

	// Get instance external IP
	ipAddress := getInstanceExternalIP(ctx, instancesClient, projectID, zone, instanceName)
	return &vm{
		ip:             ipAddress,
		instanceClient: instancesClient,
		instanceName:   instanceName,
		zone:           zone,
	}, nil
}

func (vm *vm) deployRunner(configFileName string) error {
	fmt.Println("Transferring the config to the VM...")
	// Send the binary to the VM
	fmt.Println("Transferring the binary to the VM...")
	scpCmd := exec.Command("gcloud", "compute", "scp", binaryFilePath, fmt.Sprintf("%s@%s:~/%s", "YOUR_VM_USER", vm.ip, binaryName), "--zone", vm.zone)
	return scpCmd.Run()
}

func (vm *vm) startRunner() error {
	// Execute the binary on the VM
	fmt.Println("Executing the binary on the VM...")
	execCmd := exec.Command("gcloud", "compute", "ssh", fmt.Sprintf("%s@%s", "YOUR_VM_USER", vm.ip), "--zone", vm.zone, "--command", fmt.Sprintf("chmod +x ~/%s && ./%s > %s", binaryName, binaryName, logFileName))
	err := execCmd.Run()
	if err != nil {
		log.Error("Error executing binary: %v", err)
	}
	return err
}

func (vm *vm) deleteRunner(ctx context.Context, projectID string) {
	// Download the log file
	fmt.Println("Downloading the log file...")
	scpLogCmd := exec.Command("gcloud", "compute", "scp", fmt.Sprintf("%s@%s:~/%s", "YOUR_VM_USER", vm.ip, logFileName), ".", "--zone", vm.zone)
	if err := scpLogCmd.Run(); err != nil {
		log.Info("Error downloading log file: %v", err)
	}
	// Delete the VM
	fmt.Println("Deleting the VM...")
	deleteInstance(ctx, vm.instanceClient, projectID, vm.zone, vm.instanceName)
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

	log.Info("Waiting for the operation to complete...")
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
	log.Crit("No external IP found for the instance.")
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

func listZones(projectID string) ([]string, error) {
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
	zones := make([]string, 0)
	it := c.List(ctx, req)
	for {
		zone, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list zones: %v", err)
		}

		zones = append(zones, zone.GetName())
		fmt.Println(zone.GetName())
	}

	return zones, nil
}
