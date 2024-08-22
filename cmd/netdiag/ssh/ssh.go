package ssh

import (
	"bytes"
	"os"
	"path/filepath"
	"sync"

	scp "github.com/povsister/scp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

var sshPeers sync.Map

type Client struct {
	ip  string
	ssh *ssh.Client
	scp *scp.Client
}

func New(ip, user string) *Client {
	key, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"))
	if err != nil {
		log.Error("Failed to read ssh file, err:", err)
		return nil
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Error("Failed to parse private key, err:", err)
		return nil
	}

	//userName := os.Getenv("USER")
	log.Info("Connecting to", " host: ", ip, " user: ", user)
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// use default ssh port 22
	portip := ip + ":22"
	client, err := ssh.Dial("tcp", portip, config)
	if err != nil {
		log.Error("Unable to connect, err:", err)
		return nil
	}
	// create the scp client
	scpCl, err := scp.NewClientFromExistingSSH(client, &scp.ClientOption{Sudo: false})
	if err != nil {
		log.Error("Unable to setup scp client using ssh client, err:", err)
		return nil
	}
	sshService := &Client{ip: ip, ssh: client, scp: scpCl}
	return sshService
}

func Manager(ip, user string) *Client {
	result, ok := sshPeers.Load(ip)
	if ok {
		return result.(*Client)
	}
	sshService := New(ip, user)
	if sshService == nil {
		return nil
	}
	sshPeers.Store(ip, sshService)
	return sshService
}

func (sshSvc *Client) RunCommand(cmd string) (string, error) {
	sess, err := sshSvc.ssh.NewSession()

	if err != nil {
		log.Error("Unable to create new session, err:", err)
		return "", err
	}
	defer sess.Close()
	var b bytes.Buffer
	sess.Stdout = &b
	err = sess.Run(cmd)
	return b.String(), err
}

func (sshSvc *Client) Download(src, dst string, isFile bool) error {
	var err error
	if isFile {
		err = sshSvc.scp.CopyFileFromRemote(src, dst, &scp.FileTransferOption{PreserveProp: true})
	} else {
		err = sshSvc.scp.CopyDirFromRemote(src, dst, &scp.DirTransferOption{})
	}
	return err
}

func (sshSvc *Client) Upload(src, dst string, isFile bool) error {
	var err error
	if isFile {
		err = sshSvc.scp.CopyFileToRemote(src, dst, &scp.FileTransferOption{PreserveProp: true})
	} else {
		err = sshSvc.scp.CopyDirToRemote(src, dst, &scp.DirTransferOption{})
	}
	return err
}
