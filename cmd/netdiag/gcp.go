package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"google.golang.org/api/iterator"

	"github.com/autonity/autonity/log"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

const (
	binaryName  = "netdiag"
	logFileName = "output.log"
)

type vm struct {
	id             int
	ip             string
	instanceName   string
	zone           string
	user           string
	instanceClient *compute.InstancesClient // don't forget to close ! defer instancesClient.Close()
}

func deployVM(id int, projectID, instanceName, zone, instanceTemplate, user string) (*vm, error) {
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
		id:             id,
		ip:             ipAddress,
		instanceClient: instancesClient,
		instanceName:   instanceName,
		zone:           zone,
		user:           user,
	}, nil
}

func (vm *vm) deployRunner(configFileName string) error {
	log.Info("Transferring config file to the VM...", "id", vm.id)
	// Send the binary to the VM
	cmd := exec.Command("scp", configFileName, fmt.Sprintf("%s@%s:~/%s", vm.user, vm.ip, configFileName))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Error("command failure", "err", err, "cmd", cmd.String())
		return err
	}
	log.Info("Transferring the binary to the VM...", "id", vm.id)
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	//scpCmd.Stdout = os.Stdout
	//scpCmd.Stderr = os.Stderr
	if err := exec.Command("scp", exePath, fmt.Sprintf("%s@%s:~/%s", vm.user, vm.ip, binaryName)).Run(); err != nil {
		log.Error("command failure", "err", err, "id", vm.id)
	}
	return err
}

func (vm *vm) startRunner(configFileName string) error {
	// Execute the binary on the VM
	log.Info("Executing the runner on the VM...", "id", vm.id)
	flags := fmt.Sprintf("run --config %s --id %d", configFileName, vm.id)
	localCommand := fmt.Sprintf("chmod +x ~/%s && nohup sudo -b ./%s %s > %s 2>&1 ", binaryName, binaryName, flags, logFileName)
	execCmd := exec.Command("ssh", fmt.Sprintf("%s@%s", vm.user, vm.ip), localCommand)
	if err := execCmd.Start(); err != nil {
		log.Error("Error executing binary: %v", err, "id", vm.id)
	}
	return nil
}

func (vm *vm) deleteRunner(ctx context.Context, projectID string) {
	// Download the log file
	log.Info("Downloading the log file...", "id", vm.id)
	scpLogCmd := exec.Command("gcloud", "compute", "scp", fmt.Sprintf("%s@%s:~/%s", "YOUR_VM_USER", vm.ip, logFileName), ".", "--zone", vm.zone)
	scpLogCmd.Stdout = os.Stdout
	scpLogCmd.Stderr = os.Stderr
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
		//fmt.Println(zone.GetName())
	}

	return zones, nil
}
